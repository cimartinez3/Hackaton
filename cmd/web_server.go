package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"kushki/hackaton/gateway"
	"kushki/hackaton/schemas"
	"net/http"
)

func TokensHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Init Token Service")

	mg := gateway.NewMongoService()

	req, _ := ioutil.ReadAll(r.Body)

	a := schemas.TokenRequest{}

	if err := json.Unmarshal(req, &a); err != nil {
		fmt.Println(err)
		return
	}
	res := mg.PutDocument(a)

	fmt.Println(res)

	w.Write([]byte("success"))

}

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cards Token Service")

	mg := gateway.NewMongoService()

	res := mg.GetDocument("laperra@kushki.com")

	fmt.Println(res)
}

func Email(w http.ResponseWriter, r *http.Request) {
	e := gateway.NewEmailService()

	a := e.SendOTP("")

	fmt.Println(a)
}

func main() {
	r := mux.NewRouter()
	r.Host("www.kushihackaton.com")
	r.HandleFunc("/tokens", TokensHandler).Methods("POST")
	r.HandleFunc("/cards", CardsHandler).Methods("GET")
	r.HandleFunc("/email", Email).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":80", r)
}
