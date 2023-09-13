package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"log"
	"net/http"
	"os"
	"time"
)

type formPayload struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

func main() {
	publicURL := os.Getenv("PUBLIC_URL")
	appPass := os.Getenv("PASSWORD")
	formTarget := os.Getenv("FORM_TARGET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("dist")))

	// Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("healthy!")
		w.WriteHeader(http.StatusOK)
	})

	// Handle POST requests to /send
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON formPayload
		var payload formPayload
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON formPayload", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		m := gomail.NewMessage()

		m.SetHeader("From", payload.Email)
		m.SetHeader("To", formTarget)
		m.SetHeader("Subject", fmt.Sprintf("Melding fra %s (via blizterapi.no)", payload.Name))
		m.SetBody("text/plain", fmt.Sprintf("Melding:\n%s\n\nAvsender: %s (%s)\n", payload.Message, payload.Name, payload.Email))

		// Settings for SMTP server
		d := gomail.NewDialer("smtp.gmail.com", 587, "vev.kontaktskjema@gmail.com", appPass)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// keep app alive by polling health endpoint
	req, err := http.NewRequest(http.MethodGet, publicURL+"/health", nil)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		t := time.NewTicker(time.Minute * 2)
		for range t.C {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			r := req.WithContext(ctx)

			if _, err := http.DefaultClient.Do(r); err != nil {
				log.Println(err)
			}
		}
	}() // Start the server

	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
