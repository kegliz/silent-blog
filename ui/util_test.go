package ui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConvertMdFileToHTML is a test for ConvertMdFileToHTML
func TestConvertMdFileToHTML(t *testing.T) {
	assert := assert.New(t)

	html, err := ConvertMdFileToHTML("testdata/blog_1.md")
	assert.Nil(err)
	assert.Contains(html, "<h1>Heading 1</h1>")
	assert.Contains(html, "<h2>Heading 2</h2>")
	assert.Contains(html, "<h3>Heading 3</h3>")
	// write it out into a file to check it later
	os.WriteFile("testdata/blog_1.html", []byte(html), 0644)

}
