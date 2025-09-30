package searchProduct

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func SearchProduct(keyword string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Сначала диагностируем страницу
			if err := diagnoseSearchPage(ctx); err != nil {
				log.Printf("Ошибка диагностики: %v", err)
			}
			
			log.Printf("🔍 Выполняем поиск: %s", keyword)
			
			// Даем время на полную загрузку страницы
			time.Sleep(2 * time.Second)
			
			// Пробуем разные стратегии поиска
			strategies := []struct {
				name string
				fn   func(context.Context) error
			}{
				{"Поиск через основной input", func(ctx context.Context) error { return searchWithMainInput(ctx, keyword) }},
				{"Поиск через кнопку Enter", func(ctx context.Context) error { return searchWithEnter(ctx, keyword) }},
				{"Прямой URL", func(ctx context.Context) error { return searchWithDirectURL(ctx, keyword) }},
				{"Поиск через клик по предложению", func(ctx context.Context) error { return searchWithSuggestionClick(ctx, keyword) }},
			}
			
			for i, strategy := range strategies {
				log.Printf("Попытка #%d: %s", i+1, strategy.name)
				if err := strategy.fn(ctx); err == nil {
					// Проверяем, что поиск действительно выполнился
					if isSearchSuccessful(ctx) {
						log.Printf("✅ %s успешен", strategy.name)
						return nil
					} else {
						log.Printf("⚠️ %s выполнен, но результаты не загрузились", strategy.name)
					}
				} else {
					log.Printf("❌ %s не сработал: %v", strategy.name, err)
				}
				time.Sleep(1 * time.Second)
			}
			
			return fmt.Errorf("все методы поиска не сработали")
		}),
		
		// Ждем подтверждения успешного поиска
		chromedp.ActionFunc(func(ctx context.Context) error {
			return waitForSearchResults(ctx)
		}),
	}
}

func diagnoseSearchPage(ctx context.Context) error {
	log.Println("=== ДИАГНОСТИКА ПОИСКОВОЙ СТРАНИЦЫ ===")
	
	// Проверяем текущий URL
	var currentURL string
	if err := chromedp.Evaluate("window.location.href", &currentURL).Do(ctx); err == nil {
		log.Printf("URL: %s", currentURL)
	}
	
	// Проверяем доступные поля поиска
	searchSelectors := []string{
		`#searchInput`,
		`.search-catalog__input`,
		`input[placeholder*="поиск"]`,
		`input[placeholder*="найди"]`,
	}
	
	for _, selector := range searchSelectors {
		var exists bool
		if err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector("%s")`, selector), &exists).Do(ctx); err == nil && exists {
			log.Printf("✅ Найден элемент поиска: %s", selector)
			
			// Проверяем, доступен ли элемент
			var disabled bool
			chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").disabled`, selector), &disabled).Do(ctx)
			if !disabled {
				log.Printf("   ✅ Элемент доступен для ввода")
			} else {
				log.Printf("   ❌ Элемент заблокирован")
			}
		}
	}
	
	log.Println("=== КОНЕЦ ДИАГНОСТИКИ ===")
	return nil
}

func searchWithMainInput(ctx context.Context, keyword string) error {
	searchInputSelectors := []string{
		`#searchInput`,
		`.search-catalog__input`,
		`input[placeholder*="поиск"]`,
		`input[placeholder*="найди"]`,
		`.search__input`,
	}
	
	searchButtonSelectors := []string{
		`#applySearchBtn`,
		`.search-catalog__btn`,
		`button[type="submit"]`,
		`.search-btn`,
		`.header__search-submit`,
	}
	
	for _, inputSel := range searchInputSelectors {
		// Проверяем, существует ли поле ввода
		var inputExists bool
		if err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector("%s")`, inputSel), &inputExists).Do(ctx); err != nil || !inputExists {
			continue
		}
		
		log.Printf("Найдено поле ввода: %s", inputSel)
		
		// Очищаем поле ввода
		if err := chromedp.Run(ctx,
			chromedp.Focus(inputSel, chromedp.ByQuery),
			chromedp.SetValue(inputSel, "", chromedp.ByQuery),
		); err != nil {
			continue
		}
		
		// Вводим текст постепенно
		if err := chromedp.Run(ctx,
			chromedp.SendKeys(inputSel, keyword, chromedp.ByQuery),
			chromedp.Sleep(500*time.Millisecond),
		); err != nil {
			return err
		}
		
		// Пробуем разные кнопки поиска
		for _, btnSel := range searchButtonSelectors {
			if err := chromedp.Run(ctx,
				chromedp.WaitVisible(btnSel, chromedp.ByQuery),
				chromedp.Click(btnSel, chromedp.ByQuery),
			); err == nil {
				log.Printf("Нажата кнопка поиска: %s", btnSel)
				return nil
			}
		}
		
		// Если кнопки не сработали, пробуем отправить форму
		if err := chromedp.Run(ctx,
			chromedp.Submit(inputSel, chromedp.ByQuery),
		); err == nil {
			log.Println("Форма отправлена через submit")
			return nil
		}
	}
	
	return fmt.Errorf("не удалось выполнить поиск через input")
}

func searchWithEnter(ctx context.Context, keyword string) error {
    searchInputSelectors := []string{
        `#searchInput`,
        `.search-catalog__input`,
        `input[placeholder*="поиск"]`,
        `input[name="search"]`,
    }
    
    for _, inputSel := range searchInputSelectors {
        log.Printf("Пробуем Enter в селекторе: %s", inputSel)
        
        err := chromedp.Run(ctx,
            // Ждем и фокусируемся на поле ввода
            chromedp.WaitVisible(inputSel, chromedp.ByQuery),
            chromedp.Focus(inputSel, chromedp.ByQuery),
            chromedp.Sleep(100*time.Millisecond),
            
            // Очищаем поле
            chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").value = ""`, inputSel), nil),
            chromedp.Sleep(100*time.Millisecond),
            
            // Вводим текст
            chromedp.SendKeys(inputSel, keyword, chromedp.ByQuery),
            chromedp.Sleep(500*time.Millisecond),
            
            // СПОСОБ 1: Отправляем Enter через chromedp.SendKeys
            chromedp.SendKeys(inputSel, "\n", chromedp.ByQuery),
            
            // Ждем обработки
            chromedp.Sleep(1*time.Second),
        )
        
        if err == nil {
            log.Printf("✅ Enter отправлен через SendKeys с селектором: %s", inputSel)
            
            // Проверяем, сработал ли поиск
            if isSearchSuccessful(ctx) {
                return nil
            }
        }
        
        log.Printf("❌ SendKeys Enter не сработал с %s: %v", inputSel, err)
        
        // Если не сработало, пробуем JavaScript события
        err = chromedp.Run(ctx,
            chromedp.Focus(inputSel, chromedp.ByQuery),
            chromedp.Sleep(100*time.Millisecond),
            
            // Вводим текст снова
            chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").value = "%s"`, inputSel, keyword), nil),
            chromedp.Sleep(500*time.Millisecond),
            
            // СПОСОБ 2: JavaScript события клавиатуры
            chromedp.Evaluate(fmt.Sprintf(`
                (function() {
                    var input = document.querySelector("%s");
                    
                    // Создаем полную последовательность событий Enter
                    var events = [
                        new KeyboardEvent('keydown', {
                            key: 'Enter',
                            code: 'Enter',
                            keyCode: 13,
                            which: 13,
                            bubbles: true,
                            cancelable: true
                        }),
                        new KeyboardEvent('keypress', {
                            key: 'Enter',
                            code: 'Enter',
                            keyCode: 13,
                            which: 13,
                            bubbles: true,
                            cancelable: true
                        }),
                        new KeyboardEvent('keyup', {
                            key: 'Enter',
                            code: 'Enter',
                            keyCode: 13,
                            which: 13,
                            bubbles: true,
                            cancelable: true
                        })
                    ];
                    
                    // Отправляем все события
                    events.forEach(function(event) {
                        input.dispatchEvent(event);
                    });
                    
                    // Пробуем отправить форму
                    var form = input.closest('form');
                    if (form) {
                        var submitEvent = new Event('submit', { bubbles: true });
                        form.dispatchEvent(submitEvent);
                    }
                    
                    return true;
                })()
            `, inputSel), nil),
            
            chromedp.Sleep(1*time.Second),
        )
        
        if err == nil {
            log.Printf("✅ Enter отправлен через JavaScript события с селектором: %s", inputSel)
            
            if isSearchSuccessful(ctx) {
                return nil
            }
        }
        
        log.Printf("❌ JavaScript Enter не сработал с %s: %v", inputSel, err)
        
        // СПОСОБ 3: Пробуем найти и нажать кнопку поиска после ввода текста
        err = chromedp.Run(ctx,
            chromedp.Focus(inputSel, chromedp.ByQuery),
            chromedp.Sleep(100*time.Millisecond),
            
            // Вводим текст
            chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").value = "%s"`, inputSel, keyword), nil),
            chromedp.Sleep(500*time.Millisecond),
            
            // Ищем и нажимаем кнопку поиска
            chromedp.Click(`#applySearchBtn, .search-catalog__btn, button[type="submit"]`, chromedp.ByQuery),
            chromedp.Sleep(1*time.Second),
        )
        
        if err == nil {
            log.Printf("✅ Кнопка поиска нажата после ввода текста в: %s", inputSel)
            
            if isSearchSuccessful(ctx) {
                return nil
            }
        }
    }
    
    return fmt.Errorf("поиск через Enter не сработал ни с одним селектором")
}

func searchWithDirectURL(ctx context.Context, keyword string) error {
	searchURL := fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?search=%s", url.QueryEscape(keyword))
	log.Printf("Переходим по прямому URL: %s", searchURL)
	return chromedp.Navigate(searchURL).Do(ctx)
}

func searchWithSuggestionClick(ctx context.Context, keyword string) error {
	searchInputSelectors := []string{
		`#searchInput`,
		`.search-catalog__input`,
	}
	
	for _, inputSel := range searchInputSelectors {
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(inputSel, chromedp.ByQuery),
			chromedp.SetValue(inputSel, "", chromedp.ByQuery),
			chromedp.SendKeys(inputSel, keyword, chromedp.ByQuery),
			chromedp.Sleep(1*time.Second),
		); err != nil {
			continue
		}
		
		// Пробуем кликнуть по первой подсказке
		suggestionSelectors := []string{
			`.search-suggest__item`,
			`.suggest-item`,
			`.search-suggestion`,
			`[data-suggest-item]`,
		}
		
		for _, suggestionSel := range suggestionSelectors {
			if err := chromedp.Run(ctx,
				chromedp.WaitVisible(suggestionSel, chromedp.ByQuery),
				chromedp.Click(suggestionSel, chromedp.ByQuery),
			); err == nil {
				log.Printf("Клик по подсказке: %s", suggestionSel)
				return nil
			}
		}
	}
	
	return fmt.Errorf("поиск через подсказки не сработал")
}

func isSearchSuccessful(ctx context.Context) bool {
	successIndicators := []string{
		`.searching-results`,
		`[data-tag="searchResults"]`,
		`.catalog-page`,
		`.goods-list`,
		`.product-card`,
	}
	
	for _, indicator := range successIndicators {
		var exists bool
		if err := chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector("%s")`, indicator), &exists).Do(ctx); err == nil && exists {
			return true
		}
	}
	
	// Проверяем URL
	var currentURL string
	if err := chromedp.Evaluate("window.location.href", &currentURL).Do(ctx); err == nil {
		if strings.Contains(currentURL, "search") || strings.Contains(currentURL, "Search") {
			return true
		}
	}
	
	return false
}

func waitForSearchResults(ctx context.Context) error {
	log.Println("⏳ Ожидаем загрузки результатов поиска...")
	
	timeout := time.After(15 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("таймаут ожидания результатов поиска")
		case <-ticker.C:
			if isSearchSuccessful(ctx) {
				log.Println("✅ Результаты поиска загружены")
				return nil
			}
			log.Println("⏳ Ожидаем загрузки результатов...")
		}
	}
}