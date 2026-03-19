package main

import (
	"flag"
	"log"
	"os"

	"github.com/dorgu-ai/dorgu-platform/pkg/server"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	kubeconfig := flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "Path to kubeconfig file")
	context := flag.String("context", "", "Kubernetes context to use")
	flag.Parse()

	config := &server.Config{
		Port:       *port,
		KubeConfig: *kubeconfig,
		Context:    *context,
	}

	srv := server.NewServer(config)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
