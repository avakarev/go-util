# gormutil

> Opinionated utilities for Gorm

## Install

```shell
go get github.com/avakarev/go-util/gormutil
```

## Usage
```go
package main

import (
    "log"

    "github.com/nutsdb/nutsdb"
)

func main() {
    dsn := "./data/sqlite.db"
    if fd := os.Getenv("MYAPP_DSN"); fd != "" {
        dsn = fd
    }
    sqlitedb, err := gormutil.Open(
        sqlite.Open(dsn),
        gormutil.WithLogger(zerologger.NewDefault()),
        gormutil.WithLocks(),
    )
    if err != nil {
        return err
    }
    db = sqlitedb
}

```

## License

`go-testutil` is licensed under MIT license. (see [LICENSE](./../LICENSE))
