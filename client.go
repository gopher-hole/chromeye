package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func NewClient(url string) (context.Context, context.CancelFunc, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[0:2], chromedp.DefaultExecAllocatorOptions[3:]...)

	alloctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(alloctx, chromedp.WithLogf(log.Printf))
	if err := chromedp.Run(ctx); err != nil {
		return nil, nil, err
	}

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return nil, nil, err
	}

	return ctx, cancel, nil
}
