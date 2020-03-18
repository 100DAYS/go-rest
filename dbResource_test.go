package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/100DAYS/tablegateway"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func prepareDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	//db, err := sqlx.Connect("postgres", "host=localhost user=example password=example dbname=example sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Error opening sqlx: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging sqlx: %s", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS places")
	if err != nil {
		return nil, fmt.Errorf("Error dropping Table: %s", err)
	}

	var ddl string
	if db.DriverName() == "postgres" {
		ddl = `CREATE TABLE places (
                id SERIAL PRIMARY KEY,
                country text,
                city text NULL,
                telcode integer);`
	} else {
		ddl = `CREATE TABLE places (
                id INTEGER PRIMARY KEY AUTOINCREMENT ,
                country text,
                city text NULL,
                telcode integer);`
	}

	_, err = db.Exec(ddl)
	if err != nil {
		return nil, fmt.Errorf("Error creating Schema: %s", err)
	}

	return db, err
}

type Place struct {
	Id      sql.NullInt64  `db:"id"`
	Country string         `db:"country"`
	City    sql.NullString `db:"city"`
	Telcode int            `db:"telcode"`
}

type PlaceGw struct {
	tablegateway.TableGateway
}

func NewPlaceGw(db *sqlx.DB) PlaceGw {
	return PlaceGw{tablegateway.NewGw(db, "places", "id")}
}

func (pg PlaceGw) GetStruct() interface{} {
	return &Place{}
}
func (pg PlaceGw) GetStructList() interface{} {
	return &[]Place{}
}

func TestGeneralized(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Errorf("Cannot open DB: %s", err)
		return
	}

	dao := NewPlaceGw(db)
	r := DBResourceDriver{Gw: &dao, ResourceName: "places"}

	dbr := r.NewResource(Place{Id: tablegateway.NullInt64(12), Country: "Germany"})
	id, err := r.Create(dbr)
	fmt.Printf("Inserted with Id %d\n", id)

	dbr2, err := r.Find(12)
	if err != nil {
		t.Errorf("Error Finding record: %s", err)
	}
	fmt.Printf("Record: %#v\n", dbr2)

	affected, err := r.Update(12, map[string]interface{}{"City": "Frankfurt"})
	if err != nil {
		t.Errorf("Error Updating record: %s", err)
	}
	if affected != 1 {
		t.Errorf("Affected Rows: 1 expected, %d found", affected)
	}

	dbr3, err := r.Find(12)
	if err != nil {
		t.Errorf("Error Finding record: %s", err)
	}
	pl := dbr3.GetPayload().(*Place)
	if pl.City.String != "Frankfurt" {
		t.Errorf("Frankfurt expected, found %s", pl.City.String)
	}
	fmt.Printf("Payload: %#v\n", pl)

	affected, err = r.Delete(12)
	if err != nil {
		t.Errorf("Error Deleting record: %s", err)
	}
	if affected != 1 {
		t.Errorf("Affected Rows: 1 expected, %d found", affected)
	}
}

func TestFiltering(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Errorf("Cannot open DB: %s", err)
		return
	}

	dao := NewPlaceGw(db)
	r := DBResourceDriver{Gw: &dao, ResourceName: "places"}

	dbr := r.NewResource(Place{Country: "Germany", City: tablegateway.NullString("Frankfurt"), Telcode: 1})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Germany", City: tablegateway.NullString("Stuttgart")})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Germany", City: tablegateway.NullString("München")})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Germany", City: tablegateway.NullString("Ulm"), Telcode: 1})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Germany", City: tablegateway.NullString("Hamburg")})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Italia", City: tablegateway.NullString("Roma")})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "Italia", City: tablegateway.NullString("Milano")})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "USA", City: tablegateway.NullString("New York"), Telcode: 1})
	_, err = r.Create(dbr)
	dbr = r.NewResource(Place{Country: "USA", City: tablegateway.NullString("San Fran")})
	_, err = r.Create(dbr)

	q := Query{Filters: map[string]interface{}{"telcode": 1}, Limit: 100, Order: []string{"country"}}
	fmt.Printf("Query: %#v\n", q)
	res, err := r.Query(q)
	if err != nil {
		t.Errorf("Error during FilterQuery: %s", err)
		return
	}
	if len(res.List) != 3 {
		t.Errorf("Expected 3 Results, found %d", len(res.List))
		return
	}
	if res.List[0].GetId() != 1 {
		t.Errorf("First result Id should be 1,  %d found", res.List[0].GetId())
		return
	}
	pl, err := json.Marshal(res)
	if err != nil {
		t.Errorf("Error Marshalling Response: %s", err)
		return
	}
	fmt.Printf("Res: \n%s\n", pl)
}
