package rest

import (
	"github.com/100DAYS/tablegateway"
)

type DBResource struct {
	IsNew        bool
	Id           int64
	ResourceName string
	Payload      interface{}
}

func (r DBResource) HasId() bool {
	return !r.IsNew
}

func (r DBResource) GetId() int64 {
	return r.Id
}

func (r DBResource) GetResourceName() string {
	return r.ResourceName
}

func (r DBResource) GetPayload() interface{} {
	return r.Payload
}

func (r DBResource) GetLinks() LinkList {
	return nil
}

type DBResourceDriver struct {
	Gw           tablegateway.AutomatableTableDataGateway
	ResourceName string
}

func (dbr DBResourceDriver) Find(id int64) (res Resource, err error) {
	s := dbr.Gw.GetStruct()
	err = dbr.Gw.Find(id, s)
	if err != nil {
		return
	}
	res = DBResource{Payload: s, IsNew: false, Id: id, ResourceName: dbr.ResourceName}
	return
}

func (dbr DBResourceDriver) Delete(id int64) (int64, error) {
	return dbr.Gw.Delete(id)
}

func (dbr DBResourceDriver) Create(data Resource) (int64, error) {
	return dbr.Gw.Insert(data.GetPayload())
}

func (dbr DBResourceDriver) Update(id int64, changes map[string]interface{}) (int64, error) {
	return dbr.Update(id, changes)
}

func (dbr DBResourceDriver) Query(q Query) (ListResult, error) {
	res := dbr.Gw.GetStructList()
	err := dbr.Gw.FilterQuery(q.Filters, q.Order, q.Offset, q.Limit, res)
	if err != nil {
		return ListResult{}, err
	}
	resList := make([]Resource, len(*res))
	i := 0
	for rec := range *res {
		id, err := dbr.Gw.GetId(rec)
		if err != nil {
			return ListResult{}, err
		}
		resList[i] = DBResource{Payload: rec, IsNew: false, ResourceName: dbr.ResourceName, Id: id}
		i++
	}
	return ListResult{ResourceName: dbr.ResourceName, Query: q, List: resList}, nil
}
