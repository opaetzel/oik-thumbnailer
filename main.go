package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	input := flag.String("input", "", "input folder")
	output := flag.String("output", "", "output folder")
	dimensions := flag.String("dim", "", "dimensions in <width>x<height>")

	useGui := flag.Bool("gui", false, "use gui")

	flag.Parse()

	if *useGui {
		fmt.Println("running with gui. Open browser and direct to \"http://localhost:8080\"")
		runWithGui()
	} else {
		if *input == "" || *output == "" {
			flag.Usage()
			return
		}

		createPackage(input, output, dimensions)

	}
}

func runWithGui() {
	router := NewRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
