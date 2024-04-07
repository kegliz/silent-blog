package app

import (
	"context"

	"github.com/kegliz/silent-blog/internal/config"
	"github.com/kegliz/silent-blog/internal/post"
	"github.com/kegliz/silent-blog/internal/server/logger"
	"github.com/kegliz/silent-blog/internal/server/router"

	"github.com/kegliz/silent-blog/internal/server"
)

type (
	ServerOptions struct {
		C       *config.Config
		Version string
	}

	appServer struct {
		logger   *logger.Logger
		router   *router.Router
		pService post.Service
		version  string
	}

	appServerOptions struct {
		logger   *logger.Logger
		router   *router.Router
		pService post.Service
		version  string
	}
)

// newAppServer creates a new appServer.
func newAppServer(options appServerOptions) *appServer {
	a := &appServer{
		logger:   options.logger,
		router:   options.router,
		pService: options.pService,
		version:  options.version,
	}
	a.router.SetRoutes(a.routes())
	return a
}

// Listen implements server.Server.
func (a *appServer) Listen(port int, localOnly bool, isTLS bool, domain string) error {
	a.logger.Debug().Str("version", a.version).Msg("debug silent server")
	a.logger.Info().
		Int("port", port).
		Bool("localOnly", localOnly).
		Bool("isTLS", isTLS).
		Str("domain", domain).
		Msg("Starting silent secret service")
	return a.router.Start(port, localOnly, isTLS, domain)
}

// Shutdown implements server.Server for graceful shutdown.
func (a *appServer) Shutdown(ctx context.Context) error {
	return a.router.Shutdown(ctx)
}

// NewServer creates a new server.
func NewServer(options ServerOptions) (server.Server, error) {
	l, r := server.NewLoggerAndRouter(server.EngineOptions{
		Debug: options.C.GetBool("debug"),
	})
	l.Debug().Msgf("Options: posts.file: %s, posts.mddir: %s", options.C.GetString("posts.file"), options.C.GetString("posts.mddir"))
	p, err := post.NewService(post.ServiceOptions{
		Logger:   l,
		FileName: options.C.GetString("posts.file"),
		MdDir:    options.C.GetString("posts.mddir"),
	})
	if err != nil {
		return nil, err
	}
	app := newAppServer(appServerOptions{
		logger:   l,
		router:   r,
		pService: p,
		version:  options.Version,
	})

	return app, nil
}
