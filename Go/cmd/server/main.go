package main

import(
	"log"
	"net/http"

	"github.com/gorilla/mux"

	iwant_controllers "github.com/jayschoen/iWant-slack-bot"
)

type server struct{}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write( []byte( `{"message: "GET"}` ) )
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write( []byte( `{"message: "POST"}` ) )
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write( []byte( `{"message: "PUT"}` ) )
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write( []byte( `{"message: "DELETE"}` ) )
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write( []byte( `{"message: "NOT FOUND"}` ) )
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", get).Methods(http.MethodGet)
	r.HandleFunc("/", post).Methods(http.MethodPost)
	r.HandleFunc("/", put).Methods(http.MethodPut)
	r.HandleFunc("/", delete).Methods(http.MethodDelete)
	r.HandleFunc("/", notFound)

	log.Fatal(http.ListenAndServe(":8080", nil) )
}