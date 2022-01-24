package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

func handleSsePublish(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		}

		log.Printf("publishing message: %s", bodyBytes)

		err = rdb.Publish(r.Context(), "sse", bodyBytes).Err()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleSseSubscribe(b *Broker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// we got a new client
		fmt.Printf("client connected: %v\n", r.RemoteAddr)

		// the headers for event-stream
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering:", "no")

		// since this is a long-lived connection, we need to
		// flush the data to the client with this flusher
		f := w.(http.Flusher)

		// acknowledge the connection
		sse := SSE{"", "ack", "", 1000}
		sse.WriteTo(w)
		f.Flush()

		sub := b.Subscribe()

		for {
			select {
			// if the client closes the connection
			// remove the subscription
			case <-r.Context().Done():
				b.Unsubscribe(sub)
				fmt.Printf("client disconnected: %v\n", r.RemoteAddr)
				return
			// if the broker has a message
			// send it to the client
			case msg := <-sub:
				sse := SSE{"", "sse", msg, 1000}
				sse.WriteTo(w)
				f.Flush()
			// send a ping every 30 seconds
			// to keep the connection alive
			case <-time.After(time.Second * 30):
				sse := SSE{"", "ping", "", 1000}
				sse.WriteTo(w)
				f.Flush()
			}
		}
	}
}

func handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
