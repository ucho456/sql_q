package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	l, err := net.Listen("tcp", ":38080")
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	if err := run(context.Background(), l); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context, l net.Listener) error {
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello world!")
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでhttpサーバーを起動する
	eg.Go(func() error {
		if err := srv.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を待機する
	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Goメソッドで起動した別ゴルーチンの終了を待つ
	return eg.Wait()
}
