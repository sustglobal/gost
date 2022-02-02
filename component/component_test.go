package component

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func newUnixDomainSocket(t *testing.T) (string, *http.Client) {
	dir, err := os.MkdirTemp("", "gust-*")
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("failed to remove socket temp directory: %v", err)
		}
	})

	sockname := filepath.Join(dir, "gust.sock")

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockname)
			},
		},
	}

	return sockname, httpClient
}

func TestComponentStartBindStop(t *testing.T) {
	sockname, httpClient := newUnixDomainSocket(t)

	cfg := DefaultConfig()
	cfg.BindHTTPServer = fmt.Sprintf("unix://%s", sockname)
	cfg.ExposeMetrics = true

	cmp, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := cmp.Start(); err != nil {
		t.Errorf("Component.Start failed with err=%v", err)
	}

	time.Sleep(3)

	req, err := http.NewRequest("GET", "http://unix/debug/vars", nil)
	if err != nil {
		t.Fatalf("Failed building HTTP request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Errorf("HTTP request failed with err=%v", err)
	} else if resp.StatusCode != 200 {
		t.Errorf("HTTP request returned unexpected status code: want=200 got=%d", resp.StatusCode)
	}

	if err := cmp.Stop(); err != nil {
		t.Errorf("Component.Stop failed with err=%v", err)
	}

	if err := cmp.Error(); err != nil {
		t.Errorf("Component.Error returned err=%v", err)
	}
}

func TestComponentRun(t *testing.T) {
	sockname, _ := newUnixDomainSocket(t)

	cfg := DefaultConfig()
	cfg.BindHTTPServer = fmt.Sprintf("unix://%s", sockname)

	cmp, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errc := make(chan error, 1)
	sigc := make(chan os.Signal, 1)
	go func() {
		err := cmp.runUntil(sigc)
		errc <- err
		close(errc)
	}()

	sigc <- syscall.SIGINT
	close(sigc)

	select {
	case err, ok := <-errc:
		if !ok {
			t.Fatalf("channel closed early")
		}
		if err != nil {
			t.Fatalf("expected nil err value")
		}
	case <-time.NewTimer(2 * time.Second).C:
		t.Fatalf("timed out waiting for err")
	}

	if err := cmp.Error(); err != nil {
		t.Errorf("Component.Error returned err=%v", err)
	}
}
