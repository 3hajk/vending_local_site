package web

import (
	//"github.com/3hajk/vending_machine/controller"
	"github.com/3hajk/vending_machine/database"
	"github.com/rivo/users"
	"html/template"
	"log"
	"net/http"
)

type systemHandler struct {
	db *database.Database
	t  *template.Template
}

type systemCmdHandler struct {
	db *database.Database
}

func NewSystemHandler(db *database.Database) *systemHandler {
	h := systemHandler{}
	h.db = db
	h.t = template.Must(template.ParseFiles("./web/template/system.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func NewSystemCmdHandler(db *database.Database) *systemCmdHandler {
	h := systemCmdHandler{}
	h.db = db
	return &h
}

func (h *systemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO add loading config data

	token, notifyId := h.db.GetTelegramNotify()
	mode := true // h.cnt.GetServiceMod()

	imei := h.db.GetVendingParamByName("imei")
	proto := h.db.GetVendingParamByName("proto")

	log.Println(token)

	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}

	d := struct {
		Title        string
		NotifyId     int64
		BoardAddress byte
		Imei         string
		Proto        string
		Mode         bool
		User         users.User
		UserEmail    string
	}{
		Title:        "System",
		NotifyId:     notifyId,
		BoardAddress: 1, // h.cnt.GetBoardAddress(),
		Imei:         imei,
		Proto:        proto,
		Mode:         !mode,
		User:         user,
		UserEmail:    email,
	}

	h.t.ExecuteTemplate(w, "system", &d)
}

func (h *systemCmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	if r.ContentLength > 0 {
		switch r.PostForm["cmd"][0] {
		case "proto":
			h.db.SetVendingParam("proto", r.PostForm["value"][0])
			break
		case "imei":
			h.db.SetVendingParam("imei", r.PostForm["value"][0])
			break
		case "notify":
			h.db.SetNotifyId(r.PostForm["value"][0])
			break
		case "address":
			break
		}
	}
	w.WriteHeader(http.StatusOK)
}
