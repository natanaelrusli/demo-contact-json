package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// response message type
type ResponseMsg string

const (
	ResponseOK ResponseMsg = "OK"
)

// dto
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ContactsResponse struct {
	Message ResponseMsg `json:"message"`
	Data    []Contact   `json:"data"`
}

type ContactResponse struct {
	Message ResponseMsg `json:"message"`
	Data    Contact     `json:"data"`
}

type CreateContactResponse struct {
	Message ResponseMsg `json:"message"`
	Data    Contact     `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		errorResponse := ErrorResponse{
			Error: "method not allowed",
		}

		json.NewEncoder(w).Encode(errorResponse)
	})
}

func main() {
	// start coding here...
	mux := gmux.NewRouter()
	mux.MethodNotAllowedHandler = MethodNotAllowedHandler()

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

	// GET a single contact by id
	mux.HandleFunc("/contacts/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		vars := gmux.Vars(r)
		userId := vars["id"]

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			errorResponse := ErrorResponse{
				Error: "invalid id",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		for _, user := range contacts {
			if user.ID == userIdInt {
				res := ContactResponse{
					Data:    user,
					Message: ResponseOK,
				}

				err = json.NewEncoder(w).Encode(res)
				if err != nil {
					fmt.Fprintf(w, "error: %s", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				return
			}
		}

		errorResponse := ErrorResponse{
			Error: "data not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}).Methods(http.MethodGet)

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
