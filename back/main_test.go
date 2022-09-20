package main

import (
	"context"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// キャンセル可能なcontext.Contextのオブジェクトを作る
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	rsp, err := http.Get("http://localhost:38080")
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
