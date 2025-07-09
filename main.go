package main

import (
	"bytes"
	."fmt"
	"html/template"

	//"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type Post struct {
	Filename string
	Content  []byte
	HTML     template.HTML
	meta     metadata
}
type metadata struct {
	Title string    `yaml:"title"`
	Date  time.Time `yaml:"date"`
}

func main() {
	postsDir := "posts"

	markdownPosts, err := readMarkdownFiles(postsDir)
	if err != nil {
		Printf("Error reading markdown files: %v\n", err)
		return
	}

	Printf("Found %d markdown files:\n", len(markdownPosts))
	tmpl, err := template.ParseFiles("blogt.html")
	check(err)

	P := Printf
	for _, post := range markdownPosts {
		post.HTML = template.HTML(bytes.TrimSpace([]byte(post.HTML)))
		outFile := "postsHTML/" + post.Filename + ".html"
		file,e := os.Create(outFile); check(e); defer file.Close()
		tmpl.Execute(file, post)
		P("--- File: %s parsed and executed ---\n", post.Filename)
	}
}

func check(err error) {
	if err != nil {
		Println("%w", err)
	}
}

func readMarkdownFiles(dirPath string) ([]Post, error) {
	var posts []Post

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, Errorf("failed to read directory '%s': %w", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.HasSuffix(filename, ".md") {
			file := filepath.Join(dirPath, filename)
			filedata, err := os.Open(file); check(err); defer filedata.Close()
			var matter metadata
			rest, err := frontmatter.Parse(filedata, &matter); check(err)
			var buf bytes.Buffer
			if err := goldmark.Convert(rest, &buf); err != nil {
				Printf("failed to convert markdown file '%s' to HTML: %v\n", file, err)
			}

			Printf("Metadata: %+v\n", matter)
			posts = append(posts,
				Post{Filename: strings.TrimSuffix(filename,".md"), Content: rest, HTML: template.HTML(buf.Bytes())})
		}
	}

	return posts, nil
}
