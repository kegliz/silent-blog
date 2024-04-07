package app

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/kegliz/silent-blog/internal/post"
	"github.com/kegliz/silent-blog/internal/server/logger"
	"github.com/kegliz/silent-blog/ui"
)

var badRequestErrorMsg = "Bad Request - please contact the administrator"
var internalServerErrorMsg = "Internal Server Error - please contact the administrator"

const (
	webTitleDev  = "Silent Secret DEV"
	webTitleProd = "Silent Secret"
)

var webTitle string = webTitleDev

// TODO: try templ.guide version of gin rendering

// RootHandler is the handler for the / endpoint
func (a *appServer) RootHandler(c *gin.Context) {
	log := a.logger.ContextLoggingFn(c)
	log(logger.DebugLevel).Msg("RootHandler: serving root endpoint")

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	err := ui.Page(webTitle, ui.BaseContent()).Render(c.Request.Context(), c.Writer)

	if err != nil {
		log(logger.ErrorLevel).Err(err).Msg("rendering root failed failed")
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}
}

// HealthHandler is the handler for the /health endpoint
func (a *appServer) HealthHandler(c *gin.Context) {
	log := a.logger.ContextLoggingFn(c)
	log(logger.DebugLevel).Msg("HealthHandler: serving health endpoint")
	c.String(http.StatusOK, "OK")
}

// AboutHandler is the handler for the /about endpoint
func (a *appServer) AboutHandler(c *gin.Context) {
	log := a.logger.ContextLoggingFn(c)
	log(logger.DebugLevel).Msg("AboutHandler: serving about endpoint")
	err := presentSubContent(c, ui.About())
	if err != nil {
		log(logger.ErrorLevel).Err(err).Msg("rendering about failed")
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}
}

// PostsHandler is the handler for the /posts endpoint
func (a *appServer) PostsHandler(c *gin.Context) {
	log := a.logger.ContextLoggingFn(c)
	log(logger.DebugLevel).Msg("PostHandler: serving post endpoint")
	posts, err := a.pService.GetPosts(log)
	if err != nil {
		log(logger.ErrorLevel).Err(err).Msg("getting posts failed")
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}
	err = presentSubContent(c, ui.PostList(posts))
	if err != nil {
		log(logger.ErrorLevel).Err(err).Msg("rendering posts failed")
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}
}

// PresentPost is the handler for the /post/:id endpoint
func (a *appServer) PresentPost(c *gin.Context) {
	log := a.logger.ContextLoggingFn(c)
	log(logger.DebugLevel).Msg("PresentPost: serving post/id endpoint")
	id := c.Param("id")
	if id == "" {
		log(logger.ErrorLevel).Msg("PresentPost: no id provided")
		c.String(http.StatusBadRequest, badRequestErrorMsg)
		return
	}
	log(logger.DebugLevel).Msgf("PresentPost: extracted id %s", id)

	postToPresent, err := a.pService.GetPost(log, id)
	if err != nil {
		log(logger.ErrorLevel).Err(err).Msg("getting post failed")
		var keyError *post.KeyError
		if errors.As(err, &keyError) {
			c.String(http.StatusNotFound, "Post not found")
			return
		}
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}

	var content string
	if postToPresent.FileName != "" {
		log(logger.DebugLevel).Msgf("PresentPost: converting md file to html for post/%s from the file %s", id, postToPresent.FileName)
		content, err = ui.ConvertMdFileToHTML(postToPresent.FileName)
		if err != nil {
			log(logger.ErrorLevel).Err(err).Msgf("converting md file to html failed for post/%s", id)
			c.String(http.StatusInternalServerError, internalServerErrorMsg)
			return
		}
	} else {
		content = postToPresent.Content
	}

	err = presentSubContent(c, ui.Post(postToPresent, content))
	if err != nil {
		log(logger.ErrorLevel).Err(err).Msgf("rendering post/%s failed", id)
		c.String(http.StatusInternalServerError, internalServerErrorMsg)
		return
	}
}

// presentSubContent is a helper function to present sub content
func presentSubContent(c *gin.Context, subContent templ.Component) error {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	var err error
	if c.GetHeader("HX-Request") == "true" {
		err = subContent.Render(c.Request.Context(), c.Writer)
	} else {
		err = ui.Page(webTitle, subContent).Render(c.Request.Context(), c.Writer)
	}
	return err
}
