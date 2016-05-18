package main

import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"io"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"time"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.tmpl"))
	favicon,err = ioutil.ReadFile("templates/favicon.ico")
	sessionStore = make(map[string]*Session)
)

func main() {
	http.HandleFunc("/favicon.ico", serveFavicon)
	http.HandleFunc("/login", login)
	http.HandleFunc("/", handlerWrapper(sayHello))
	if err := http.ListenAndServe(":8080",nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveFavicon(w http.ResponseWriter, r *http.Request){
	w.Write(favicon)
}


func sayHello(w http.ResponseWriter, r *http.Request, s Session){
	fmt.Fprintf(w, "Hello " + s.user + "!")
}

func login(w http.ResponseWriter, r *http.Request){
	token, _ := GenerateToken(32)
	if r.Method == "GET" {
		renderTemplate(w,"login",token)
	} else {

		if err := r.ParseForm(); err != nil {
			fmt.Printf("%#v\n",err)
		}
		sessionStore[token] = createNewSession(r.FormValue("username"))
		fmt.Printf("%#v\n",r.Form)
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
		fmt.Println("Correct token: ", r.Form["token"][0] )
		http.Redirect(w,r,"/"+token,http.StatusFound)
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
	return base64.URLEncoding.EncodeToString(b),nil
}

func handlerWrapper(inputFunction func(http.ResponseWriter, *http.Request, Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println(r.RequestURI)
		token := strings.Split(r.RequestURI,"/")[1]
		session := sessionStore[token]
		// TODO Really need to look into how long a session should be kept active
		if session == nil || time.Now().After(session.timestamp.Add(8*time.Hour)){
			delete(sessionStore,token)
			fmt.Println("No valid session..")
			http.Redirect(w,r,"/login",http.StatusFound)
			return
		}
		session.updateSession()
		inputFunction(w,r,*session)
	}
}


type Session struct {
	timestamp time.Time
	token string
	user string
}

func (sess *Session) updateSession(){
	sess.timestamp = time.Now()
}

func createNewSession(user string) *Session{
	token, _ := GenerateToken(32)
	return &Session{time.Now(),token,user}
}

