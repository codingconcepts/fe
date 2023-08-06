package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	version string

	databases = []string{"postgres", "cockroachdb", "oracle", "sqlserver", "mysql"}
	languages = []string{"go"}
)

func main() {
	log.SetFlags(0)

	rootCmd := &cobra.Command{
		Use:   "fe",
		Short: "Extract functions from databases into code",
	}

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

	}
}
