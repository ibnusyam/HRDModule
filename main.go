package main

import (
	handlers "HRD/handler"
	"HRD/internal/repository"
	"HRD/internal/service"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main(){

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	db, err := repository.ConnectDB()
	if err != nil {
		fmt.Println("gagal conect ke db :", err)
	}
	fmt.Println(db)

	cleaningLogRepository := repository.NewCleaningLogsRepository(db)
	cleaningLogService := service.NewCleaningLogService(cleaningLogRepository)
	cleaningLogHandler := handlers.NewCleaningLogHandler(cleaningLogService)

	dashboardRepository := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepository)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	e := echo.New()

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))


	e.Static("/uploads", "uploads")
	e.GET("/logs", cleaningLogHandler.GetAllLogs)
    e.POST("/logs", cleaningLogHandler.CreateFullLog)
	e.GET("/form-options", cleaningLogHandler.GetFormOptions)
	e.GET("/dashboard/stats", dashboardHandler.GetCleanerStats)

	e.Logger.Fatal(e.Start(":8081"))
}