package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Default read/write timeout for http server
const (
	DefaultReadTimeout     = 32 * time.Second
	DefaultWriteTimeout    = 64 * time.Second
	DefaultShutdownTimeout = 30 * time.Second
)

type App struct {
	server      *http.Server
	mux         *mux.Router
	Middlewares []Middleware
	Config      *AppConfig
	redis       *RedisWrapper
	serviceInfo *ServiceInfo
}

type ServiceInfo struct {
	CommitHash string
	NodeID     string
}

func NewApp(ctx context.Context) (*App, error) {
	config := NewAppConfig()
	app := &App{
		Config: config,
	}
	return app, app.setup(ctx)
}

func (app *App) setup(ctx context.Context) error {
	var err error
	// Create a mux for routing incoming requests
	app.mux = mux.NewRouter()
	app.setupMiddlewares(ctx)
	app.setupRoute(ctx)
	app.redis, err = NewRedisWrapper(app.Config.RedisConfig)
	if err != nil {
		return err
	}

	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Config.Port),
		Handler: app.mux,
	}
	hostName, err := os.Hostname()
	if err != nil {
		return err
	}
	app.serviceInfo = &ServiceInfo{
		CommitHash: CommitHash,
		NodeID:     hostName,
	}
	app.server = s

	return nil

}

func (app *App) setupRoute(ctx context.Context) {
	routeHandlers := []*RouterHandler{
		{
			Route: indexRoute, Handler: app.InfoHandler,
		},
		{
			Route: pingRoute, Handler: app.InfoHandler,
		},
	}
	app.AddPublicRouteHandlers(routeHandlers...)
}

func (app *App) setupMiddlewares(ctx context.Context) {
	app.Middlewares = []Middleware{
		&ParseInfoMiddleware{}, &CorsMiddleware{CORS: DefaultCORS},
		&MonitoringMiddleware{},
	}
}

func (app *App) Start(ctx context.Context) error {
	// Continue to process new requests until an error occurs
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	idleConnsClosed := make(chan struct{})

	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer done()

		log.Print(ctx, "prepare to shut down...")
		if err := app.server.Shutdown(shutdownCtx); err != nil {
			log.Print(ctx, "prepare to shut down...", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("server running in port %v", app.Config.Port)
	return app.server.ListenAndServe()
}

//Stop to clear usage
func (app *App) Stop(ctx context.Context) {
	log.Print("application is shutting down")
}

// AddPublicRouteHandlers adds more public route handlers to the app
func (app *App) AddPublicRouteHandlers(routeHandlers ...*RouterHandler) {
	for _, handler := range routeHandlers {
		app.addRouteWithMiddleware(*handler)
	}
}

// addRouteWithMiddleware decorates the Route handler with all middlewares then add to the system
func (app *App) addRouteWithMiddleware(routeHandler RouterHandler) {
	var middlewareFn []func(*RequestContext)
	for _, r := range app.Middlewares {
		middlewareFn = append(middlewareFn, r.Handle)
	}
	middlewareFn = append(middlewareFn, routeHandler.Handler)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		ctx := &RequestContext{
			ctx:         r.Context(),
			Writer:      w,
			Request:     r,
			middlewares: middlewareFn,
			Route:       routeHandler.Route,
			idx:         -1,
		}

		ctx.Next()
	}

	app.mux.HandleFunc(routeHandler.Route.Path, handlerFunc)
}
