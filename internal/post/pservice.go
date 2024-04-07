package post

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/kegliz/silent-blog/internal/server/logger"
)

// pService is the implementation of the Service interface.
type pService struct {
	store map[string]Post
	mdDir string
	sync.RWMutex
	logger *logger.Logger
}

// NewService returns a new Service.
// If a filename is provided in options, the service will attempt to initialize the store from the file.
func NewService(opts ServiceOptions) (Service, error) {
	p := pService{
		store:  make(map[string]Post),
		logger: opts.Logger,
		mdDir:  opts.MdDir,
	}
	if opts.FileName != "" {
		if err := p.initPostsFromJson(opts.FileName, opts.MdDir); err != nil {
			p.logger.Error().Err(err).Msg("initPostsFromJson")
			return nil, fmt.Errorf("NewService: cannot init posts from json: %v", err)
		}
	}
	return &p, nil
}

// GetPost implements Service.
func (s *pService) GetPost(l logger.LoggingFn, id string) (Post, error) {
	l(logger.DebugLevel).Str("id", id).Msg("PostService::GetPost")
	s.RLock()
	defer s.RUnlock()
	post, ok := s.store[id]
	if !ok {
		return Post{}, &KeyError{Key: id, Err: ErrKeyNotExist}
	}
	return post, nil
}

// GetPosts implements Service.
func (s *pService) GetPosts(l logger.LoggingFn) ([]Post, error) {
	l(logger.DebugLevel).Msg("PostService::GetPosts")
	s.RLock()
	defer s.RUnlock()
	posts := make([]Post, 0, len(s.store))
	for _, post := range s.store {
		posts = append(posts, post)
	}

	// order posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts, nil
}

// GetPostsByTag implements Service.
func (s *pService) GetPostsByTag(l logger.LoggingFn, tag string) ([]Post, error) {
	l(logger.DebugLevel).Str("tag", tag).Msg("PostService::GetPostsByTag")
	s.RLock()
	defer s.RUnlock()
	posts := make([]Post, 0, len(s.store))
	for _, post := range s.store {
		for _, t := range post.Tags {
			if t == tag {
				posts = append(posts, post)
				break
			}
		}
	}

	// order posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	return posts, nil
}

// initPostsFromJson initializes the store from a json file.
func (s *pService) initPostsFromJson(fileName string, mdDir string) error {
	s.logger.Debug().Str("filename", fileName).Msg("initPostsFromJson")
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("initPostsFromJson: cannot open file : %v", err)
	}
	defer file.Close()

	var posts []Post
	if err = json.NewDecoder(file).Decode(&posts); err != nil {
		return fmt.Errorf("initPostsFromJson: cannot unmarshal json file : %v", err)
	}

	s.Lock()
	defer s.Unlock()
	for _, p := range posts {
		if p.ID != "" {
			if p.FileName != "" && mdDir != "" {
				p.FileName = mdDir + "/" + p.FileName
			}

			s.store[p.ID] = p
		}
	}

	return nil

}
