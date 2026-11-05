package main

import (
	"backend/internal/assistant"
	"encoding/json"
	"log"
	"net/http"
)

type MessageRequest struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Reply string `json:"reply"`
}

func main() {
	ai, err := assistant.New()
	if err != nil {
		log.Fatal(err)
	}
	
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../frontend")))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"app":    "noah",
		})
	})

	mux.HandleFunc("POST /api/chat", func(w http.ResponseWriter, r *http.Request) {
		var req MessageRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if req.Message == "" {
			http.Error(w, `{"error":"message is required"}`, http.StatusBadRequest)
			return
		}

		reply, err := ai.Reply(r.Context(), req.Message)
		if err != nil {
			log.Printf("AI request failed: %v", err)
			http.Error(w, `{"error":"AI request failed"}`, http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageResponse{Reply: reply})
	})

	log.Println("noah API listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}