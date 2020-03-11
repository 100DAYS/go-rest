package rest

import "fmt"

type Link interface {
	Render() string
}

type ResourceLink struct {
	ResourceName string
	Id           int
}

func (rl *ResourceLink) Render() string {
	return fmt.Sprintf("/%s/%d", rl.ResourceName, rl.Id)
}

type CollectionLink struct {
	ResourceName string
	Query        Query
}

func (cl *CollectionLink) Render() string {
	return fmt.Sprintf("/%s?%s", cl.ResourceName, cl.Query.Render())
}
