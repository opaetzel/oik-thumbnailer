package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
)

func main() {
	input := flag.String("input", "", "input folder")
	output := flag.String("output", "", "output folder")
	dimensions := flag.String("dim", "", "dimensions in <width>x<height>")

	useGui := flag.Bool("gui", false, "use gui")

	flag.Parse()

	if *useGui {
		fmt.Println("running with gui. Your browser should open")
		go openBrowser()
		runWithGui()
	} else {
		if *input == "" || *output == "" {
			flag.Usage()
			return
		}

		createPackage(input, output, dimensions)

	}
}

func openBrowser() {
	err := browser.OpenURL("http://localhost:8080")
	if err != nil {
		fmt.Println("It seems like your browser did not opent. Start it and direct to \"http://localhost:8080\"")
	}
}

func runWithGui() {
	router := NewRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
