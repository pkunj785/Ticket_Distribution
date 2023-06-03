package main

import (
	"log"
	"net/http"

	"com.routee/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/arena/api/create", controllers.CreateTicket)
		r.Get("/arena/api/get/{id}", controllers.GetTicketData)
		r.Delete("/arena/api/cancel/{id}", controllers.CancelTicket)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route dosen't exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	log.Println("Server running at 8050 port")
	http.ListenAndServe(":8050", r)

}
