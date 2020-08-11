package main

import (
	"context"
	"flag"
	"fmt"
	//"github.com/3hajk/vending_machine/controller"
	"github.com/3hajk/vending_machine/database"
	//"github.com/3hajk/vending_machine/printer"
	"github.com/3hajk/vending_site/web"
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

	//controller := &controller.Controller{
	//
	//}
	//labelPrinter := &printer.Printer{
	//
	//}

	err := database.CreateDatabases(sqlitePath)
	if err != nil {
		log.Fatal("Cant connect to db", err.Error())
	}

	db, err := database.ConnectToDB(sqlitePath)
	if err != nil {
		log.Fatal("Cant connect to db: ", err.Error())
	}

	rh, err := web.NewRootHandler(db) //labelPrinter, controller, db)
	if err != nil {
		log.Fatalf("Could not make root handler: %v", err)
	}
	http.Handle("/", rh)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/template/assets"))))

	exph := web.NewExpendHandler(db)
	http.Handle("/expend", exph)

	brdh := web.NewBoardHandler()
	http.Handle("/board", brdh)
	outsh := web.NewOutsHandler()
	http.Handle("/outs", outsh)
	fillh := web.NewFillingHandler()
	http.Handle("/filling", fillh)
	tareh := web.NewTareHandler()
	http.Handle("/tare", tareh)
	sntzh := web.NewSanitizeHandler()
	http.Handle("/sanitize", sntzh)

	lph := web.NewPrinterHandler()
	http.Handle("/printer", lph)

	mdmh := web.NewModemHandler(db)
	http.Handle("/modem", mdmh)
	sh := web.NewSystemHandler(db)
	http.Handle("/system", sh)
	sch := web.NewSystemCmdHandler(db)
	http.Handle("/system_cmd", sch)
	mmdh := web.NewMediaHandler(db)
	http.Handle("/media", mmdh)

	ah := web.NewAboutHandler()
	http.Handle("/about", ah)

	bvh := web.NewBoardViewHandler()
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
