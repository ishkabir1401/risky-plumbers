package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"risky-plumbers/middleware"
	"risky-plumbers/routes"
	"risky-plumbers/settings"
	"syscall"
	"time"
)

func main() {
	settings.InitializeRiskStore()
	router := gin.Default()
	router.Use(middleware.RecoverFromPanic)
	routes.SetupRoutes(router)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_, cancel := context.WithCancel(context.Background())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start server", err, nil)
			os.Exit(1)
		}
	}()
	<-sig
	fmt.Println("Shutting down gracefully", nil)
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Failed to gracefully shutdown server", err, nil)
		os.Exit(1)
	}
}
