package getProduct

import (
	product "WebCrauler/internal/Product"
	shortImg "WebCrauler/pkg/short"
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func GetProducts(nodes []*cdp.Node, ctx context.Context, short *bool) []*product.Product {
	products := make([]*product.Product, 0, len(nodes))
	
	log.Printf("Начинаем обработку %d товаров...", len(nodes))

	successCount := 0
	for i, v := range nodes {
		log.Printf("Обрабатываем товар %d/%d", i+1, len(nodes))
		
		product, err := GetProductData(ctx, v)
		if err != nil {
			log.Printf("❌ Ошибка обработки товара %d: %v", i+1, err)
			continue
		}
		
		if *short {
			product.Img = shortImg.ShortImg(product.Img)
		}
		
		// Проверяем, что есть основные данные
		if product.Title != "" && product.Price != "" {
			products = append(products, product)
			successCount++
		} else {
			log.Printf("⚠️ Товар %d пропущен - отсутствуют данные", i+1)
		}
	}
	
	log.Printf("✅ Успешно обработано %d из %d товаров", successCount, len(nodes))
	return products
}

func GetProductData(ctx context.Context, node *cdp.Node) (*product.Product, error) {
	product := &product.Product{}
	
	// Получаем ID из data-nm-id атрибута
	product.ID = node.AttributeValue("data-nm-id")
	if product.ID == "" {
		// Если нет data-nm-id, пробуем извлечь из класса или других атрибутов
		product.ID = node.AttributeValue("id")
	}
	
	err := chromedp.Run(ctx,
		// Изображение - основные селекторы WB
		chromedp.AttributeValue(`
			.product-card__img img,
			.j-card-img img,
			img[alt*="товар"],
			img
		`, "src", &product.Img, nil, chromedp.ByQuery, chromedp.FromNode(node)),
		
		// Название товара
		chromedp.TextContent(`
			.product-card__name,
			.j-card-name,
			.goods-name,
			.card__name
		`, &product.Title, chromedp.ByQuery, chromedp.FromNode(node)),
		
		// Цена
		chromedp.TextContent(`
			.price__lower-price,
			.lower-price,
			.final-cost,
			.j-final-price,
			.price
		`, &product.Price, chromedp.ByQuery, chromedp.FromNode(node)),
	)
	
	if err != nil {
		log.Printf("❌ Ошибка получения данных товара %s: %v", product.ID, err)
		return nil, err
	}
	
	// Очистка данных
	product.Title = strings.TrimSpace(product.Title)
	product.Price = strings.TrimSpace(product.Price)
	
	// Если изображение не найдено, пробуем data-src
	if product.Img == "" || strings.Contains(product.Img, "data:image") {
		var dataSrc string
		chromedp.AttributeValue(`
			.product-card__img img,
			.j-card-img img
		`, "data-src", &dataSrc, nil, chromedp.ByQuery, chromedp.FromNode(node)).Do(ctx)
		if dataSrc != "" {
			product.Img = dataSrc
		}
	}
	
	log.Printf("✅ Товар обработан: %s - %s", product.Title, product.Price)
	return product, nil
}