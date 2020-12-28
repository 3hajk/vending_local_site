package web

import (
	"github.com/3hajk/vending_machine/printer"
	"github.com/rivo/users"
	"golang.org/x/text/message"
	"html/template"
	"log"
	"net/http"
)

type printerHandler struct {
	t *template.Template
}

func NewPrinterHandler(locale *message.Printer) *printerHandler {
	rh := printerHandler{}
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("printer")
	rh.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/printer.html", "./web/template/header.html", "./web/template/footer.html"))
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
	status := "text"
	pStatus := &printer.Status{
		0,
		0,
		0, 0,
		2, 8,
		42,
		40,
		0,
		9,
		8,
	}

	pSensor := &PrinterSensor{
		"4.2",
		"4.1",
		"0.9",
	}
	pState := &printer.State{
		true,
		false,
		true,
		true,
		false,
		true,
		true,
		false,
		true,
		true,
	}
	pInit := true

	var d = struct {
		Title         template.HTML
		Status        string
		PrinterInit   bool
		Printer       *printer.Status
		PrinterSensor *PrinterSensor
		PrinterState  *printer.State
		User          users.User
		UserEmail     string
	}{
		template.HTML("Vending status"),
		status,
		pInit,
		pStatus,
		pSensor,
		pState,
		user,
		email,
	}
	ah.t.ExecuteTemplate(w, "printer", d)
}
