package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	users := loadUsers()
	userStore := NewUserStore(users)
	invoiceStore := NewInvoiceStore()

	mux := http.NewServeMux()

	// JSON middleware wrapping mux
	handler := jsonMiddleware(mux)

	// User routes
	mux.HandleFunc("GET /api/users", handleGetUsers(userStore))

	// Invoice routes
	mux.HandleFunc("POST /api/invoice", handlePostInvoice(invoiceStore))
	mux.HandleFunc("POST /api/invoice-inspector/reset", handleResetInvoices(invoiceStore))
	mux.HandleFunc("GET /api/invoice-inspector/{email}", handleInspectInvoice(invoiceStore))

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("mock-users listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loadUsers() []User {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Resolve db.json relative to this source file
		_, filename, _, _ := runtime.Caller(0)
		dbPath = filepath.Join(filepath.Dir(filename), "..", "db.json")
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		log.Printf("warning: could not read %s: %v — starting with empty user store", dbPath, err)
		return nil
	}

	var db struct {
		Users []User `json:"users"`
	}
	if err := json.Unmarshal(data, &db); err != nil {
		log.Fatalf("failed to parse db.json: %v", err)
	}
	return db.Users
}
