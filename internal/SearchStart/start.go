package searchstart

import (
	"context"

	"github.com/chromedp/chromedp"
)

func SearchStart(url string) (context.Context, context.CancelFunc){
	var baseCtx context.Context

	if url != "" {
		baseCtx, _ = chromedp.NewRemoteAllocator(context.Background(), url)
	} else {
		baseCtx = context.Background()
	}

	return chromedp.NewContext(baseCtx)
}