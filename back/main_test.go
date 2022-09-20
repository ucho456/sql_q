package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	// キャンセル可能なcontext.Contextのオブジェクトを作る
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	url := fmt.Sprintf("http://%s", l.Addr().String())
	t.Logf("try request to %q", url)

	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()

	// httpサーバーの戻り値を検証する
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	want := "Hello world!\n"
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// 終了通知を送信して戻り値を検証する
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
