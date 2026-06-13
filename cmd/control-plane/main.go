package main

import (
	"kopeika/internal/app"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	dsn := "postgres://user:admin@pg:5432/kopeika?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer db.Close()

	AppRepo := app.NewAppRepository(db)
	AppService := app.NewService(AppRepo)
	AppHandler := app.NewHandler(AppService)
	r.GET("/api/v1", AppHandler.GetByID)
	r.GET("/api/v1/name", AppHandler.GetByName)
	r.POST("/api/v1", AppHandler.Create)
	r.PUT("/api/v1/envs", AppHandler.UpdateEnvVars)
	r.PUT("/api/v1/status", AppHandler.UpdateStatus)

	r.Run(":8000")
}
