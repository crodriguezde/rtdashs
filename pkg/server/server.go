package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/crodriguezde/rtdashs/pkg/sink"
	"github.com/crodriguezde/rtdashs/static"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	ctx context.Context
	sk  *sink.Server
}

func NewHandler(ctx context.Context, sk *sink.Server) *mux.Router {
	h := &handler{
		ctx: ctx,
		sk:  sk,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/stream/cpu", h.cpuStream)
	router.PathPrefix("/").Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.FS(static.FS)).ServeHTTP(w, r)
		}))

	return router
}

func (h *handler) cpuStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("stream connected %p", r)
	id := uuid.New().String()

	recv, err := h.sk.Register(id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		log.Printf("stream disconnected %p", r)
		h.sk.Deregister(id)
	}()

	w.Header().Set("content-Type", "text/event-stream")
	w.Header().Set("cache-Control", "no-store")

	for {
		select {
		case avgEvent := <-recv:
			if buf, err := json.Marshal(avgEvent); err != nil {
				log.Printf("cannot marshal event: %s", err)
				return
			} else if _, err := fmt.Fprintf(w, "event: update\ndata: %s\n\n", buf); err != nil {
				log.Printf("cannot write update event: %s", err)
				return
			}
			w.(http.Flusher).Flush()
		case <-h.ctx.Done():
			return
		case <-r.Context().Done():
			return
		}
	}
}
