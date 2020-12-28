package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/3hajk/vending_local_site/translation"
	"github.com/3hajk/vending_local_site/web"
	"github.com/3hajk/vending_machine/database"
	"github.com/rivo/sessions"
	"github.com/rivo/users"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	sqlitePath = flag.String("sqlitepatch", "./.vending.storage.sqlite", "vending storage sqlite path")
	port       = flag.Int("port", 8080, "The port the server should listen on")
)

func main() {

	flag.Parse()
	log.Println("Starting server")

	err := database.CreateDatabases(sqlitePath)
	if err != nil {
		log.Fatal("Cant connect to db", err.Error())
	}

	db, err := database.ConnectToDB(sqlitePath)
	if err != nil {
		log.Fatal("Cant connect to db: ", err.Error())
	}

	lang := db.GetVendingParamByName("lang")
	if len(lang) < 2 {
		lang = "en"
	}
	locale := translation.GetLocale(lang)

	rh, err := web.NewRootHandler(db, locale)
	if err != nil {
		log.Fatalf("Could not make root handler: %v", err)
	}
	http.Handle("/", rh)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/template/assets"))))

	db.InitUsers()
	users.Config.Internationalization = true
	users.Config.HTMLTemplateDir = "web/template/users"
	users.Config.NewUser = func() users.User {
		return db.NewUser(sessions.CUID())
	}
	users.Config.LoadUserByEmail = func(email string) (users.User, error) {

		log.Print("TEst sdfsdfsdf")
		return db.LoadUserByEmail(email)
	}

	http.HandleFunc(users.Config.RouteLogIn, users.LogIn)
	http.HandleFunc(users.Config.RouteLogOut, users.LogOut)

	exph := web.NewExpendHandler(db, locale)
	http.Handle("/expend_setup", exph)

	brdh := web.NewBoardHandler(locale)
	http.Handle("/board_setup", brdh)
	outsh := web.NewOutsHandler()
	http.Handle("/outs", outsh)
	fillh := web.NewFillingHandler(locale)
	http.Handle("/filling", fillh)
	tareh := web.NewTareHandler(locale)
	http.Handle("/tare", tareh)
	sntzh := web.NewSanitizeHandler()
	http.Handle("/sanitize", sntzh)

	lph := web.NewPrinterHandler(locale)
	http.Handle("/printer_setup", lph)

	mdmh := web.NewModemHandler(db, locale)
	http.Handle("/modem_setup", mdmh)
	sh := web.NewSystemHandler(db, locale)
	http.Handle("/vending_setup", sh)
	sch := web.NewSystemCmdHandler(db)
	http.Handle("/system_cmd", sch)
	mmdh := web.NewMediaHandler(db, locale)
	http.Handle("/media_setup", mmdh)

	ah := web.NewAboutHandler(locale)
	http.Handle("/about", ah)

	bvh := web.NewBoardViewHandler(locale)
	http.Handle("/board_view", bvh)
	bch := web.NewBoardCmdHandler()
	http.Handle("/board_cmd", bch)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	log.Println("press CTRL+C Properly stop the server.")

	// create a channel to receive OS signals
	c := make(chan os.Signal)

	// relay os.Interrupt to channel (os.Interrupt = CTRL+C)
	signal.Notify(c, os.Interrupt)

	// block main routine until a signal is received
	// as long as user doesn't press CTRL+C, our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	log.Println("Stopping the server...")

	log.Println("Closing db connection...")
	//db.Close()
	server := http.Server{Handler: rh}
	server.Shutdown(context.TODO())
	log.Println("END")
}
