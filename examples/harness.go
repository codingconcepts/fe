package main

import (
	"context"
	"log"
	"time"

	repo "github.com/codingconcepts/fe/examples/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.SetFlags(0)

	db, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	df := repo.NewDatabaseFunctions(db)

	log.Println("AddPerson")
	if err = df.AddPerson(context.Background(), "Rob Reid", time.Date(1986, 9, 28, 0, 0, 0, 0, time.UTC), "uk"); err != nil {
		log.Fatalf("error adding person: %v", err)
	}

	log.Println("\nGetOldestPerson")
	name, err := df.GetOldestPerson(context.Background())
	if err != nil {
		log.Fatalf("error getting oldest person: %v", err)
	}
	log.Printf("oldest person: %s", name)

	log.Println("\nPeopleBornOn")
	count, err := df.PeopleBornOn(context.Background(), time.Date(1986, 9, 28, 0, 0, 0, 0, time.UTC))
	if err != nil {
		log.Fatalf("error getting people born on: %v", err)
	}
	log.Printf("people born on: %d", count)

	log.Println("\nPersonById")
	byID, err := df.PersonById(context.Background(), "a58933a1-c24f-43d9-bb53-6a1aa3170a12")
	if err != nil {
		log.Fatalf("error getting person by id: %v", err)
	}
	log.Println(byID)

	log.Println("\nPeopleBetween")
	between, err := df.PeopleBetween(context.Background(), "a58933a1-c24f-43d9-bb53-6a1aa3170a12", "d86506a2-e186-4b89-97be-3294cb86d53a")
	if err != nil {
		log.Fatalf("error getting people born between: %v", err)
	}
	for _, person := range between {
		log.Println(person)
	}

	log.Println("\nNamesBetween")
	namesBetween, err := df.NamesBetween(context.Background(), "a58933a1-c24f-43d9-bb53-6a1aa3170a12", "d86506a2-e186-4b89-97be-3294cb86d53a")
	if err != nil {
		log.Fatalf("error getting names between: %v", err)
	}
	for _, name := range namesBetween {
		log.Println(name)
	}
}
