package main

import (
	"github.com/aries-auto/appimport/app/importer"
	"log"
)

func main() {

	log.Println(importer.DoImport("csvs/grille_guards.csv", "grille guards"))
}
