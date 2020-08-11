package web

import (
	"github.com/3hajk/vending_machine/controller"
	"github.com/3hajk/vending_machine/controller/board"
	"html/template"
	"log"
	"net/http"
)

type boardViewHandler struct {
	t *template.Template
}

type boardCmdHandler struct {
	cnt *controller.Controller
}

func NewBoardViewHandler() *boardViewHandler {
	h := boardViewHandler{}
	h.t = template.Must(template.ParseFiles("./web/template/board_cmd.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func NewBoardCmdHandler() *boardCmdHandler {
	rh := boardCmdHandler{}
	return &rh
}

func (h *boardViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.t.ExecuteTemplate(w, "board_cmd", &Page{Title: "Board Control"})
}

func (h *boardCmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("err := r.ParseForm()")
	}
	if r.ContentLength > 0 {
		log.Println(r.PostForm)
		switch r.PostForm["cmd"][0] {
		case "enableLock":
			h.cnt.Command(board.Command.EnableLock, nil)
			return
		case "disableLock":
			h.cnt.Command(board.Command.DisableLock, nil)
			return
		case "oldHead":
			h.cnt.Command(board.Command.UseOldHead, nil)
			return
		case "newHead":
			h.cnt.Command(board.Command.UseNewHead, nil)
			return
		case "usePressureSwitch":
			h.cnt.Command(board.Command.UsePressureSwitch, nil)
			return
		case "unUsePressureSwitch":
			h.cnt.Command(board.Command.UnUsePressureSwitch, nil)
			return
		case "vol1":
			h.cnt.Command(board.Command.StartCalibrationVol1, nil)
			return
		case "vol2":
			h.cnt.Command(board.Command.StartCalibrationVol2, nil)
			return
		case "vol3":
			h.cnt.Command(board.Command.StartCalibrationVol3, nil)
			return
		case "vol4":
			h.cnt.Command(board.Command.StartCalibrationVol4, nil)
			return
		case "220":
			h.cnt.Command(board.Command.StartCalibrationMode, nil)
			return
		case "230":
			h.cnt.Command(board.Command.StartSanitise, nil)
			return
		case "240":
			h.cnt.Command(board.Command.StartFillingSystem, nil)
			return
		case "250":
			h.cnt.Command(board.Command.StopFillingSystem, nil)
			return
		case "241":
			h.cnt.Command(board.Command.Start, nil)
			return
		case "251":
			h.cnt.Command(board.Command.Stop, nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
