package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templateHandler represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	// Obtener el puerto de la variable de entorno PORT asignada por Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Valor por defecto para uso local
	}

	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// Iniciar la sala
	go r.run()

	// Iniciar el servidor web
	log.Println("Starting web server on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

