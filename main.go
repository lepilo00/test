package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	//"prijava"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gocolly/colly"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
	Email    string
	IsAdmin  bool
}

type DB struct {
	*sql.DB
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func InsertIntoDB(usr User) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testdb") //ne dela
	if err != nil {
		fmt.Println("Error1")
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO user (username, password, email, isAdmin) VALUES (?, ?, ?, ?)`, usr.Username, usr.Password, usr.Email, usr.IsAdmin)
	if err != nil {
		fmt.Println("error2")
	}
}

func HashPass(pass string) (string,error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
    return string(bytes), err
}

func CheckPassHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
    return err == nil
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}


func Registracija(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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

	fmt.Println(usr1.Password)

	//kreiranje hash oblike passworda
	hashP, err:= HashPass(usr1.Password)
	if err!= nil{
		fmt.Println("Napaka pri hashu passworda!")
	}

	checkPass:=CheckPassHash(usr1.Password,hashP)

	
	fmt.Println("Password: ",usr1.Password,"\nHash: ",hashP,"\nMatch: ",checkPass)

	tpl.ExecuteTemplate(w, "login.html", nil)

	InsertIntoDB(usr1)
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/login", Registracija)
	http.ListenAndServe(":9090", nil)

}
