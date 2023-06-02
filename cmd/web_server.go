package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"kushki/hackaton/gateway"
	"kushki/hackaton/schemas"
	"kushki/hackaton/service"
	"net/http"
)

func TokensHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Init Token Service")

	req, _ := ioutil.ReadAll(r.Body)

	a := &schemas.TokenRequest{}

	if err := json.Unmarshal(req, a); err != nil {
		fmt.Println(err)
		return
	}
	service := service.NewTokenService()
	service.ProcessToken(a)

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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r := mux.NewRouter()
	r.HandleFunc("/tokens", TokensHandler).Methods("POST")
	r.HandleFunc("/cards", CardsHandler).Methods("GET")
	r.HandleFunc("/email", Email).Methods("GET")
	handler := c.Handler(r)
	http.ListenAndServe(":80", handler)
}
