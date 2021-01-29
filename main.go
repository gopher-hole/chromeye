package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func onClose(cancel context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("cleaning up... please wait...")
		cancel()
		os.Exit(0)
	}()
}

func main() {
	addr := flag.String("addr", "http://localhost:3000", "-addr [address]")
	extensions := flag.String("ext", "gohtml,html,css,js", "-ext [gohtml,html,css,js]")
	flag.Parse()

	ctx, cancel, err := NewClient(*addr)
	if err != nil {
		log.Fatal(err)
	}
	onClose(cancel)

	extcache := map[string]bool{}
	for _, v := range strings.Split(*extensions, ",") {
		extcache[v] = true
	}

	if err := watch(ctx, extcache); err != nil {
		cancel()
		log.Println(err)
	}
}

func watch(ctx context.Context, extensions map[string]bool) error {
	files := map[string]time.Time{}
	for {
		var refresh bool
		err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !extensions[strings.TrimPrefix(filepath.Ext(path), ".")] {
				return nil
			}

			t, ok := files[path]
			if !ok {
				refresh = true
				files[path] = info.ModTime()
			}

			if !t.Equal(info.ModTime()) {
				refresh = true
				files[path] = info.ModTime()
			}

			return nil
		})

		if err != nil {
			return err
		}

		if refresh {
			log.Println("refreshing...")
			if err := chromedp.Run(ctx, chromedp.Reload()); err != nil {
				return err
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}
