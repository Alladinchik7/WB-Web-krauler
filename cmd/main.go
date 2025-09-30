package main

import (
	getproduct "WebCrauler/internal/Product/getProduct"
	"WebCrauler/internal/Product/printProduct"
	searchrun "WebCrauler/internal/SearchRun"
	searchstart "WebCrauler/internal/SearchStart"
	websocket "WebCrauler/internal/WebSocket"
	helpFlag "WebCrauler/pkg/help"
	"flag"
	"log"
)

func main() {
	short := flag.Bool("short", true, "shorten image URLs in output")
	help := flag.Bool("help", false, "help flags")
	productFlag := flag.String("product", "gopher", "your product name")
	websocketFlag := flag.String("webSocketDebuggerUrl", "", "your websocket url")

	flag.Parse()

	product := *productFlag
	websocketURL := *websocketFlag

	if websocketURL == "" {
		var err error

		websocketURL, err = websocket.InitWebSocket()
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
	}
	defer websocket.KillExistingChrome()

	if *help {
		helpFlag.Help()
	}

	if websocketURL != "" {
		log.Printf("Using remote browser: %s\n", websocketURL)
	}

	ctx, cancel := searchstart.SearchStart(websocketURL)
	defer cancel()

	nodes, err := searchrun.RunSearch(ctx, product)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Получено узлов для обработки: %d", len(nodes))
	
	products := getproduct.GetProducts(nodes, ctx, short)
	log.Printf("Успешно собрано данных товаров: %d", len(products))

	printProduct.PrintProducts(products, product)
}