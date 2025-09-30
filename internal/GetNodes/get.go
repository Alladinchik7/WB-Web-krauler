// internal/GetNodes/getNodes.go
package getNodes

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func GetNodesWithScroll(nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			return collectNodesWithScroll(ctx, nodes)
		}),
	}
}

func collectNodesWithScroll(ctx context.Context, nodes *[]*cdp.Node) error {
	// Сначала собираем начальные узлы
	productSelectors := []string{
		`.product-card`,
		`.j-card-item`, 
		`div[data-nm-id]`,
	}
	
	var allNodes []*cdp.Node
	var workingSelector string
	
	// Находим рабочий селектор
	for _, selector := range productSelectors {
		var currentNodes []*cdp.Node
		err := chromedp.Nodes(selector, &currentNodes, chromedp.ByQueryAll).Do(ctx)
		if err == nil && len(currentNodes) > 0 {
			log.Printf("✅ Селектор %s нашел %d товаров", selector, len(currentNodes))
			allNodes = currentNodes
			workingSelector = selector
			break
		}
	}
	
	if len(allNodes) == 0 {
		log.Println("❌ Не найдено товаров")
		*nodes = allNodes
		return nil
	}
	
	log.Printf("Начальное количество товаров: %d", len(allNodes))
	
	// Скроллим с таймаутами
	for i := 0; i < 3; i++ { // Уменьшаем количество попыток
		log.Printf("Скроллим страницу... (попытка %d/3)", i+1)
		
		// Используем chromedp.ScrollIntoView вместо Evaluate
		err := chromedp.Run(ctx,
			// Скроллим к низу страницы
			chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight)`, nil),
			// Ждем немного
			chromedp.Sleep(2*time.Second),
		)
		
		if err != nil {
			log.Printf("❌ Ошибка при скролле: %v", err)
			break
		}
		
		// Проверяем, появились ли новые товары
		var newNodes []*cdp.Node
		err = chromedp.Nodes(workingSelector, &newNodes, chromedp.ByQueryAll).Do(ctx)
		if err != nil {
			log.Printf("❌ Ошибка при сборе узлов после скролла: %v", err)
			break
		}
		
		log.Printf("После скролла #%d: всего %d товаров", i+1, len(newNodes))
		
		if len(newNodes) > len(allNodes) {
			log.Printf("✅ Загружено новых товаров: %d", len(newNodes)-len(allNodes))
			allNodes = newNodes
		} else {
			log.Println("📌 Новых товаров не загружено, завершаем скроллинг")
			break
		}
		
		// Проверяем, есть ли еще контент для загрузки
		var canScrollMore bool
		err = chromedp.EvaluateAsDevTools(`
			window.innerHeight + window.pageYOffset < document.body.scrollHeight - 100
		`, &canScrollMore).Do(ctx)
		
		if err != nil || !canScrollMore {
			log.Println("📌 Достигнут конец страницы")
			break
		}
	}
	
	*nodes = allNodes
	log.Printf("✅ Итог: собрано %d товаров", len(allNodes))
	return nil
}