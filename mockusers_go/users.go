package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	CardID  string  `json:"card_id"`
	Email   *string `json:"email,omitempty"`
}

type UserStore struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewUserStore(initial []User) *UserStore {
	s := &UserStore{users: make(map[string]User)}
	for _, u := range initial {
		s.users[u.ID] = u
	}
	return s
}

func (s *UserStore) GetAll() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *UserStore) GetByID(id string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *UserStore) Create(u User) User {
	if u.ID == "" {
		u.ID = randomID()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
	return u
}

func (s *UserStore) Update(id string, u User) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return User{}, false
	}
	u.ID = id
	s.users[id] = u
	return u, true
}

func (s *UserStore) Patch(id string, fields map[string]any) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.users[id]
	if !ok {
		return User{}, false
	}
	// Marshal existing to map, merge, unmarshal back
	data, _ := json.Marshal(existing)
	merged := make(map[string]any)
	json.Unmarshal(data, &merged) //nolint:errcheck
	for k, v := range fields {
		merged[k] = v
	}
	data, _ = json.Marshal(merged)
	var updated User
	json.Unmarshal(data, &updated) //nolint:errcheck
	updated.ID = id
	s.users[id] = updated
	return updated, true
}

func (s *UserStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

func randomID() string {
	b := make([]byte, 8)
	rand.Read(b) //nolint:errcheck
	return fmt.Sprintf("%x", b)
}

// Handlers

func handleGetUsers(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(store.GetAll())
	}
}

func handleGetUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		u, ok := store.GetByID(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

func handleCreateUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		created := store.Create(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}

func handleUpdateUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updated, ok := store.Update(id, u)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(updated)
	}
}

func handlePatchUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var fields map[string]any
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updated, ok := store.Patch(id, fields)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(updated)
	}
}

func handleDeleteUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !store.Delete(id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
