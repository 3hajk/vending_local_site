package web

import (
	"fmt"
	"github.com/3hajk/vending_machine/controller/board"
	"github.com/3hajk/vending_machine/database"
	"github.com/3hajk/vending_machine/printer"
	"github.com/fatih/structs"
	"github.com/rivo/users"
	"golang.org/x/text/message"
	"html/template"
	"net/http"
	"time"
)

type rootHandler struct {
	db *database.Database
	t  *template.Template
}

type aboutHandler struct {
	t *template.Template
}

type FillLog struct {
	Title       template.HTML
	TableHeader []string
	Data        []database.FillStatus
}

type Expend struct {
	Title       template.HTML
	TableHeader []string
	Data        *database.ExpendData
}

type Page struct {
	Title     string
	User      users.User
	UserEmail string
}

type PrinterSensor struct {
	SensorPE  string
	SensorLAB string
	SensorRUL string
}

type Context struct {
	Name                 string
	Active               bool
	Type                 string
	Protocol             string
	AccessPointName      string
	Username             string
	Password             string
	AuthenticationMethod string
	Settings             string
	Interface            Interface
}

type Interface struct {
	Interface         string `json:"Interface"`
	Method            string `json:"Method"`
	Address           string `json:"Address"`
	Netmask           string `json:"Netmask"`
	DomainNameServers string `json:"DomainNameServers"`
}

func (rh *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("lang")
	currentLang := "ru"
	if err != nil || c.Value != currentLang {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "lang", Value: "ua", Expires: expiration}
		http.SetCookie(w, &cookie)
	}

	expend, err := rh.db.GetExpend()
	if err != nil {
		w.WriteHeader(502)
		fmt.Fprintf(w, "Error writing end-of-print: %v", err)
		return
	}
	expendStat := &Expend{
		Title:       template.HTML("Expend Material"),
		TableHeader: structs.Names(&database.ExpendData{})[1:8],
		Data:        expend,
	}
	fillStatus, err := rh.db.GetFillStatus(10)
	if err != nil {
		w.WriteHeader(502)
		fmt.Fprintf(w, "Error writing end-of-print: %v", err)
		return
	}
	fill := &FillLog{
		Title:       template.HTML("Fill status log"),
		TableHeader: structs.Names(&database.FillStatus{})[1:],
		Data:        fillStatus,
	}
	//status, err := rh.p.GetPrinterStatus()
	status := &printer.Status{
		CommandCode:      42,
		ErrorCode:        0,
		PrinterStatus:    77,
		PrinterStatus2:   77,
		FirmwareVersionH: 77,
		FirmwareVersionL: 77,
		SensorPE:         77,
		SensorLAB:        77,
		SensorTM:         77,
		SensorRUL:        77,
		PRMDarkness:      77,
	}
	var msgError string
	serviceMode := true // rh.cnt.GetServiceMod()

	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}

	//cfg := users.Config

	boardStatus := &board.VendingStatus{
		Id:                 100,
		StatusCode:         100,
		Status:             "100",
		CurrentVolume:      100,
		CurrentPressure:    100,
		CurrentImpulses:    100,
		ErrorCode:          100,
		Error:              "100",
		CurrentTemperature: 100,
		KegPressure:        100,
	}

	outs := board.SensorConfig{
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		false,
		1,
		true,
		false,
		false,
		false,
		false,
	}

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

	d := struct {
		Title         template.HTML
		ExpendData    *Expend
		FillData      *FillLog
		Outs          *board.SensorConfig
		PrinterInit   bool
		Printer       *printer.Status
		PrinterSensor *PrinterSensor
		Msg           string
		Mode          bool
		BoardStatus   *board.VendingStatus
		Modem         *Context
		Date          string
		User          users.User
		UserEmail     string
	}{
		Title:       template.HTML("Vending status"),
		ExpendData:  expendStat,
		FillData:    fill,
		Printer:     status,
		Msg:         msgError,
		Mode:        serviceMode,
		BoardStatus: boardStatus,
		User:        user,
		UserEmail:   email,
		Date:        time.Now().Format("2006.01.02 15:04"),
		Outs:        &outs,
		Modem:       cnx,
	}
	rh.t.ExecuteTemplate(w, "main", d)
}

func (ah *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}
	ver := "14 77 sdfsdf sdf sdf ds"
	mac := "00:ed:45:Le:ww:12"
	var d = struct {
		Title     string
		Mac       string
		Status    string
		User      users.User
		UserEmail string
	}{
		Title:     "About",
		User:      user,
		UserEmail: email,
		Mac:       mac,
		Status:    ver,
	}

	ah.t.ExecuteTemplate(w, "about", &d)
}

func NewRootHandler(db *database.Database, locale *message.Printer) (*rootHandler, error) {
	h := rootHandler{}
	h.db = db
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("main")
	h.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/main.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h, nil
}

func NewAboutHandler(locale *message.Printer) *aboutHandler {
	h := aboutHandler{}
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("main")
	h.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/header.html", "./web/template/footer.html", "./web/template/about.html"))
	return &h
}
