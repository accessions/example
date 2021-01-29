package main

import (
	_ "embed"
	"log"
	"net/http"
	"github.com/flowchartsman/swaggerui"

)
//go:embed docs/swagger.json
var spec []byte

func main()  {
	log.SetFlags(0)

	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec, swaggerui.SpecTypeJSON)))
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
