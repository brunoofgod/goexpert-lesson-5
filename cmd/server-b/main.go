package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/brunoofgod/goexpert-lesson-5/cmd/server-b/docs"
	"github.com/brunoofgod/goexpert-lesson-5/cmd/server-b/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Cidade API
// @version 1.0
// @description API que recebe um CEP e retorna o nome da cidade equivalente ao CEP.
// @BasePath /
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Rotas do Swagger
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://"+os.Getenv("HOSTNAME")+"/swagger/doc.json")))

	// Rotas da aplicacao
	r.Get("/get-city-by-zip", handlers.GetCityByZip)

	port := os.Getenv("PORT")

	log.Printf("Servidor rodando na porta %s...", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
