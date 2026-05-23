package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/frankfoo/doraemon/internal/api"
	mcpserver "github.com/frankfoo/doraemon/internal/mcp"
	"github.com/frankfoo/doraemon/internal/store"
	"github.com/frankfoo/doraemon/web"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	fmt.Println("Doraemon starting...")

	s, err := store.New("data")
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	defer s.Close()

	addr := os.Getenv("DORAEMON_HTTP_PORT")
	if addr == "" {
		addr = ":8080"
	} else {
		addr = ":" + addr
	}

	jwtSecret := getJWTSecret()

	webFS, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		log.Fatalf("Failed to load embedded web assets: %v", err)
	}

	apiServer := api.New(s, addr, jwtSecret, webFS)
	go func() {
		log.Printf("Admin panel listening on %s", apiServer.Addr())
		if err := apiServer.ListenAndServe(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	srv := mcpserver.NewServer(s)
	go func() {
		fmt.Println("MCP server serving via stdio...")
		if err := server.ServeStdio(srv); err != nil {
			log.Printf("MCP server stopped: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Shutting down...")
}

func getJWTSecret() []byte {
	if secret := os.Getenv("DORAEMON_JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	b := make([]byte, 32)
	rand.Read(b)
	return []byte(hex.EncodeToString(b))
}
