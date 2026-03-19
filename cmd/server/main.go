package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dorgu-ai/dorgu-platform/pkg/server"
)

func main() {
	// Placeholder for standalone server binary
	// Agent 2 will implement the actual server logic

	fmt.Println("Dorgu Platform Server")
	fmt.Println("This is a placeholder. Run 'make build' after Agent 2 completes.")

	srv := &server.Server{
		Port:       8080,
		KubeConfig: os.Getenv("KUBECONFIG"),
		Context:    "",
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
