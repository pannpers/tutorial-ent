package main

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pannpers/tutorial-ent/ent"
	"github.com/pannpers/tutorial-ent/ent/todo"
)

func main() {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed to open connection to sqlite: %v", err)
	}
	defer client.Close()

	client = client.Debug()

	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed to create schema resources: %v", err)
	}

	task1, err := client.Todo.Create().
		SetText("Add GraphQL example").
		Save(ctx)
	if err != nil {
		log.Fatalf("failed to create a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task1.ID, task1.Text)

	task2, err := client.Todo.Create().
		SetText("Add tracing example").
		Save(ctx)
	if err != nil {
		log.Fatalf("failed to create a todo: %v", err)
	}
	fmt.Printf("%d: %q\n", task2.ID, task2.Text)

	if err := task2.Update().SetParent(task1).Exec(ctx); err != nil {
		log.Fatalf("failed to connect todo2 to its parent: %v", err)
	}

	items, err := client.Todo.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed to query todos: %v", err)
	}
	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}

	// Query all todo items that depend on other items.
	items, err = client.Todo.Query().Where(todo.HasParent()).All(ctx)
	if err != nil {
		log.Fatalf("failed to query todos: %v", err)
	}
	for _, t := range items {
		fmt.Printf("%d: %q\n", t.ID, t.Text)
	}
}
