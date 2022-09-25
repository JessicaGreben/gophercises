package handlers

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler returns a http.HandlerFunc that will redirect
// to a different URL for specific endpoints.
func MapHandler(fallback http.Handler) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("map endpoint", r.URL.Path)
		v, ok := pathsToUrls[r.URL.Path]
		if ok {
			fmt.Println("redirecting to", v)
			http.Redirect(w, r, v, http.StatusTemporaryRedirect)
			return
		}
		fmt.Println("fallback to default mux")
		fallback.ServeHTTP(w, r)
	}
}

type YamlData struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

// YAMLHandler parses the provided YAML and will redirect to the URL
// if the path matches from the YAML.
// YAML is expected to be in the format:
//     - path: /some-path
//       url: https://www.some-url.com/demo
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedYaml []YamlData
	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		fmt.Println(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, y := range parsedYaml {
			if r.URL.Path == y.Path {
				fmt.Println("redirecting to", y.Url)
				http.Redirect(w, r, y.Url, http.StatusTemporaryRedirect)
				return
			}
		}
		fmt.Println("fallback to map handler")
		fallback.ServeHTTP(w, r)
	}, err
}
