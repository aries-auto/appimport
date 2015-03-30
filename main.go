package main

import (
	"github.com/aries-auto/appimport/app/importer"
	"log"
)

func main() {

	// err := importer.DoImport("csvs/grille_guards.csv", "grille guards")
	// if err != nil {
	// 	log.Printf("Errored on Grille Guard: %s\n", err.Error())
	// } else {
	// 	log.Println("Finished Grille Guard")
	// }

	err := importer.DoImport("csvs/bull_bar.csv", "bull bar")
	if err != nil {
		log.Printf("Errored on Bull Bar: %s\n", err.Error())
	} else {
		log.Println("Finished Bull Bar")
	}

}
