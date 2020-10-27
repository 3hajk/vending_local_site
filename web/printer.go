package web

import (
	"github.com/rivo/users"
	//"github.com/3hajk/vending_machine/printer"
	"html/template"
	"log"
	"net/http"
)

type printerHandler struct {
	t *template.Template
}

func NewPrinterHandler() *printerHandler {
	rh := printerHandler{}

	rh.t = template.Must(template.ParseFiles("./web/template/printer.html", "./web/template/header.html", "./web/template/footer.html"))
	return &rh
}

func (ah *printerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	if r.ContentLength > 0 {
		log.Println(r.PostForm)

		//file := r.PostForm["file"][0]

		//ah.p.PrintLabelBMP2(file)
	}
	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}
	ah.t.ExecuteTemplate(w, "printer", &Page{Title: "Printer", User: user, UserEmail: email})
}
