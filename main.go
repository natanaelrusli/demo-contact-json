package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gmux "github.com/gorilla/mux"
)

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

var contacts = []Contact{
	{
		ID:    1,
		Name:  "John Doe",
		Phone: "+628912345678",
		Email: "johndoe@example.com",
	},
	{
		ID:    2,
		Name:  "Samantha Jane",
		Phone: "+628912345999",
		Email: "samantha@example.com",
	},
}

const (
	Addr = ":8081"
)

type ResponseMsg string

const (
	ResponseOK ResponseMsg = "OK"
)

type ContactsResponse struct {
	Message ResponseMsg `json:"message"`
	Data    []Contact   `json:"data"`
}

type CreateContactResponse struct {
	Message ResponseMsg `json:"message"`
	Data    Contact     `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// start coding here...
	mux := gmux.NewRouter()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello!")
	})

	// GET all contacts
	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		res := ContactsResponse{
			Data:    contacts,
			Message: ResponseOK,
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Fprintf(w, "error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// POST create new contact
	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		contact := Contact{}
		err := json.NewDecoder(r.Body).Decode(&contact)

		if err != nil {
			err := ErrorResponse{
				Error: "invalid payload",
			}

			json.NewEncoder(w).Encode(err)
			return
		}

		contacts = append(contacts, contact)

		res := CreateContactResponse{
			Message: ResponseOK,
			Data:    contact,
		}
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(res)
	}).Methods(http.MethodPost)

	srv := http.Server{
		Addr:        Addr,
		Handler:     mux,
		ReadTimeout: 30 * time.Second,
	}

	log.Printf("Server Listening on %s...", Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error when listen and serve: %s", err.Error())
	}
	log.Println("Done...")
}
