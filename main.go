package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type (
	D struct{ T template.HTML }
	M struct {
		Title string    `yaml:"title"`
		Date  time.Time `yaml:"date"`
	}
)

func main() {
	renderBlog()
	renderIndex()
}

func renderIndex() {
	t, e := template.ParseFiles("s/t/base.html")
	if e != nil {
		log.Fatal("couldn't open base template")
	}
	f, e := os.ReadFile("s/t/index.html")
	if e != nil {
		log.Fatal("couldn't open index.html")
	}
	var buf bytes.Buffer
	if e := t.Execute(&buf, D{T: template.HTML(f)}); e != nil {
		log.Fatal("couldn't execute index template")
	}
	os.WriteFile("s/index.html", buf.Bytes(), 0o644)
}

func renderBlog() {
	t, e := template.ParseFiles("s/t/base.html")
	if e != nil {
		log.Fatal("couldn't open base template")
	}
	post, e := template.ParseFiles("s/t/post.html")
	if e != nil {
		log.Fatal("couldn't open post template")
	}
	filepath.WalkDir("posts", func(p string, d fs.DirEntry, e error) error {
		if e != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return e
		}
		f, e := os.Open(p)
		if e != nil {
			return e
		}
		var meta M
		rest, e := frontmatter.Parse(f, &meta)
		if e!=nil{return e}
		var buf, out bytes.Buffer
		if e := post.Execute(&buf, meta); e!=nil {return e}
		if e := goldmark.Convert(rest, &buf); e != nil {
			return e
		}
		b := bytes.NewBuffer(bytes.TrimSpace(buf.Bytes()))
		if e := t.Execute(&out, D{T: template.HTML(b.Bytes())}); e != nil {
			return e
		}
		o := filepath.Join("s", strings.TrimSuffix(d.Name(), ".md")+".html")
		os.WriteFile(o, out.Bytes(), 0o644)

		return nil
	})
}
