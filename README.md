# Credit Card Validator

## 📋 Описание проекта

Основанный на Go инструмент скрапинга товаров веб-магазина Wildberis, с использованием удалённой отладки chrome с помощью websocket.

## 🛠 Основные функции

- Получение названия продукта с помощью флага
- Автоматезировованое получение WebSocketa для удалённой отладки
- Флаг help
- Вывод всех найденых товаров с информацие в консоль в виде таблички 

## 🚀 Установка и запуск

### Требования

- Go 1.22 или выше

### Установка

1. Клонируйте репозиторий

    ```markdown
        
    ```bash
    git clone https://github.com/Alladinchik7/WebKrauler.git
    ```

2. Перейдите в директорию проекта

    ```markdown

    ```bash
    cd WebKrauler
    ```

3. Запустите программу

    ```markdown

    ```bash
    go run cmd/main.go
    или
    .\scrap.exe
    ```

## 🎮 Примеры использования

1. Скрапинг с указание товара  

    ```markdown

    ```bash
    .\scrap.exe -product="iphone"
    ```

2. Флаг help

    ```markdown

    ```bash
    .\scrap.exe -help
    ```
