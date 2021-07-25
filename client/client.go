package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var mySigningKey = []byte("supersecret")

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	ID          string `json:"ID"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Namalengkap string `json:"namalengkap"`
	Foto        string `json:"foto"`
}

type Image struct {
	Id       int       `json:"id"`
	Location string    `json:"location"`
	Path     string    `json:"path"`
	Date     time.Time `json:"date"`
}

var db *gorm.DB
var err error

func Login(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User

	db.First(&user, user.Username)

	if user.Password != credentials.Password || user.Username != credentials.Username {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	http.Redirect(w, r, "/", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("uploadfile")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err3 := ioutil.ReadAll(file)
	if err3 != nil {
		fmt.Println(err3)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	namalengkap := r.FormValue("namalengkap")

	var base64Encoding string
	mimeType := http.DetectContentType(fileBytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += b64.StdEncoding.EncodeToString([]byte(fileBytes))

	foto := base64Encoding

	user := &User{"", username, password, namalengkap, foto}

	db.Create(&user)

	res := Result{Code: 200, Data: user, Message: "Success Register"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))
}

func ClearSession(r http.ResponseWriter, w *http.Request) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(r, cookie)
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/login", Login)
	myRouter.HandleFunc("/home", Home)
	myRouter.HandleFunc("/logout", ClearSession)
	myRouter.HandleFunc("/register", Register).Methods("POST")

	log.Fatal(http.ListenAndServe(":9998", myRouter))
}

func main() {
	fmt.Println("Client")

	db, err = gorm.Open("mysql", "root:@/testdb?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection Success")
	}
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	handleRequests()

}
