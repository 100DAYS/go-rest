package rest

type LinkList map[string]string

type Resource interface {
	GetId() int64
	HasId() bool
	GetResourceName() string
	GetLinks() LinkList
	GetPayload() interface{}
}

type ListResult struct {
	ResourceName string
	Query        Query
	List         []Resource
}

type ResourceDriver interface {
	Find(id int64) (Resource, error)
	Delete(id int64) error
	Create(data Resource) (int64, error)
	Update(id int64, changes map[string]interface{}) error
	Query(q Query) (ListResult, error)
}
