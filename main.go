///var/lib/webserver/survey_results.txt
// A főbb részek maradnak ugyanazok
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
    "cloud.google.com/go/storage"
	"context"
	"io"
	"time"

)

var (
	visitorCount int
	mu           sync.Mutex
)

// Írás Cloud Storage bucketbe
func writeToBucket(bucketName, objectName, content string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.WriteString(wc, content); err != nil {
		return fmt.Errorf("write: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("close: %v", err)
	}
	return nil
}


func countMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		visitorCount++
		log.Printf("New visitor: %s", r.RemoteAddr)
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
func appendToBucketFile(bucketName, objectName, newLine string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("client: %v", err)
	}
	defer client.Close()

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(objectName)

	// 1. Meglévő fájl tartalmának beolvasása (ha létezik)
	var currentContent string
	reader, err := obj.NewReader(ctx)
	if err == nil {
		defer reader.Close()
		data, err := io.ReadAll(reader)
		if err == nil {
			currentContent = string(data)
		}
	}

	// 2. Új tartalom összeállítása
	newContent := currentContent + newLine

	// 3. Fájl felülírása az új tartalommal
	writer := obj.NewWriter(ctx)
	_, err = writer.Write([]byte(newContent))
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}
	return writer.Close()
}

// Végpont, ahová a felmérés adatai érkeznek
func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	ageGroup := r.FormValue("ageGroup")
	log.Printf("Felmérés válasz érkezett: %s", ageGroup)

	line := fmt.Sprintf("%s - %s\n", r.RemoteAddr, ageGroup)
	err := appendToBucketFile("webserver-data", "/var/lib/webserver/valaszok.txt", line)
	if err != nil {
		log.Printf("Hiba írás közben: %v", err)
		http.Error(w, "Hiba a mentés közben", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Köszönjük a választ!")
}

func main() {

	// Statikus fájlok
	fs := http.FileServer(http.Dir("/static"))
	http.Handle("/", countMiddleware(fs))

	// Felmérés beküldésének kezelése
	http.HandleFunc("/submit", submitHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

