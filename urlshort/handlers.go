package urlshort

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

//this functions creates and returns a mapping of urls
func NewBaseUrlMapper(urls map[string]string) func(string) (string, bool) {
	return func(path string) (string, bool) {
		url, ok := urls[path]
		return url, ok
	}
}

//this function parses yaml,  and returns a base mapper with its content
func NewYamlUrlMapper(filename string) (func(string) (string, bool), error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	yml := []map[string]string{}
	err = yaml.Unmarshal(content, &yml)

	if err != nil {
		return nil, err
	}

	mapping := make(map[string]string)

	for _, m := range yml {
		mapping[m["path"]] = m["url"]
	}

	return NewBaseUrlMapper(mapping), nil
}

//this is the function that makes redirects
func NewHttpRedirectHandler(mapper func(string) (string, bool), fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := mapper(r.URL.Path); ok {
			log.Printf("Redirecting %s to %s\n", r.URL.Path, url)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
