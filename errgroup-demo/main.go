package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	srv := NewServer("8080")
	if err := srv.Start(); err != nil {
		log.Fatalln(err)
	}
}

type Server struct {
	HttpSrv *http.Server
}

func NewServer(port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/user", GetUser)
	return &Server{
		HttpSrv: &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	eg.Go(func() error {
		if s.HttpSrv != nil {
			log.Printf("http server listen address: %s\n", s.HttpSrv.Addr)
			return s.HttpSrv.ListenAndServe()
		}
		return nil
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("context done.")
		}
		return nil
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigs:
			cancel()
			return s.Stop(ctx)
		}
	})
	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("http server shutdown")
	if s.HttpSrv != nil {
		return s.HttpSrv.Shutdown(ctx)
	}
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	w.WriteHeader(200)
}
