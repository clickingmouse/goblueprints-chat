package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/clickingmouse/blueprints/chat/trace"
)

//template represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// serveHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	//t.templ.Execute(w, nil)
	t.templ.Execute(w, r)

}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	//
	r := newRoom()

	// removable to turn off tracer
	r.tracer = trace.New(os.Stdout)

	//newRoom
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(`<html><head><title>Chat</title><head/><body>Let's Chat</body></html>`))

	// })

	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		//	if err := http.ListenAndServe(":8080", nil); err != nil {

		log.Fatal("ListenAndServe:", err)
	}
}
