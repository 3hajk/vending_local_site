package web

import (
	"fmt"
	"github.com/3hajk/vending_machine/database"
	"github.com/google/logger"
	"github.com/rivo/users"
	"golang.org/x/text/message"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type systemHandler struct {
	db *database.Database
	t  *template.Template
}

type systemCmdHandler struct {
	db *database.Database
}

func NewSystemHandler(db *database.Database, locale *message.Printer) *systemHandler {
	h := systemHandler{}
	h.db = db
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("modem")
	h.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/system.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func NewSystemCmdHandler(db *database.Database) *systemCmdHandler {
	h := systemCmdHandler{}
	h.db = db
	return &h
}

type Option struct {
	Value    int8
	Id, Text string
	Selected bool
}

type Scheduler struct {
	Start    string
	Stop     string
	Timezone int8
}

func (h *systemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO add loading config data

	_, notifyId := h.db.GetTelegramNotify()
	mode := false
	imei := h.db.GetVendingParamByName("imei")
	proto := h.db.GetVendingParamByName("proto")
	start := h.db.GetVendingParamByName("enable")
	stop := h.db.GetVendingParamByName("disable")
	zone := h.db.GetVendingParamByName("timezone")
	verify := h.db.GetVendingParamByName("verify")
	expends := h.db.GetVendingParamByName("expends")
	labels := h.db.GetVendingParamByName("labels")
	product := h.db.GetVendingParamByName("product")
	t_max := h.db.GetVendingParamByName("t_max")
	board_watchdog := h.db.GetVendingParamByName("board_watchdog")

	i, err := strconv.ParseInt(zone, 10, 8) //Is third parameter not honored??
	if err != nil {
		logger.Error("error ParseInt: ", err.Error())
	}

	scheduler := &Scheduler{
		start,
		stop,
		int8(i),
	}

	stemplate := make([]Option, 26)
	for i := 0; i < 26; i++ {
		text := ""
		value := int8(i - 12)
		if value < 0 {
			text = fmt.Sprintf("(GMT%03d:00)", value)
		} else {
			text = fmt.Sprintf("(GMT+%02d:00)", value)
		}
		stemplate[i] = Option{
			Value: value,
			Text:  text,
		}
	}

	lang := h.db.GetVendingParamByName("lang")
	ltemplate := make([]Option, 3)
	for i, l := range []string{"en", "ua", "ru"} {
		ltemplate[i] = Option{
			Id:   l,
			Text: l,
		}
	}
	expendNotify := h.db.GetVendingParamByName("expend-notify")
	expendTemplate := make([]Option, 2)
	for i, l := range []string{"yes", "no"} {
		expendTemplate[i] = Option{
			Id:   l,
			Text: l, //h.cnt.Translate(l),
		}
	}

	user, _, _ := users.IsLoggedIn(w, r)
	email := ""
	if user != nil {
		email = user.GetEmail()
	} else {
		http.Redirect(w, r, users.Config.RouteLogIn, 302)
	}

	d := struct {
		Title         string
		NotifyId      int64
		BoardAddress  byte
		Imei          string
		Proto         string
		Verify        string
		Mode          bool
		Lang          string
		LangTpl       *[]Option
		Scheduler     *Scheduler
		STempl        *[]Option
		ExpendNotify  string
		ExpendTpl     *[]Option
		Expends       string
		Labels        string
		Product       string
		Temperature   string
		BoardWatchdog string
		User          users.User
		UserEmail     string
		UserRole      int
	}{
		Title:         "System",
		NotifyId:      notifyId,
		BoardAddress:  1,
		Imei:          imei,
		Proto:         proto,
		Verify:        verify,
		Mode:          !mode,
		Lang:          lang,
		LangTpl:       &ltemplate,
		Scheduler:     scheduler,
		STempl:        &stemplate,
		ExpendNotify:  expendNotify,
		ExpendTpl:     &expendTemplate,
		Expends:       expends,
		Labels:        labels,
		Product:       product,
		Temperature:   t_max,
		BoardWatchdog: board_watchdog,
		User:          user,
		UserEmail:     email,
		UserRole:      3,
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
