package main

import (
	"github.com/aries-auto/appimport/helpers/database"
	"github.com/aries-auto/appimport/importer"
	"gopkg.in/mgo.v2"
	"log"
)

var (
	DataFiles = map[string]string{
		// "4 in round side bars, big step": "csvs/four_in_big_step_side_bars.csv",
		// "4 in oval side bars":            "csvs/four_in_oval_side_bars.csv",
		// "4 in oval side bars, wheel to wheel":       "csvs/four_in_oval_side_bars_wheel_to_wheel.csv",
		"grille guards": "csvs/grille_guards.csv",
		"bull bars":     "csvs/bull_bars.csv",
		// "3 in round side bars, pro series":          "csvs/three_in_pro_series_side_bars.csv",
		// "3 in round side bars":                      "csvs/three_in_side_bars.csv",
		// "6 in oval side bars and mounting brackets": "csvs/six_in_oval_side_bars_and_mounting_brackets.csv",
		// "jeep bumper kits and replacement parts":    "csvs/jeep_bumper_kits_and_replacement_parts.csv",
		// "jeep accessories":                          "csvs/jeep_accessories.csv",
		// "interior": "csvs/interior.csv",
	}
)

func main() {

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		panic(err)
	}

	for name, file := range DataFiles {
		tmp := sess.Clone()
		err := importer.DoImport(file, name, tmp)
		if err != nil {
			log.Printf("Errored on %s: %s\n", name, err.Error())
		} else {
			log.Printf("Finished %s", name)
		}
		tmp.Close()
	}

}
