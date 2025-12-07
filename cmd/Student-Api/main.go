package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/akshayjha21/Student-Api/internal/config"
	"github.com/akshayjha21/Student-Api/internal/http/handler/student"
	"github.com/akshayjha21/Student-Api/internal/storage/sqlite"
)

func main() {
	//TODO - load config

	cfg := config.MustLoad()

	//TODO - database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage intializer ", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//TODO - setup route
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	//TODO - setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("Address", cfg.Addr))
	// fmt.Println("Server has started ")

	//applying graceful shutdown

	/*	Start server in background goroutine
		↓
		Wait for OS signal (CTRL+C)
		↓
		Signal received → Start graceful shutdown
		↓
		Create 5 sec timeout context
		↓
		server.Shutdown(ctx)
				- Stop taking new requests
				- Finish current requests
				- Close connections safely
		↓
		Shutdown complete
	*/

	// 1. We first create a 'done' channel to receive OS shutdown signals
	done := make(chan os.Signal, 1)

	// 2. We tell Go: “Whenever these signals come → send them into 'done' channel”
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 3. Start the HTTP server in a separate goroutine
	//    because ListenAndServe() blocks forever while serving requests
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to connect to the server")
		}
	}()

	// 4. Now main goroutine waits here until a signal arrives (like CTRL+C)
	//    This is blocking — program stops here until OS sends a shutdown signal
	<-done

	// -------------------------
	// Now SHUTDOWN STARTS HERE
	// -------------------------

	slog.Info("Shutting down the server")

	// 5. Create a context with a 5-second timeout
	//    Why? Because we want the server to shutdown gracefully but not forever.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 6. Politely ask the server to shutdown.
	//    - It will stop taking new requests
	//    - It will finish ongoing requests
	//    - Uses the timeout context above
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
