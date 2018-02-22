package structs

type Scope struct {
	Client Client
	Id       string        `json:"id"`
	Children []Scope `json:"children"`
}