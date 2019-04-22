package main

import (
	"go-contacts/gql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
        "go-contacts/controllers"
        "go-contacts/args"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

        my_graphql_handler := handler.GraphQL(
                gql.NewExecutableSchema(gql.NewRootResolvers()),
                handler.ComplexityLimit(args.Graphql_complexity),
            )

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{}})))
	http.HandleFunc("/graphql", controllers.WeakSecureMiddleware(my_graphql_handler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
