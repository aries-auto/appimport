package main

import (
	"github.com/aries-auto/appimport/importer"
	"log"
)

var (
	DataFiles = map[string]string{
		"4 in bull bars, big horn":                  "csvs/four_inch_big_horn_bull_bars.csv",
		"4 in round side bars, big step":            "csvs/four_inch_big_step_side_bars.csv",
		"4 in oval side bars":                       "csvs/four_inch_oval_side_bars.csv",
		"grille guards":                             "csvs/grille_guards.csv",
		"pro series grille guards":                  "csvs/pro_series_grille_guards.csv",
		"3 in bull bars":                            "csvs/three_inch_bull_bars.csv",
		"3 in bull bars, pro series":                "csvs/three_inch_pro_series_bull_bars.csv",
		"3 in round side bars, pro series":          "csvs/three_inch_pro_series_side_bars.csv",
		"3 in round side bars":                      "csvs/three_inch_side_bars.csv",
		"6 in oval side bars and mounting brackets": "csvs/six_inch_oval_side_bars.csv",
	}
)

func main() {

	for name, file := range DataFiles {
		err := importer.DoImport(file, name)
		if err != nil {
			log.Printf("Errored on %s: %s\n", name, err.Error())
		} else {
			log.Printf("Finished %s", name)
		}
	}

}
