package main

import (
	"github.com/aries-auto/appimport/app/importer"
	"log"
)

func main() {
	var err error
	// err = importer.DoImport("csvs/exteriorsCompleteApp.csv", "exteriors")
	// log.Print(err)

	err = importer.DoImport("csvs/interiorsWIPfix.csv", "interiors")
	log.Print(err)
}
