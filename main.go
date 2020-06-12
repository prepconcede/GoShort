package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"pivdoggo/urlShortener/urlshort"
)

func main() {
	yamlFilename := flag.String("yaml-file", "redirect.yaml", "Yaml file name with redirection URLs")
	flag.Parse()

	mux := defaultMux()

	mapHandler := urlshort.NewHttpRedirectHandler(
		urlshort.NewBaseUrlMapper(map[string]string{
			"/pivdoggo-vk":        "https://vk.com/zubaastik",
			"/pivdoggo-shortener": "https://github.com/prepconcede/GoShort",
		}), mux)
	yamlUrlMapper, err := urlshort.NewYamlUrlMapper(*yamlFilename)
	if err != nil {
		log.Fatalf("Cant create YAML redirect URL provider. %v", err)
	}

	yamlHandler := urlshort.NewHttpRedirectHandler(yamlUrlMapper, mapHandler)

	fmt.Println("starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
