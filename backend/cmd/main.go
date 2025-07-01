package main

import (
	"log"

	// "fyne.io/fyne/v2/storage/repository"
	"github.com/gin-gonic/gin"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/config"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/handler"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/repository"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/service"
)

func main() {
	cfg := config.Load()
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()
	h.InitRoutes(r)

	r.Run(":8080")
}
