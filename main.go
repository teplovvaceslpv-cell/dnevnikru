package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	install := exec.Command("pip", "install", "-r", "requirements.txt")
	install.Stdout = io.Discard
	install.Stderr = io.Discard
	install.Run()

	cmd := exec.Command("python", "loggerbot.py")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}