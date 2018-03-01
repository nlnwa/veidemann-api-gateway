package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	api "github.com/nlnwa/veidemann-ws-api-gateway/veidemann_api"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"net/http"
	"path"
	"strings"
)

var (
	controllerEndpoint = flag.String("controller_endpoint", "localhost:50051", "endpoint of Controller")
	listenAddress      = flag.String("listenAddress", ":3010", "Address this server listens to. Can be just ':port' to listen on all interfaces")
	swaggerDir         = flag.String("swagger_dir", "html", "path to the directory which contains swagger definitions")
	pathPrefix         = flag.String("path_prefix", "/", "http router prefix path")
)

func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	if err := api.RegisterControllerHandlerFromEndpoint(ctx, mux, *controllerEndpoint, dialOpts); err != nil {
		return nil, err
	}
	if err := api.RegisterReportHandlerFromEndpoint(ctx, mux, *controllerEndpoint, dialOpts); err != nil {
		return nil, err
	}
	if err := api.RegisterStatusHandlerFromEndpoint(ctx, mux, *controllerEndpoint, dialOpts); err != nil {
		return nil, err
	}

	return mux, nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		glog.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	glog.Infof("Serving %s", r.URL.Path)
	p := path.Join(*swaggerDir, r.URL.Path)
	http.ServeFile(w, r, p)
}

func serveSwaggerUi(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swaggerui/")
	p = path.Join(*swaggerDir, r.URL.Path)
	http.ServeFile(w, r, p)
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}

func run(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc(*pathPrefix+"/swagger/", serveSwagger)
	mux.HandleFunc(*pathPrefix+"/swaggerui/", serveSwaggerUi)

	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle(*pathPrefix, gw)

	return http.ListenAndServe(address, allowCORS(wsproxy.WebsocketProxy(mux)))
}

//go:generate scripts/get-dependencies.sh
//go:generate scripts/build-protobuf.sh
func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("Starting api-gateway")
	glog.Info("Connecting to ", *controllerEndpoint)
	glog.Info("Listening on ", *listenAddress)
	glog.Info("Serving from ", *pathPrefix)
	if err := run(*listenAddress); err != nil {
		glog.Fatal(err)
	}
}
