package app

import (
	"net/http"

	"github.com/kegliz/silent-blog/internal/server/router"
)

func (a *appServer) routes() []*router.Route {
	return []*router.Route{
		{
			Name:        "root",
			Method:      http.MethodGet,
			Pattern:     "/",
			HandlerFunc: a.RootHandler,
		},
		{
			Name:        "health",
			Method:      http.MethodGet,
			Pattern:     "/health",
			HandlerFunc: a.HealthHandler,
		},
		{
			Name:        "about",
			Method:      http.MethodGet,
			Pattern:     "/about",
			HandlerFunc: a.AboutHandler,
		},
		{
			Name:        "posts",
			Method:      http.MethodGet,
			Pattern:     "/posts",
			HandlerFunc: a.PostsHandler,
		},
		{
			Name:        "post",
			Method:      http.MethodGet,
			Pattern:     "/post/:id", // /post/13 ---- c.Param("id")
			HandlerFunc: a.PresentPost,
		},
	}
}
