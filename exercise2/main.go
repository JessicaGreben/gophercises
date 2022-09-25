package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jessicagreben/gophercises/exercise2/pkg/handlers"
)

var yamlFileName = flag.String("file", "example.yaml", "path to the yaml file containing path and url redirects")

func main() {
	flag.Parse()

	mux := defaultMux()
	mapHandler := handlers.MapHandler(mux)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("os.Getwd:", err)
	}
	yml, err := ioutil.ReadFile(filepath.Join(dir, *yamlFileName))
	if err != nil {
		fmt.Println(err)
	}
	yamlHandler, err := handlers.YAMLHandler(yml, mapHandler)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
