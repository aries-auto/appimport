package importer

import (
	"github.com/aries-auto/appimport/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"

	"database/sql"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Input struct {
	Make      string
	Model     string
	Style     string
	YearRange string
	Parts     []string
}

type Application struct {
	Year  int    `bson:"year"`
	Make  string `bson:"make"`
	Model string `bson:"model"`
	Style string `bson:"style"`
	Part  int    `bson:"part"`
}

func DoImport(filename string, collectionName string) error {
	var applications []Application
	es, err := CaptureCsv(filename)
	if err != nil {
		return err
	}

	for _, e := range es {
		apps, err := ConvertToApplication(e)
		if err != nil {
			return err
		}
		applications = append(applications, apps...)
	}
	for _, a := range applications {
		err = IntoDB(a, collectionName)
	}
	return err
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
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		return es, err
	}

	for _, line := range lines {
		e.Make = line[0]
		e.Model = line[1]
		e.Style = line[2]
		e.YearRange = line[3]
		e.Parts = line[4:reader.FieldsPerRecord]
		es = append(es, e)
	}
	return es, nil
}

//Convert Input ot Applications array
func ConvertToApplication(e Input) ([]Application, error) {
	var err error
	var app Application
	var apps []Application
	shortYears := strings.Split(e.YearRange, "-")
	var longYear int

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return apps, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select partID from Part where oldPartNumber = ?")
	if err != nil {
		return apps, err
	}
	defer stmt.Close()

	for _, shortYear := range shortYears {
		if utf8.RuneCountInString(shortYear) == 2 {
			//add 19 or 20
			shortYearInt, err := strconv.Atoi(shortYear)
			if err != nil {
				return apps, err
			}
			if shortYearInt >= 0 && shortYearInt < 16 {
				longYear = 2000 + shortYearInt
			} else {
				longYear = 1900 + shortYearInt
			}

			for _, part := range e.Parts {
				if part != "" && part != "-" && part != "--" {
					//is there an old/new part number?
					var num int
					part = strings.TrimSpace(part)
					err = stmt.QueryRow(part).Scan(&num)
					if err == nil {
						app.Part = num
					} else {
						app.Part, err = strconv.Atoi(part)
						if err != nil {
							//non-existent part
							continue
						}
					}
					app.Make = LowerInitial(e.Make)
					app.Model = LowerInitial(e.Model)
					app.Style = e.Style
					app.Year = longYear
					apps = append(apps, app)
				}
			}

		}
	}
	return apps, err
}

//Dump into mongo
func IntoDB(app Application, collectionName string) error {
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database.MongoConnectionString().Database).C(collectionName)
	err = c.Insert(app)
	return err

}

//Proper capitalization
func LowerInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + strings.ToLower(str[i+1:])
	}
	return ""
}
