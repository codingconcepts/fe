# fe
Export functions out of databases and into code

### Workflow

``` sh
$ fe [database] [language] args
```

So, for example, to create Go code from the functions in a Postgres database:
``` sh
$ fe postgres go --url 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable'
```

### Todos

* Support for SETOF scalar values (e.g. `SETOF VARCHAR`)
* Run gofmt before saving output file