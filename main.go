package main

import (
	"github.com/aries-auto/appimport/app/exterior"
	"log"
)

func main() {
	err := exterior.DoExteriors("csvs/exteriorsCompleteApp.csv")
	log.Print(err)
}
