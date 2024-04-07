package post

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kegliz/silent-blog/internal/server/logger"
	"github.com/stretchr/testify/suite"
)

type PostServiceTestSuite struct {
	suite.Suite
	Logger      *logger.Logger
	TestService Service
	LogFn       logger.LoggingFn
}

var testPostData = []Post{
	{
		ID:       "1",
		Title:    "Example Post 1",
		Content:  "This is the body of example post 1.\n\nIt has multiple lines.",
		Tags:     []string{"example", "post"},
		Date:     "2020-01-01",
		FileName: "testdata/blog_1.md",
	},
	{
		ID:       "2",
		Title:    "Example Post 2",
		Content:  "This is the body of example post 2.\n\nIt has multiple lines.",
		Tags:     []string{"example"},
		Date:     "2020-01-02",
		FileName: "",
	},
}

func (s *PostServiceTestSuite) SetupSuite() {
	logger := logger.NewLogger(logger.LoggerOptions{
		Debug: true,
	})
	s.Logger = logger
	s.LogFn = logger.ContextLoggingFn(&gin.Context{})

	// order posts by date
	sort.Slice(testPostData, func(i, j int) bool {
		return testPostData[i].Date > testPostData[j].Date
	})
}

// TestNewService tests the NewService method of the post service with initailization from file
func (s *PostServiceTestSuite) TestNewService() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "testdata/test_posts.json",
	})
	s.Require().NoError(err)
	s.TestService = testService
}

// TestNewServiceNoFile tests the NewService method of the post service when wrong file is provided
func (s *PostServiceTestSuite) TestNewServiceNoFile() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "wrongfilename.json",
	})

	s.Require().Error(err)
	s.Require().ErrorContains(err, "cannot init posts from json") // TODO check for specific error
	s.Require().Nil(testService)
}

// TestInitFromJson tests the initFromJson method of the post service
// TODO: test the case when mdDir is provided
func (s *PostServiceTestSuite) TestInitPostsFromJson() {
	f, err := os.CreateTemp("", "post-*.json")
	s.Require().NoError(err)

	defer f.Close()
	defer os.Remove(f.Name())

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	s.Require().NoError(err)

	err = encoder.Encode(testPostData)
	s.Require().NoError(err)

	testService, err := NewService(ServiceOptions{
		Logger: s.Logger,
	})

	s.Require().NoError(err)
	s.TestService = testService

	err = testService.(*pService).initPostsFromJson(f.Name(), "")
	s.Require().NoError(err)

	posts, err := testService.GetPosts(s.LogFn)
	s.Require().NoError(err)
	s.Require().Len(posts, 2)

	post, err := testService.GetPost(s.LogFn, testPostData[0].ID)
	s.Require().NoError(err)
	s.Require().Equal(testPostData[0], post)
}

// TestInitFromJsonWhenWrongFileName tests the initFromJson method of the post service when wrong file is provided
func (s *PostServiceTestSuite) TestInitPostsFromJsonWhenWrongFileName() {
	testService, err := NewService(ServiceOptions{
		Logger: s.Logger,
	})

	s.Require().NoError(err)
	s.TestService = testService

	err = testService.(*pService).initPostsFromJson("wrongfilename.json", "")
	s.Require().Error(err)
	s.Require().ErrorContains(err, "no such file or directory") // TODO check for specific error
}

// TestInitFromJsonWhenWrongJson tests the initFromJson method of the post service when wrong json is provided
func (s *PostServiceTestSuite) TestInitPostsFromJsonWhenWrongJson() {
	f, err := os.CreateTemp("", "post-*.json")
	s.Require().NoError(err)

	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.WriteString("wrong json")
	s.Require().NoError(err)

	testService, err := NewService(ServiceOptions{
		Logger: s.Logger,
	})

	s.Require().NoError(err)
	s.TestService = testService

	err = testService.(*pService).initPostsFromJson(f.Name(), "")
	s.Require().Error(err)
	s.Require().ErrorContains(err, "invalid character") // TODO check for specific error
}

// TestGetPost tests the GetPost method of the post service
func (s *PostServiceTestSuite) TestGetPost() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "testdata/test_posts.json",
	})
	s.Require().NoError(err)
	s.TestService = testService

	post, err := testService.GetPost(s.LogFn, testPostData[0].ID)
	s.Require().NoError(err)
	s.Require().Equal(testPostData[0], post)
}

// TestGetPostNotFound tests the GetPost method of the post service when the post is not found
func (s *PostServiceTestSuite) TestGetPostNotFound() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "testdata/test_posts.json",
	})
	s.Require().NoError(err)
	s.TestService = testService

	expectedKey := "3"
	_, err = testService.GetPost(s.LogFn, expectedKey)
	s.Require().Error(err, "function should return an error")
	keyError := &KeyError{Key: "3", Err: ErrKeyNotExist}

	// one way to check for specific error
	s.Require().ErrorContains(err, keyError.Error(), "error should contain the key error")
	// another way to check for specific error
	// for faking the error uncomment the following: // err = fmt.Errorf("key does not exist: %s", "3")
	if errors.As(err, &keyError) {
		s.Require().Equal(keyError, err)
		s.Require().Equal(expectedKey, keyError.Key, "key should be 3")
		s.Require().Equal(ErrKeyNotExist, keyError.Err, "error should be ErrKeyNotExist")
	} else {
		s.Require().Fail("error should be a KeyError", "Received: %T", err)
	}
}

// TestGetPosts tests the GetPosts method of the post service
func (s *PostServiceTestSuite) TestGetPosts() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "testdata/test_posts.json",
	})
	s.Require().NoError(err)
	s.TestService = testService

	posts, err := testService.GetPosts(s.LogFn)
	s.Require().NoError(err)
	s.Require().Len(posts, 2)
	s.Require().Equal(testPostData, posts)
}

// TestGetPostsByTag tests the GetPostsByTag method of the post service
func (s *PostServiceTestSuite) TestGetPostsByTag() {
	testService, err := NewService(ServiceOptions{
		Logger:   s.Logger,
		FileName: "testdata/test_posts.json",
	})
	s.Require().NoError(err)
	s.TestService = testService

	posts, err := testService.GetPostsByTag(s.LogFn, "example")
	s.Require().NoError(err)
	s.Require().Len(posts, 2)

	posts, err = testService.GetPostsByTag(s.LogFn, "post")
	s.Require().NoError(err)
	s.Require().Len(posts, 1)
	s.Require().Equal(testPostData[1], posts[0])
}

// TestKeyError tests the KeyError error type
func (s *PostServiceTestSuite) TestKeyError() {
	keyError := KeyError{Key: "test", Err: ErrKeyNotExist}
	s.Require().Equal("key does not exist: test", keyError.Error())
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}
