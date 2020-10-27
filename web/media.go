package web

import (
	"github.com/3hajk/vending_machine/database"
	"github.com/rivo/users"
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
	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}
	h.t.ExecuteTemplate(w, "media", &Page{Title: "Media", User: user, UserEmail: email})
}
