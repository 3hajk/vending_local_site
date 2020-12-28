package web

import (
	"github.com/3hajk/vending_machine/database"
	"github.com/rivo/users"
	"golang.org/x/text/message"
	"html/template"
	"net/http"
)

type modemHandler struct {
	db *database.Database
	t  *template.Template
}

type Modem struct {
	Apn    string
	User   string
	Passwd string
}

type Wifi struct {
	Ssid   string
	Passwd string
}

func NewModemHandler(db *database.Database, locale *message.Printer) *modemHandler {
	h := modemHandler{}
	h.db = db
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("modem")
	h.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/modem.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func (h *modemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}

	apn := h.db.GetVendingParamByName("apn")
	modemUser := h.db.GetVendingParamByName("user")
	passwd := h.db.GetVendingParamByName("passwd")

	ssid := h.db.GetVendingParamByName("ssid")
	secret := h.db.GetVendingParamByName("secret")

	mode := false

	i := Interface{
		"ppp0", "ip", "127.0.0.1", "/24", "google",
	}
	cnx := &Context{
		Name:                 "test",
		Active:               true,
		Type:                 "",
		Protocol:             "",
		AccessPointName:      "",
		Username:             "user",
		Password:             "passwd",
		AuthenticationMethod: "string",
		Settings:             "",
		Interface:            i,
	}
	modem := &Modem{
		apn,
		modemUser,
		passwd,
	}

	wifi := &Wifi{
		ssid,
		secret,
	}

	roleId := 3
	d := struct {
		Title        string
		Mode         bool
		Modem        *Modem
		Wifi         *Wifi
		ModemContext *Context
		User         users.User
		UserEmail    string
		UserRole     int
	}{
		Title:        "Modem",
		Mode:         !mode,
		Modem:        modem,
		Wifi:         wifi,
		ModemContext: cnx,
		User:         user,
		UserEmail:    email,
		UserRole:     roleId,
	}

	h.t.ExecuteTemplate(w, "modem", &d)
}
