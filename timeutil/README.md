# timeutil

> Utilities around time with respect of `TZ` environment variable

## Install

```shell
go get github.com/avakarev/go-util/timeutil
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/avakarev/go-util/timeutil"
)

func respectTimezone(t time.Time) {
	// `TZ` env var is set to "Europe/Berlin"
	fmt.Println(t.Format(time.RFC3339))                 // => 2022-06-03T16:26:15Z
	fmt.Println(timeutil.Local(t).Format(time.RFC3339)) // => 2022-06-03T18:26:15+02:00
}

func mockTimeNow(t time.Time) {
	timeutil.MockNow(func() time.Time {
		return t
	})
	defer timeutil.UnmockNow()
	fmt.Println(timeutil.Now().Format(time.RFC3339)) // => 2022-06-03T16:26:15Z
}

func main() {
	t, _ := time.Parse(time.RFC3339, "2022-06-03T16:26:15Z")

	respectTimezone(t)
	mockTimeNow(t)
}
```


## License

`go-timeutil` is licensed under MIT license. (see [LICENSE](./../LICENSE))
