package main

import (
	"context"
	"log"
	"net/http"

	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
	tutorial_ent "github.com/pannpers/tutorial-ent"
	"github.com/pannpers/tutorial-ent/ent"
	"github.com/pannpers/tutorial-ent/ent/migrate"
)

func main() {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed to open ent client: %v", err)
	}

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatalf("failed to create ent schema: %v", err)
	}

	// configure the server and start listening on :8081
	srv := handler.NewDefaultServer(tutorial_ent.NewSchema(client))

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)

	log.Println("listening on :8081")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
