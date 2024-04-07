package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma/v2"
	"github.com/yuin/goldmark"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ConvertMdFileToHTML markodwn file to HTML
func ConvertMdFileToHTML(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("ConvertMdFileToHTML: cannot open file : %v", err)
	}
	defer file.Close()

	markdown, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("ConvertMdFileToHTML: cannot read file : %v", err)
	}
	custom := chroma.MustNewStyle("custom", chroma.StyleEntries{
		chroma.Background:           "bg:#1d1d1d",
		chroma.Comment:              "#7ec699",
		chroma.Keyword:              "#cc99cd",
		chroma.KeywordDeclaration:   "#cc99cd",
		chroma.KeywordNamespace:     "#cc99cd",
		chroma.KeywordType:          "#cc99cd",
		chroma.Operator:             "#67cdcc",
		chroma.OperatorWord:         "#cdcd00",
		chroma.NameClass:            "#f08d49",
		chroma.NameBuiltin:          "#f08d49",
		chroma.NameFunction:         "#f08d49",
		chroma.NameException:        "bold #666699",
		chroma.NameVariable:         "#21212c",
		chroma.LiteralString:        "#999999",
		chroma.LiteralNumber:        "#f08d49",
		chroma.LiteralStringBoolean: "#f08d49",
		chroma.Text:                 "#21212c",
		chroma.Name:                 "#21212c",
		chroma.Generic:              "#21212c",
	})

	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithCustomStyle(custom))),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	// convert to HTML
	err = md.Convert(markdown, &buf)
	if err != nil {
		return "", fmt.Errorf("ConvertMdFileToHTML: cannot convert file : %v", err)
	}
	return buf.String(), nil
}
