package web

import (
	"github.com/3hajk/vending_machine/database"
	"html/template"
	"net/http"
)

type modemHandler struct {
	db *database.Database
	t  *template.Template
}

func NewModemHandler(db *database.Database) *modemHandler {
	h := modemHandler{}
	h.db = db

	h.t = template.Must(template.ParseFiles("./web/template/modem.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func (h *modemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.t.ExecuteTemplate(w, "modem", &Page{Title: "Modem"})
}
