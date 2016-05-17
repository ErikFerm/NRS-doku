package main

import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"io"
	"crypto/rand"
	"encoding/base64"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.tmpl"))
)

func main() {
	http.HandleFunc("/",sayHello)
	http.HandleFunc("/login", login)
	if err := http.ListenAndServe(":8080",nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayHello(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello World!")
}

func login(w http.ResponseWriter, r *http.Request){
	// Needs to be saved in a session.. Getting there
	token, _ := GenerateToken(32)
	if r.Method == "GET" {
		renderTemplate(w,"login",token)
	} else {

		if err := r.ParseForm(); err != nil {
			fmt.Printf("%#v\n",err)
		}

		fmt.Printf("%#v\n",r.Form)
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ",r.Form["password"])
		fmt.Println("Correct token: ", r.Form["token"][0] )
	}
}

func renderTemplate(w io.Writer, file string, entity interface{}){
	if err := templates.ExecuteTemplate(w,file+".tmpl",entity); err != nil {
		fmt.Fprint(w,err.Error())
	}
}

// Generates a token of n bytes
func GenerateToken(n int)(string, error){
	b := make([]byte, n)
	_, err :=rand.Read(b)
	if err != nil {
		return "", err
	}
	fmt.Printf("%#v\n",b)
	return base64.URLEncoding.EncodeToString(b),nil
}