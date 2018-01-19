# env

Unmarshal environment variables into structures.

## Usage

Annotate structure properties with the `env:"VAR_NAME"` tag and have  them be populated from environment variables with the same name.

## Methods

```go
func Unmarshal(s interface{})
```

## Example

```go
//main.go
package main

import (
    "fmt"

    "github.com/magento-mcom/go-common/env"
)

// Configuration holder structure
type Configuration struct {
    SettingOne   string `env:"SETTING_ONE"`
    SettingTwo   string `env:"SETTING_TWO"`
    SettingThree string `env:"SETTING_THREE"`
}

func main() {

    config := &Configuration{
        SettingOne:   "default_1",
        SettingTwo:   "default_2",
        SettingThree: "default_3",
    }

    env.Unmarshal(config)

    fmt.Printf("One:   %s\n", config.SettingOne)
    fmt.Printf("Two:   %s\n", config.SettingTwo)
    fmt.Printf("Three: %s\n", config.SettingThree)
}
```

Output:

```sh
$ SETTING_TWO="from environment" go run main.go
One:   default_1
Two:   from environment
Three: default_3
```
