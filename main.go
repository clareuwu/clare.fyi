package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type D struct{ T template.HTML }

func main() {
	t, e := template.ParseFiles("blogt.html")
	if e != nil {
		log.Fatal("couldn't open blog template")
	}
	filepath.WalkDir("posts", func(p string, d fs.DirEntry, e error) error {
		if e != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return e
		}
		f, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		var buf bytes.Buffer
		if e := goldmark.Convert(f, &buf); e != nil {
			return e
		}
		buf = *bytes.NewBuffer(bytes.TrimSpace(buf.Bytes()))
		if e := t.Execute(&buf, D{T: template.HTML(buf.Bytes())}); e != nil {
			return e
		}
		o := filepath.Join("postsHTML", strings.TrimSuffix(d.Name(), ".md")+".html")
		os.WriteFile(o, buf.Bytes(), 0o644)

		return nil
	})
}
