package post

import (
	"errors"

	"github.com/kegliz/silent-blog/internal/server/logger"
)

var ErrKeyNotExist = errors.New("key does not exist")

type (

	// ServiceOptions is a struct that contains the options for constructing a Service.
	ServiceOptions struct {
		Logger   *logger.Logger
		FileName string
		MdDir    string
	}

	// Service is an interface that defines the methods of the Service.
	Service interface {
		// GetPosts returns all posts. No pagination is implemented.
		GetPosts(l logger.LoggingFn) ([]Post, error)
		// GetPost returns a post by its ID.
		GetPost(l logger.LoggingFn, id string) (Post, error)
		// GetPostsByTag returns all posts with a given tag.
		GetPostsByTag(l logger.LoggingFn, tag string) ([]Post, error)
	}

	// Post is a struct that contains the fields of a post.
	Post struct {
		ID       string   `json:"id"`
		Title    string   `json:"title"`
		Tags     []string `json:"tags"`
		Date     string   `json:"date"`
		Content  string   `json:"content"`
		FileName string   `json:"filename"`
	}

	// KeyError is an error type that is returned when a key is not found in the store.
	KeyError struct {
		Key string
		Err error
	}
)

// Error implements the error interface.
func (e *KeyError) Error() string {
	return e.Err.Error() + ": " + e.Key
}
