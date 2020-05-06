package main

import(
	"log"
	"net/http"
)

type server struct{}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			w.Write( []byte( `{ "message": "GET" }` ) )

		case "POST":
			w.WriteHeader(http.StatusCreated)
			w.Write( []byte( `{"message": "POST"}` ) )
		
		case "PUT":
			w.WriteHeader(http.StatusAccepted)
			w.Write( []byte( `{"message": "PUT"}` ) )

		case "DELETE":
			w.WriteHeader(http.StatusOK)
			w.Write( []byte( `{"message": "DELETE"}` ) )

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write( []byte( `{"message": "NOT FOUND"}` ) )
	}
}

func main() {

	http.HandleFunc("/", root)

	log.Fatal(http.ListenAndServe(":8080", nil) )
}