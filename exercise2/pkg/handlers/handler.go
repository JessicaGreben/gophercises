package handlers

import (
	"log"
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
		log.Println("map endpoint", r.URL.Path)
		v, ok := pathsToUrls[r.URL.Path]
		if ok {
			log.Println("redirecting to", v)
			http.Redirect(w, r, v, http.StatusTemporaryRedirect)
			return
		}
		log.Println("fallback to default mux")
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler parses the provided YAML and will redirect to
// if defined in the YAML.
// YAML is expected to be in the format:
//     - path: /some-path
//       url: https://www.some-url.com/demo
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		y := parseYaml(yml)
		for _, x := range y {
			if r.URL.Path == x.Path {
				log.Println("redirecting to", x.Url)
				http.Redirect(w, r, x.Url, http.StatusTemporaryRedirect)
				return
			}
		}
		log.Println("fallback to map handler")
		fallback.ServeHTTP(w, r)
	}, nil
}

type YamlStruct struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

func parseYaml(yml []byte) []YamlStruct {
	var y []YamlStruct
	err := yaml.Unmarshal(yml, &y)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return y
}
