package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/brunoofgod/goexpert-lesson-5/cmd/server-b/docs"
	"github.com/brunoofgod/goexpert-lesson-5/cmd/server-b/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {
	zipkinURL := os.Getenv("ZIPKIN_URL")
	if zipkinURL == "" {
		zipkinURL = "http://zipkin:9411/api/v2/spans" // URL do Zipkin no docker-compose
	}

	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatalf("Erro ao criar exportador Zipkin: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("server-b"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		_ = tp.Shutdown(context.Background())
	}
}

func main() {
	cleanup := initTracer()
	defer cleanup()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Rotas do Swagger
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://"+os.Getenv("HOSTNAME")+"/swagger/doc.json")))

	// Rotas da aplicação (agora instrumentadas com OpenTelemetry)
	r.Get("/get-temperature-by-zipcode", otelhttp.NewHandler(http.HandlerFunc(handlers.GetTemperatureByZipCode), "GetTemperatureByZipCode").ServeHTTP)

	port := os.Getenv("PORT")

	log.Printf("Servidor rodando na porta %s...", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
