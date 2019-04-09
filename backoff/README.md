# Backoff package

Backoff is a package which helps to handle the exponentially increase of delay for proceses 

## Usage

Just instantiate new backoff with `backoff.New()` and call `.Timeout()` everytime you need to encrease delay.
For custom backoff time periods you can use `backoff.NewWithParameters()` in which you can specify the maximum time delay and after how long it has to be reset

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/magento-mcom/go-common/backoff"
)

func main() {
	boff := backoff.New()

	for {
		t := boff.Timeout()

		select {
		case <-time.After(t):
			fmt.Println("Timeout increase for next iteration: ", t)
		}
	}
}

```

Output:

```sh
$ ~/home » go run main.go
------------------------------------------------
Timeout increase for next iteration:  1s
Timeout increase for next iteration:  2s
Timeout increase for next iteration:  4s
Timeout increase for next iteration:  8s
Timeout increase for next iteration:  16s

```
