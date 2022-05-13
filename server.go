package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aniruddha2000/hackernews/internal/auth"
	"github.com/aniruddha2000/hackernews/internal/pkg/db/migrations/mysql"
	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aniruddha2000/hackernews/graph"
	"github.com/aniruddha2000/hackernews/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	mysql.InitDB()
	mysql.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
