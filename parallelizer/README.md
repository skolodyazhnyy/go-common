# Parallelizer package

Parallelizer helps to limit and control how many parallel jobs you have. After reached the limit, next jobs will wait for the previous one to finish.  

## Usage

Just instantiate new parallelizer with `parallelizer.New()` specifying the number of jobs that you want to limit to. Don't forget to use `.Wait()` to make the process wait until the last job finishes. 

## Example
Here we are going to execute (inifitly) a parallel task with 2 concurrent threads so the third one will wait until the previous ones finish:

```go
package main

import (
	"fmt"
	"time"

	"github.com/magento-mcom/go-common/parallelizer"
)

func main() {
	numberOfParallelJobs := 2

	prlzr := parallelizer.New(numberOfParallelJobs)
	defer prlzr.Wait()

	for {
		prlzr.Execute(func() {
			fmt.Println("Executing in parallel")

			time.Sleep(5 * time.Second)
		})
	}
}
```