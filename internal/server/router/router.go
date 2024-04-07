package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kegliz/silent-blog/internal/server/logger"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

type (
	Router struct {
		*gin.Engine
		Logger      *logger.Logger
		Routes      []*Route
		BasePath    string
		HTTPServer  *http.Server
		HTTPSServer *http.Server
	}

	RouterOptions struct {
		Logger   *logger.Logger
		BasePath string
	}

	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc gin.HandlerFunc
	}

	ErrNoServerToShutdown struct{}
)

func (e *ErrNoServerToShutdown) Error() string {
	return "no server to shutdown"
}

// NewRouter creates a new router
func NewRouter(options RouterOptions) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Static("/static", "./public") //TODO it should be configurable

	engine.Use(gin.Recovery())
	engine.Use(requestWrapper(options.Logger))

	router := &Router{
		Engine: engine,
		Routes: []*Route{},
		Logger: options.Logger,
	}
	router.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"error": "not found"}) })
	return router
}

// Start starts the server
func (r *Router) Start(nonTLSport int, localOnly bool, isTLS bool, domain string) error {
	if !isTLS {
		var ip string
		if localOnly {
			ip = "127.0.0.1"
		}
		r.HTTPServer = &http.Server{
			Addr:    fmt.Sprintf(ip+":%d", nonTLSport),
			Handler: r,
		}
		return r.HTTPServer.ListenAndServe()
	} else {
		return r.runTLSLetsEncrypt(domain)
	}
}

// TLS runs the server with Let'sEncrypt provided TLS on default ports
func (r *Router) runTLSLetsEncrypt(domain ...string) error {
	var g errgroup.Group

	r.HTTPServer = &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(redirect),
	}
	r.HTTPSServer = &http.Server{
		Handler: r,
	}

	g.Go(func() error {
		return r.HTTPServer.ListenAndServe()
	})
	g.Go(func() error {
		return r.HTTPSServer.Serve(autocert.NewListener(domain...))
	})

	return g.Wait()
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.RequestURI
	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

// Shutdown gracefully shuts down the server without interrupting any active connections
func (r *Router) Shutdown(ctx context.Context) error {
	if r.HTTPSServer != nil {
		var gShutdown errgroup.Group
		gShutdown.Go(func() error {
			return r.HTTPSServer.Shutdown(ctx)
		})
		gShutdown.Go(func() error {
			return r.HTTPServer.Shutdown(ctx)
		})
		return gShutdown.Wait()
	} else if r.HTTPServer != nil {
		return r.HTTPServer.Shutdown(ctx)
	} else {
		return new(ErrNoServerToShutdown)
	}
}

// SetRoutes sets the routes for the router and registers them in the gin engine
func (r *Router) SetRoutes(routes []*Route) {
	r.Routes = routes
	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			r.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			r.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			r.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			r.DELETE(route.Pattern, route.HandlerFunc)
		}
		r.Logger.Info().Msgf("Route %s %s registered", route.Method, r.BasePath+route.Pattern)
	}

}
