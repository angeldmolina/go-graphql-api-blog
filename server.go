package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-blog-graphql/database"
	"go-blog-graphql/graph"
	"go-blog-graphql/graph/generated"
	"go-blog-graphql/utils"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("mydatabase")
	database.ConnectDb(db)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedHeaders:   []string{"Accept", "Accept-Encoding", "Authorization", "Content-Length", "Content-Type", "X-CSRF-Token"},
		AllowedOrigins:   []string{"http://localhost:3000", "https://myapp.com"},
		AllowCredentials: true,
	}))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{},
		Directives: generated.DirectiveRoot{
			Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				// Implement your authentication logic here
				return next(ctx)
			},
		},
	}))
	// Add custom GraphQL extensions
	srv.Use(extension.FixedComplexityLimit(100))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	// Set up the GraphQL playground
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Client().Ping(context.Background(), nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Use a middleware to check if the user is authorized to access the GraphQL endpoint
	router.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-KEY")
			if apiKey != utils.GetAPIKey() {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Invalid API key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}).Handle("/query", srv)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
