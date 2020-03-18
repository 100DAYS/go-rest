package rest

import (
	"reflect"

	"github.com/100DAYS/tablegateway"
)

type DBResource struct {
	IsNew        bool `json:"-"`
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
	return dbr.Gw.Update(id, changes)
}

func (dbr DBResourceDriver) Query(q Query) (ListResult, error) {
	res := dbr.Gw.GetStructList()
	err := dbr.Gw.FilterQuery(q.Filters, q.Order, q.Offset, q.Limit, res)
	if err != nil {
		return ListResult{}, err
	}
	// walking through lists: https://gist.github.com/ahmdrz/5448c37063d3757b5581fb8df53800e0
	val := reflect.ValueOf(res)
	if val.Kind() != reflect.Ptr {
		panic("Result was no Pointer")
	}
	rows := val.Elem()
	if rows.Kind() != reflect.Slice {
		panic("Value was no slice")
	}
	numrows := rows.Len()
	resList := make([]Resource, numrows)
	for i := 0; i < numrows; i++ {
		rec := rows.Index(i)
		if rec.Kind() != reflect.Struct {
			panic("Rows must be Structs!")
		}
		recIf := rec.Interface()
		id, err := dbr.Gw.GetId(recIf)
		if err != nil {
			return ListResult{}, err
		}
		resList[i] = DBResource{Payload: recIf, IsNew: false, ResourceName: dbr.ResourceName, Id: id}
	}
	return ListResult{ResourceName: dbr.ResourceName, Query: q, List: resList}, nil
}

func (dbr DBResourceDriver) NewResource(payload interface{}) DBResource {
	return DBResource{IsNew: true, Id: 0, ResourceName: dbr.ResourceName, Payload: payload}
}

// @todo: make sure values of type sql.NullXXX are rendered correctly in JSON response
