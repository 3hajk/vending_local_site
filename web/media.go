package web

import (
	"github.com/3hajk/vending_machine/database"
	"html/template"
	"net/http"
)

type mediaHandler struct {
	db *database.Database
	t  *template.Template
}

func NewMediaHandler(db *database.Database) *mediaHandler {
	h := mediaHandler{}
	h.db = db

	h.t = template.Must(template.ParseFiles("./web/template/media.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func (h *mediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.t.ExecuteTemplate(w, "media", &Page{Title: "Media"})
}
