package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
)

type URLShortener struct {
	db map[string]*string
	mu sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		db: make(map[string]*string, 1024),
	}
}

func (us *URLShortener) store(shortURL, longURL string) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.db[shortURL] = &longURL
}

func (us *URLShortener) get(shortURL string) (*string, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	longURL, exists := us.db[shortURL]
	return longURL, exists
}

func (us *URLShortener) redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURLPtr, exists := us.get(shortURL)
	if !exists || longURLPtr == nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, *longURLPtr, http.StatusFound)
}

func (us *URLShortener) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortURL := shortenedURL()
	us.store(shortURL, longURL)

	fmt.Fprintf(w, "shortened URL: http://localhost:8080/%s", shortURL)
}

func shortenedURL() string {
	const chars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASFGHJKLZXCVBNM0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[rand.IntN(len(chars))]
	}
	return string(result)
}

func main() {
	us := NewURLShortener()

	http.HandleFunc("/", us.redirectHandler)
	http.HandleFunc("/shorten", us.shortenHandler)

	port := ":8080"
	fmt.Println("local")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("failed to listen and serve on %s", port)
	}
}
