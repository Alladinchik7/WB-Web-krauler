// internal/SearchRun/searchRun.go
package searchrun

import (
	getNodes "WebCrauler/internal/GetNodes"
	searchProduct "WebCrauler/internal/Product/SearchProduct"
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// internal/SearchRun/searchRun.go
func RunSearch(ctx context.Context, product string) ([]*cdp.Node, error) {
	var nodes []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.wildberries.ru/`),
		searchProduct.SearchProduct(product),
		// Ждем результаты поиска
		chromedp.WaitReady(".product-card", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		// Собираем узлы БЕЗ скроллинга для теста
		getNodes.GetNodesWithScroll(&nodes),
	)
	
	if err != nil {
		return nil, err
	}

	log.Printf("Поиск завершен. Найдено товаров: %d", len(nodes))
	return nodes, nil
}