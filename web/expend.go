package web

import (
	"fmt"
	"github.com/3hajk/vending_machine/database"
	"github.com/gorilla/schema"
	"github.com/rivo/users"
	"html/template"
	"log"
	"net/http"
)

type expendHandler struct {
	db *database.Database
	t  *template.Template
}

func NewExpendHandler(db *database.Database) *expendHandler {
	h := expendHandler{}
	h.db = db

	h.t = template.Must(template.ParseFiles("./web/template/expend.html", "./web/template/header.html", "./web/template/footer.html"))
	return &h
}

func (h *expendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var success string
	var err error
	var valid bool
	expend := new(database.ExpendData)
	errors := make(map[string]string)
	log.Println(r.ParseForm(), r.PostForm, r.ContentLength)
	if r.ContentLength > 0 {

		err := r.ParseForm()
		if err != nil {
			log.Println("err := r.ParseForm()")
		}

		decoder := schema.NewDecoder()
		err = decoder.Decode(expend, r.PostForm)
		if err != nil {
			log.Println("decoder.Decode(person, r.PostForm)")
		}
		log.Println(expend)

		errors, valid = expend.Validate()
		//TODO add saving new data
		if valid {

			err = h.db.SetExpendData(expend)
			if err != nil {
				log.Println("SetExpendData: ", err.Error())
			}
			success = "Saving expend success"
		}

	} else {

		expend, err = h.db.GetExpend()
		if err != nil {
			w.WriteHeader(502)
			fmt.Fprintf(w, "Error writing end-of-print: %v", err)
			return
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
		Title     template.HTML
		Data      *database.ExpendData
		Success   template.HTML
		Errors    map[string]string
		User      users.User
		UserEmail string
	}{
		template.HTML("Expend"),
		expend,
		template.HTML(success),
		errors,
		user,
		email,
	}
	h.t.ExecuteTemplate(w, "expend", &d)
}
