package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

var port = os.Getenv("PORT")

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	if port == "" {
		port = "3000"
	}
}

func main() {
	// componenets
	server := &http.Server{Addr: ":" + port}
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "", DB: 0})
	b := NewBroker()

	// handlers
	http.HandleFunc("/publish", MW(handleSsePublish(rdb)))
	http.HandleFunc("/subscribe", MW(handleSseSubscribe(b)))
	http.HandleFunc("/ping", MW(handlePing()))

	// start the broker
	b.Run()

	// subscribe to redis
	backgroundCtx := context.Background()
	pubsub := rdb.Subscribe(backgroundCtx, "sse")
	_, err := pubsub.Receive(backgroundCtx)
	if err != nil {
		panic(err)
	}

	// receive messages from redis and
	// dispatch to all clients via the broker
	go func() {
		messages := pubsub.Channel()
		for msg := range messages {
			log.Printf("receiving message: %s", msg.Payload)
			b.Publish(msg.Payload)
		}
	}()

	// trap signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start the server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("server started on port " + port)

	// wait for a signal
	<-done

	// run teardown with timeout
	ctx, cancel := context.WithTimeout(backgroundCtx, 5*time.Second)
	defer func() {
		defer cancel()
		// add teardown code here
	}()

	// try to shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	// all good
	log.Println("server stopped")

}
