# buildmeta

> Tiny go module to hold build and runtime information

## Install

```shell
go get github.com/avakarev/go-util/buildmeta
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/avakarev/go-util/buildmeta"
)

func main() {
	fmt.Println(buildmeta.Compiler()) // => go1.18.3
}
```


## License

`go-buildmeta` is licensed under MIT license. (see [LICENSE](./../LICENSE))
