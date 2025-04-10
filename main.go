// A főbb részek maradnak ugyanazok
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	visitorCount int
	mu           sync.Mutex
)

func countMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		visitorCount++
		log.Printf("New visitor: %s", r.RemoteAddr)
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// Végpont, ahová a felmérés adatai érkeznek
func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Beolvasás
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	ageGroup := r.FormValue("ageGroup")

	// Logolás vagy mentés
	log.Printf("Felmérés válasz érkezett: %s", ageGroup)

	// (opcionálisan fájlba mentés)
	f, err := os.OpenFile("survey_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		fmt.Fprintln(f, ageGroup)
	}

	// Válasz vissza a kliensnek
	fmt.Fprintf(w, "Köszönjük a választ!")
}

func main() {
	// Statikus fájlok
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", countMiddleware(fs))

	// Felmérés beküldésének kezelése
	http.HandleFunc("/submit", submitHandler)

	port := ":8000"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

