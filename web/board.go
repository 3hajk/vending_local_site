package web

import (
	"fmt"
	"github.com/3hajk/vending_local_site/controller"
	"github.com/3hajk/vending_machine/controller/board"
	"github.com/gorilla/schema"
	"golang.org/x/text/message"
	"html/template"
	"log"
	"net/http"
)

type boardHandler struct {
	t *template.Template
}

type boardOutsHandler struct {
}

type boardFillingHandler struct {
	t *template.Template
}

type boardTareHandler struct {
	t *template.Template
}

type boardSanitizeHandler struct {
}

func NewBoardHandler(locale *message.Printer) *boardHandler {
	h := boardHandler{}
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("board")
	h.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/board.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func NewOutsHandler() *boardOutsHandler {
	rh := boardOutsHandler{}
	return &rh
}

func NewFillingHandler(locale *message.Printer) *boardFillingHandler {
	rh := boardFillingHandler{}

	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("board")
	rh.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/board.html", "./web/template/header.html", "./web/template/footer.html"))
	return &rh
}

func NewTareHandler(locale *message.Printer) *boardTareHandler {
	rh := boardTareHandler{}
	fmap := template.FuncMap{
		"translate": locale.Sprintf,
	}
	t := template.New("board")
	rh.t = template.Must(t.Funcs(fmap).ParseFiles("./web/template/board.html", "./web/template/header.html", "./web/template/footer.html"))
	return &rh
}

func NewSanitizeHandler() *boardSanitizeHandler {
	rh := boardSanitizeHandler{}

	return &rh
}

func (h *boardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := genBoardConfig()

	stemplate := make([]controller.Option, 15)
	for i := 0; i < 15; i++ {
		stemplate[i] = controller.Option{
			Value: int8(i + 1),
			Text:  fmt.Sprint(float32(i+1)*0.5, " Liter"),
		}
	}
	d.SanitizeTemplate = &stemplate
	h.t.ExecuteTemplate(w, "board", &d)
}

func (h *boardOutsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outForm := new(board.SensorConfig)
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(outForm, r.PostForm)
	if err != nil {
		log.Println("decoder.Decode(drain, r.PostForm)")
	}
	err = nil //  h.cnt.SetSensorConfig(outForm)
	if err != nil {
		log.Println("h.cnt.SetPureParams(drainForm) error: ", err.Error())
	}
	http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
}

//Save drain params
func (h *boardFillingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	drainForm := new(board.PourConfig)
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(drainForm, r.PostForm)
	if err != nil {
		log.Println("decoder.Decode(drain, r.PostForm)")
	}

	msg, err := make(map[string]string), nil //h.cnt.SetPureParams(drainForm)
	if err != nil {

		d := genBoardConfig()
		stemplate := make([]controller.Option, 15)
		for i := 0; i < 15; i++ {
			stemplate[i] = controller.Option{
				Value: int8(i + 1),
				Text:  fmt.Sprint(float32(i+1)*0.5, " Liter"),
			}
		}
		d.SanitizeTemplate = &stemplate
		d.Errors = msg
		d.Success = template.HTML("Error validate params")
		if err != nil {
			w.WriteHeader(502)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		h.t.ExecuteTemplate(w, "board", &d)
	}
	http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
}

func (h *boardTareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tareForm := new(board.TareConfig)
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(tareForm, r.PostForm)
	if err != nil {
		log.Println("decoder.Decode(tareForm, r.PostForm)")
	}
	msg, err := make(map[string]string), nil //h.cnt.SetTareParams(tareForm)
	if err != nil {
		//d, err := h.cnt.GetBoardConfig()
		d := genBoardConfig()
		stemplate := make([]controller.Option, 15)
		for i := 0; i < 15; i++ {
			stemplate[i] = controller.Option{
				Value: int8(i + 1),
				Text:  fmt.Sprint(float32(i+1)*0.5, " Liter"),
			}
		}
		d.SanitizeTemplate = &stemplate
		d.Errors = msg
		d.Success = template.HTML("Error validate params")
		if err != nil {
			w.WriteHeader(502)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		h.t.ExecuteTemplate(w, "board", &d)
	}
	http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
}

func (h *boardSanitizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	decoder := schema.NewDecoder()
	sanitize := new(board.SanitizeConfig)
	err = decoder.Decode(sanitize, r.PostForm)
	if err != nil {
		log.Println("decoder.Decode(sanitize, r.PostForm)", err.Error())
	}
	err = nil //h.cnt.SetSanitizeConfig(sanitize)
	if err != nil {
		log.Println("h.cnt.SetPureParams(sanitize) error: ", err.Error())

	}
	http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
}

func genBoardConfig() *controller.BoardConfig {

	var success string
	errors := make(map[string]string)

	volume := board.VolumeConfig{
		Impulses:       5000,
		StartPressure:  5000,
		DrainTime1:     5000,
		DrainPressure1: 5000,
		DrainTime2:     5000,
		DrainPressure2: 5000,
		CycleSwitch:    5000,
		MaxTime:        5000,
	}

	drain := board.PourConfig{
		Volume1:          &volume,
		Volume2:          &volume,
		Volume3:          &volume,
		Volume4:          &volume,
		ReleaseDuration:  300,
		ReleasePeriod:    300,
		KegPressureDelta: 300,
	}
	serviceMode := true
	tare := board.TareConfig{
		TareMinPressure:   3000,
		TareDetectionTime: 3000,
		TareMaxPressure:   3000,
		TarePressureDelta: 3000,
		TareVol1Min:       3000,
		TareVol1Max:       3000,
		TareVol2Min:       3000,
		TareVol2Max:       3000,
		TareVol3Min:       3000,
		TareVol3Max:       3000,
		TareVol4Min:       3000,
		TareVol4Max:       3000,
	}

	outs := board.SensorConfig{
		ValveProduct: true,
		ValveHolder:  false,
		ValveLock:    true,
		ValveCo2:     false,
		//ValveCarbonation: true,
		ValveWater:    false,
		ValveDrainage: true,
	}

	sntz := board.SanitizeConfig{
		CleaningWater1:    10,
		CleaningDetergent: 10,
		CleaningWater2:    10,
		CleaningProduct:   10,
	}

	status := board.VendingStatus{
		1,
		50,
		"Нет давления",
		795,
		111,
		1,
		1,
		409,
		"Необходимое давление не установлено",
		20,
		96,
		1,
		1,
		1600,
		100500,
		1,
	}

	return &controller.BoardConfig{
		template.HTML("Board"),
		&drain,
		&tare,
		&outs,
		&sntz,
		&status,
		template.HTML(success),
		!serviceMode,
		nil,
		errors,
		nil,
		"",
		1,
	}

}
