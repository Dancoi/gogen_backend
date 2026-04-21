package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация маршрутизатора Gin
	router := gin.Default()

	// Здесь будут добавлены маршруты

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
