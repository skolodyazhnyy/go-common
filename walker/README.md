# walker

Recursively walk structures applying a defined callback via reflection.

## Methods

```go
func Walk(v interface{}, f callback)
```

## Example

```go
//main.go
package main

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/magento-mcom/go-common/walker"
)

// Object holder structure
type Object struct {
    PropertyOne   string
    PropertyTwo   int
    PropertyThree Nested
}

type Nested struct {
    DeepOne string
}

func main() {

    obj := &Object{
        PropertyOne:   "one",
        PropertyTwo:   2,
        PropertyThree: Nested{DeepOne: "deep value"},
    }

    walker.Walk(obj, func(v interface{}, sf reflect.StructField) {
        kind := sf.Type.Kind()
        val := v.(*reflect.Value)

        // Find which properties are strings, and upper case them
        if kind == reflect.String {
            val.SetString(strings.ToUpper(val.String()))
        }

    })

    fmt.Println(obj.PropertyOne)
    fmt.Println(obj.PropertyTwo)
    fmt.Println(obj.PropertyThree.DeepOne)
}
```

Output:

```sh
$ go run main.go
ONE
2
DEEP VALUE
```
