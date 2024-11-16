package server_config

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "ruantiengo/log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	http   *http.Server
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		router: router,
	}
}

func (s *Server) SetupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	// Add more routes here
}

func (s *Server) Start() {
	PORT := os.Getenv("API_PORT")
	s.http = &http.Server{
		Addr:    ":" + PORT,
		Handler: s.router,
	}

	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Message(logger.Error, "Failed to start server: %v", err)
		}
	}()

	logger.Message(logger.Info, "✔️ Server started on http://localhost:"+PORT)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Message(logger.Info, "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		logger.Message(logger.Error, "Server forced to shutdown: %v", err)
	}

	logger.Message(logger.Info, "Server exiting")
}

func (s *Server) Shutdown() {
	if s.http != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.http.Shutdown(ctx); err != nil {
			logger.Message(logger.Error, "Server forced to shutdown: %v", err)
		}
	}
}
