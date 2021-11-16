package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	//"prijava"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gocolly/colly"
)

type User struct {
	Username string
	Password string
	Email    string
	IsAdmin  bool
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func Collect(usr User) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		fmt.Println("Error1")
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO user (username, password, email, isAdmin) VALUES (?, ?, ?, ?)`, usr.Username, usr.Password, usr.Email, usr.IsAdmin)
	if err != nil {
		fmt.Println("error2")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Napaka1")
	}

	usr1 := User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
		IsAdmin:  false,
	}
	fmt.Println(usr1)

	tpl.ExecuteTemplate(w, "login.html", usr1)

	Collect(usr1)
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/login", Login)
	http.ListenAndServe(":9090", nil)

}
