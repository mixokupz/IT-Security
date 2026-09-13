package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imageURL := r.URL.Query().Get("url")
	if imageURL == "" {
		http.Error(w, "missing 'url' parameter", http.StatusBadRequest)
		return
	}

	// Скачиваем картинку по URL
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, imageURL, nil)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Обязательно: описательный User-Agent с контактом
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("User-Agent", "GO")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to fetch image: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "remote returned "+resp.Status, http.StatusBadGateway)
		return
	}

	// Определяем MIME-тип
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = guessMimeFromURL(imageURL)
	}

	// Отдаём картинку клиенту
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, resp.Body); err != nil {
		// уже начали писать ответ — логгируем, но не пытаемся http.Error
		fmt.Println("copy error:", err)
	}
}

func guessMimeFromURL(u string) string {
	lower := strings.ToLower(u)
	switch {
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".webp"):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("User-Agent") == "GO" {
		fmt.Fprintln(w, "SECRETS")
	}
	if r.Header.Get("User-Agent") != "GO" {
		fmt.Fprintln(w, "NO SECRETS")
	}
}

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	dbPass = "secret" // ваш пароль
	dbName = "mydb"   // имя базы
)

func main() {

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/secrets", InfoHandler)
		log.Println("internal on localhost:8081")
		log.Fatal(http.ListenAndServe("localhost:8081", mux))
	}()

	/*connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	defer sqlDB.Close()
	db := database.NewDatabase(sqlDB)
	ctrl := controller.NewProductController(db)
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", UploadHandler)
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/info", InfoHandler)
	http.ListenAndServe("localhost:8080", mux)
}
