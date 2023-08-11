package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/codingconcepts/fe/internal/pkg/code"
	"github.com/codingconcepts/fe/internal/pkg/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var (
	version string

	databases = []string{"postgres"}
	languages = []string{"go"}

	flagURL             string
	flagOutputFile      string
	flagGoOutputPackage string
)

func main() {
	log.SetFlags(0)

	rootCmd := &cobra.Command{
		Use:   "fe",
		Short: "Extract functions from databases into code",
	}
	rootCmd.PersistentFlags().StringVarP(&flagURL, "url", "u", "", "full database url/connection string")
	rootCmd.PersistentFlags().StringVarP(&flagOutputFile, "output", "o", "", "absolute or relative path to the output file")
	rootCmd.PersistentFlags().StringVar(&flagGoOutputPackage, "go-package", "", "package name of the output Go code")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("output")

	databaseCmds := databaseCommands(databases, languages)
	rootCmd.AddCommand(databaseCmds...)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func databaseCommands(databases []string, languages []string) []*cobra.Command {
	var databaseCmds []*cobra.Command

	for _, database := range databases {
		databaseCmd := &cobra.Command{
			Use:   database,
			Short: fmt.Sprintf("extract from %s", database),
		}

		for _, language := range languages {
			languageCmd := &cobra.Command{
				Use:   language,
				Short: fmt.Sprintf("Generate %s code from extracted functions", language),
				Run:   runCommand(database, language),
			}

			databaseCmd.AddCommand(languageCmd)
		}

		databaseCmds = append(databaseCmds, databaseCmd)
	}

	return databaseCmds
}

func runCommand(database, language string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// TODO: Implement more repos.
		db, err := pgxpool.New(context.Background(), flagURL)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}
		defer db.Close()

		repo := repo.NewPostgresRepo(db)
		functions, err := repo.GetFunctions()
		if err != nil {
			log.Fatalf("error getting functions: %v", err)
		}

		// TODO: Implement more languages.
		cg, err := code.NewGoCodeGenerator()
		if err != nil {
			log.Fatalf("error creating code generator: %v", err)
		}

		file, err := os.Create(flagOutputFile)
		if err != nil {
			log.Fatalf("error creating output file: %v", err)
		}

		if err = cg.Generate(functions, file, flagGoOutputPackage); err != nil {
			log.Fatalf("error generating function file: %v", err)
		}
	}
}
