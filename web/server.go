package web

import (
	"fmt"
	"github.com/3hajk/vending_machine/controller/board"
	"github.com/3hajk/vending_machine/database"
	"github.com/3hajk/vending_machine/printer"
	"github.com/fatih/structs"
	"github.com/rivo/users"
	"html/template"
	"net/http"
	"time"
)

type rootHandler struct {
	//p   *printer.Printer
	//cnt *controller.Controller
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
	//if err != nil {
	//	w.WriteHeader(502)
	//	fmt.Fprintf(w, "Error writing end-of-print: %v", err)
	//	return
	//}
	var msgError string
	//if status.ErrorCode != 0 {
	//	_, err := printer.GetError(status.ErrorCode)
	//	msgError = err.Error()
	//}
	serviceMode := true // rh.cnt.GetServiceMod()
	//boardStatus := rh.cnt.GetStatus()

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
	d := struct {
		Title       template.HTML
		ExpendData  *Expend
		FillData    *FillLog
		Printer     *printer.Status
		Msg         string
		Mode        bool
		BoardStatus *board.VendingStatus
		User        users.User
		UserEmail   string
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
	ah.t.ExecuteTemplate(w, "about", &Page{Title: "About", UserEmail: email, User: user})
}

func NewRootHandler(db *database.Database) (*rootHandler, error) {
	rh := rootHandler{}
	//rh.p = pl
	//rh.cnt = board
	rh.db = db
	rh.t = template.Must(template.ParseFiles("./web/template/main.html", "./web/template/header.html", "./web/template/footer.html"))
	return &rh, nil
}

func NewAboutHandler() *aboutHandler {
	ah := aboutHandler{}
	ah.t = template.Must(template.ParseFiles("./web/template/header.html", "./web/template/footer.html", "./web/template/about.html"))
	return &ah
}
