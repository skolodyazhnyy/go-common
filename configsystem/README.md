# ConfigSystem

Config system package provides an easy way communicate against orygn. 

## Implemented endpoints

### Value

Gets a single node/value:

```go
client := configsystem.NewClient(
	"https://orygn.local.magento.com",
	"environment",
	&http.Client{},
)

var value string

err := client.Value("LUMA", "scope", "key", &value)
if err != nil {
	panic(err)
}

fmt.Println("Value retrieved is", value)
```
