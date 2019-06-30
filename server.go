package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gomarkdown/markdown"
)

// Page is represents a simple webpage with a Title and the Body
type Page struct {
	Title string
	Body  string
}

func renderBlogPage(w http.ResponseWriter, name string) error {
	fn := name + ".md"
	t, err := template.ParseFiles("view/template.html")
	if err != nil {
		log.Printf("err on parsing template %v", err)
		return err
	}
	p := &Page{name, string(renderMD(fn))}
	t.Execute(w, p)
	return nil
}

func renderMainPage(w http.ResponseWriter) error {
	fn := "index.md"
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Printf("err on parsing template %v", err)
		return err
	}
	p := &Page{"Home", string(renderMD(fn))}
	t.Execute(w, p)
	return nil
}
func renderMD(filename string) []byte {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("reading file: %v", err)
		return nil
	}

	//extensions := parser.MathJax
	//parser := parser.NewWithExtensions(extensions)

	// htmlFlags := html.CommonFlags | html.HrefTargetBlank
	// opts := html.RendererOptions{Flags: htmlFlags}
	// renderer := html.NewRenderer(opts)

	return markdown.ToHTML(md, nil, nil)

}
func blogPageHandler(w http.ResponseWriter, r *http.Request) {
	n := r.URL.Path[1:]
	log.Printf("hit : %v", r.URL.Path)
	if err := renderBlogPage(w, n); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit : %v", r.URL.Path)
	if err := renderMainPage(w); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/view/", blogPageHandler)
	http.HandleFunc("/", mainPageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
