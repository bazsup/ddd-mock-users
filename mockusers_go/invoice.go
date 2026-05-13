package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type InvoiceEntry struct {
	Body       map[string]any
	ReceivedAt time.Time
}

type InvoiceInspectorResponse struct {
	Email       string     `json:"email"`
	Status      string     `json:"status"`
	FulfilledAt *time.Time `json:"fulfilledAt"`
	Body        any        `json:"body"`
}

type InvoiceStore struct {
	mu      sync.RWMutex
	entries []InvoiceEntry
}

func NewInvoiceStore() *InvoiceStore {
	return &InvoiceStore{}
}

func (s *InvoiceStore) Add(body map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, InvoiceEntry{Body: body, ReceivedAt: time.Now()})
}

func (s *InvoiceStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = nil
}

func (s *InvoiceStore) FindByEmail(email string, timeout time.Duration) InvoiceInspectorResponse {
	deadline := time.Now().Add(timeout)
	for {
		s.mu.RLock()
		for i := len(s.entries) - 1; i >= 0; i-- {
			e := s.entries[i]
			if v, ok := e.Body["email"]; ok && v == email {
				s.mu.RUnlock()
				t := e.ReceivedAt
				return InvoiceInspectorResponse{
					Email:       email,
					Status:      "fulfilled",
					FulfilledAt: &t,
					Body:        e.Body,
				}
			}
		}
		s.mu.RUnlock()

		if time.Now().After(deadline) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	return InvoiceInspectorResponse{Email: email, Status: "waiting"}
}

// Handlers

func handlePostInvoice(store *InvoiceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		store.Add(body)
		json.NewEncoder(w).Encode(map[string]any{"status": "OK", "body": body})
	}
}

func handleResetInvoices(store *InvoiceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		json.NewEncoder(w).Encode(map[string]string{"status": "reset"})
	}
}

func handleInspectInvoice(store *InvoiceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.PathValue("email")
		timeoutMs, _ := strconv.ParseInt(r.URL.Query().Get("timeout"), 10, 64)
		timeout := time.Duration(timeoutMs) * time.Millisecond
		result := store.FindByEmail(email, timeout)
		json.NewEncoder(w).Encode(result)
	}
}
