package importer

import (
	"fmt"
	"github.com/aries-auto/appimport/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"log"

	"database/sql"
	"encoding/csv"
	"os"
	"strings"
)

var (
	VehicleApplications map[string]Application
	PartConversion      map[string]int
	Session             *mgo.Session
	inf                 = database.MongoConnectionString()
)

type Input struct {
	Year  string
	Make  string
	Model string
	Style string
	Part  string
}

type Application struct {
	Year  string `bson:"year"`
	Make  string `bson:"make"`
	Model string `bson:"model"`
	Style string `bson:"style"`
	Parts []int  `bson:"parts"`
}

func DoImport(filename string, collectionName string, sess *mgo.Session) error {
	PartConversion = make(map[string]int, 0)
	VehicleApplications = make(map[string]Application, 0)
	Session = sess

	es, err := CaptureCsv(filename)
	if err != nil {
		return err
	}

	for _, e := range es {
		if err := ConvertToApplication(e); err != nil {
			log.Printf("Conversion Error: %s\n", err.Error())
		}
	}

	ClearCollection(collectionName)

	for _, app := range VehicleApplications {
		if err := IntoDB(app, collectionName); err != nil {
			log.Println(err)
		}
	}

	return nil
}

//Csv to Struct
func CaptureCsv(filename string) ([]Input, error) {
	var e Input
	var es []Input
	file, err := os.Open(filename)
	if err != nil {
		return es, err
	}

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		return es, err
	}

	for _, line := range lines {
		if len(line) < 5 {
			log.Println(line)
			continue
		}
		e = Input{
			Make:  strings.ToLower(strings.TrimSpace(line[0])),
			Model: strings.ToLower(strings.TrimSpace(line[1])),
			Style: strings.ToLower(strings.TrimSpace(line[2])),
			Part:  strings.TrimSpace(line[3]),
			Year:  strings.ToLower(strings.TrimSpace(line[4])),
		}

		es = append(es, e)
	}
	return es, nil
}

//Convert Input ot Applications array
func ConvertToApplication(e Input) error {
	var partID int

	if partID = PartConversion[e.Part]; partID == 0 {

		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare("select partID from Part where oldPartNumber = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		if err := stmt.QueryRow(e.Part).Scan(&partID); err != nil || partID == 0 {
			return fmt.Errorf("invalid part: %s", e.Part)
		}

		PartConversion[e.Part] = partID
	}

	tmp := Application{
		Parts: []int{partID},
		Year:  e.Year,
		Make:  e.Make,
		Model: e.Model,
		Style: e.Style,
	}

	idx := VehicleApplications[tmp.string()]
	if idx.Year == "" {
		VehicleApplications[tmp.string()] = tmp
		return nil
	}

	idx.Parts = append(idx.Parts, partID)
	VehicleApplications[tmp.string()] = idx

	return nil
}

func (a *Application) string() string {
	return fmt.Sprintf("%s%s%s%s", a.Year, a.Make, a.Model, a.Style)
}

//Dump into mongo
func IntoDB(app Application, collectionName string) error {
	// session, err := mgo.DialWithInfo(database.MongoConnectionString())
	// if err != nil {
	// 	return err
	// }
	// defer session.Close()

	return Session.DB(database.MongoConnectionString().Database).C(collectionName).Insert(app)
}

func ClearCollection(name string) error {

	// session, err := mgo.DialWithInfo(inf)
	// if err != nil {
	// 	return err
	// }
	// defer session.Close()

	return Session.DB(inf.Database).C(name).DropCollection()
}
