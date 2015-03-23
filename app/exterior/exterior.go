package exterior

import (
	"encoding/csv"
	"github.com/aries-auto/appimport/helpers/database"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ExteriorInput struct {
	Make      string
	Model     string
	Style     string
	YearRange string
	Parts     []string
}

type Application struct {
	Year  int
	Make  string
	Model string
	Style string
	Part  string
}

func DoExteriors(filename string) error {
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
		err = IntoDB(a)
	}
	return err
}

//Csv to Struct
func CaptureCsv(filename string) ([]ExteriorInput, error) {
	var e ExteriorInput
	var es []ExteriorInput
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
		e.Make = line[0]
		e.Model = line[1]
		e.Style = line[2]
		e.YearRange = line[3]
		e.Parts = line[4:70]
		es = append(es, e)
	}
	return es, nil
}

//Convert ExteriorInput ot Applications array
func ConvertToApplication(e ExteriorInput) ([]Application, error) {
	var err error
	var app Application
	var apps []Application
	shortYears := strings.Split(e.YearRange, "-")
	var longYear int
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
				if part != "" {
					app.Make = LowerInitial(e.Make)
					app.Model = LowerInitial(e.Model)
					app.Style = e.Style
					app.Year = longYear
					app.Part = part
					apps = append(apps, app)
				}
			}

		}
	}
	return apps, err
}

//Dump into mongo
func IntoDB(app Application) error {
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database.MongoConnectionString().Database).C("exteriors")
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
