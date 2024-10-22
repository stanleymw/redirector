package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

func getPaste(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	dat, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	ret := string(dat)

	log.Printf("Enlace nuevo: %s", ret)
	return ret, nil
}

func main() {
	target := flag.String("url", "https://pastebin.com/raw/hE40FhFF", "URL of which to obtain the URL to redirect to")
	endpoint := flag.String("endpoint", "/redirect", "Redirect endpoint")
	addr := flag.String("addr", ":8080", "Address to start server on")
	cooldown := flag.Int64("lifetime", 300, "Lifetime to cache url (in seconds)")
	flag.Parse()

	last_cached_at := time.Now()
	url_cached, err := getPaste(*target)

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if time.Now().After(last_cached_at.Add(time.Duration(*cooldown) * time.Second)) {
			url_cached, err = getPaste(*target)
		}

		if err != nil {
			http.Error(w, "Could not access redirect link!", http.StatusInternalServerError)
		}
		http.Redirect(w, &http.Request{}, url_cached, http.StatusSeeOther)
	}

	http.HandleFunc(*endpoint, helloHandler)
	log.Println("Empezando...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
