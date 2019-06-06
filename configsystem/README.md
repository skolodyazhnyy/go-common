# ConfigSystem
Config system package provides an easy way communicate against orygn. 

## Implemented endpoints

### `Clients`
Gets all configured clients:

```go
client := configsystem.NewClient(
	"https://orygn.local.magento.com",
	"environment",
	&http.Client{},
)

clients, err := client.Clients()
if err != nil {
	panic(err)
}

fmt.Println("List of retrieved clients is", clients)
```


### `Value`
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
