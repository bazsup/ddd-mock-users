package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

//go:embed db.json
var dbJSON []byte

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
	var db struct {
		Users []User `json:"users"`
	}
	if err := json.Unmarshal(dbJSON, &db); err != nil {
		log.Fatalf("failed to parse embedded db.json: %v", err)
	}
	return db.Users
}
