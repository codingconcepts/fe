# fe
Export functions out of databases and into code

### Workflow

``` sh
$ fe [database] [language] args
```

So, for example, to create Go code from the functions in a Postgres database:
``` sh
$ fe postgres go \
  -u "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" \
  --go-package repo \
  -o examples/postgres/out.go
```

### Todos

* Support delete functions