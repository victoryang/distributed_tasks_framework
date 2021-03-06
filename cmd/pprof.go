package cmd

import (
    "context"
    "net/http"
    "net/http/pprof"
    "time"

    "github.com/victoryang/api-gateway/config"
    Log "github.com/victoryang/api-gateway/log"

    "github.com/gorilla/mux"
)

const (
    host = "0.0.0.0"
    port = "9090"

    TimeOutDuration = 5 * time.Second
)

var adminServer *http.Server

type handlerFunc func(http.Handler) http.Handler
var handlerFns = []handlerFunc{
//  SetJwtMiddlewareHandler,
}

func RegisterHandlers (r *mux.Router, handlerFns ...handlerFunc) http.Handler {
    var f http.Handler
    f =r
    for _, hFn := range handlerFns {
        f = hFn(f)
    }
    return f
}

func configureAdminHandler() http.Handler {
    r := mux.NewRouter()
    apiRouter := r.NewRoute().PathPrefix("/").Subrouter()

    /*TODO: get some status of controller back to admin*/
    /*admin := apiRouter.PathPrefix("/admin").Subrouter()
    admin.Methods("GET").Path("/usage").HandlerFunc(SetJwtMiddlewareFunc(getUsage))*/

    apiRouter.Path("/debug/cmdline").HandlerFunc(pprof.Cmdline)
    apiRouter.Path("/debug/profile").HandlerFunc(pprof.Profile)
    apiRouter.Path("/debug/symbol").HandlerFunc(pprof.Symbol)
    apiRouter.Path("/debug/trace").HandlerFunc(pprof.Trace)
    apiRouter.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)

    return RegisterHandlers(r, handlerFns...)
}

func startAdminServer(c *config.AdminSever) {
    serverAddress := c.ListenAddress
    adminServer = &http.Server{
        Addr:           serverAddress,
        // Adding timeout of 10 minutes for unresponsive client connections.
        ReadTimeout:    10 * time.Minute,
        WriteTimeout:   10 * time.Minute,
        Handler:        configureAdminHandler(),
        MaxHeaderBytes: 1 << 20,
    }

    go func() {
        // Configure TLS if certs are available.
        err := adminServer.ListenAndServe()
        if err!= nil {
            Log.Error("Admin server error.")
        }
        Log.Info("Admin server running...")
    }()
}

func stopAdminServer() {
    //TODO
    ctx, cancel := context.WithTimeout(context.Background(), TimeOutDuration)
    defer cancel()

    adminServer.Shutdown(ctx)
}