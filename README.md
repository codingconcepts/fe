# fe
Export functions out of databases and into code

### Usage

``` sh
$ fe
Extract functions from databases into code

Usage:
  fe [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  postgres    extract from postgres
  version     Show the application version

Flags:
      --go-package string   package name of the output Go code
  -h, --help                help for fe
  -o, --output string       absolute or relative path to the output file
  -u, --url string          full database url/connection string

Use "fe [command] --help" for more information about a command.
```

So, for example, to create Go code from the functions in a Postgres database:
``` sh
$ fe postgres go \
  -u "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" \
  --go-package repo \
  -o examples/postgres/out.go
```

### Todos

* MySQL function extraction
* Oracle function extraction
* SQL Server function extraction
* Sybase function extraction