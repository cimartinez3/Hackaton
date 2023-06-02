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

	a := &schemas.TokenRequestSource{}

	if err := json.Unmarshal(req, a); err != nil {
		fmt.Println(err)
		return
	}
	service := service.NewTokenService()
	res, _ := service.ProcessToken(a)
	resOut, _ := json.Marshal(res)
	w.Write(resOut)

}

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := ioutil.ReadAll(r.Body)

	a := &schemas.CardInfoRequest{}

	if err := json.Unmarshal(req, a); err != nil {
		fmt.Println(err)
		return
	}
	res := service.GetCardsInfo(a.Email)

	ans, _ := json.Marshal(res)

	w.Write(ans)
}
func Charge(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("true"))
}

func Email(w http.ResponseWriter, r *http.Request) {
	srv := gateway.NewEmailService()

	req, _ := ioutil.ReadAll(r.Body)

	a := &schemas.CardInfoRequest{}

	if err := json.Unmarshal(req, a); err != nil {
		fmt.Println(err)
		return
	}

	err := srv.SendOTP(a.Email)

	if err == nil {
		w.Write([]byte(fmt.Sprintf("%v", "Message sent successfully")))
	}
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r := mux.NewRouter()
	r.HandleFunc("/tokens", TokensHandler).Methods("POST")
	r.HandleFunc("/cards", CardsHandler).Methods("POST")
	r.HandleFunc("/email", Email).Methods("POST")
	r.HandleFunc("/charge", Charge).Methods("POST")
	handler := c.Handler(r)
	http.ListenAndServe(":80", handler)
}
