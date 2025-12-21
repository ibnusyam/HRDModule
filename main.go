package main

import (
	"HRD/handler"
	handlers "HRD/handler"
	"HRD/internal/repository"
	"HRD/internal/service"
	"HRD/middleware" // 1. Import package middleware kamu
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware" // Rename biar gak bentrok
)

func main() {

    err := godotenv.Load()
    if err != nil {
        fmt.Println(err)
    }

    db, err := repository.ConnectDB()
    if err != nil {
        fmt.Println("gagal conect ke db :", err)
    }
    
    // Init layers
    cleaningLogRepository := repository.NewCleaningLogsRepository(db)
    cleaningLogService := service.NewCleaningLogService(cleaningLogRepository)
    cleaningLogHandler := handlers.NewCleaningLogHandler(cleaningLogService)

    dashboardRepository := repository.NewDashboardRepository(db)
    dashboardService := service.NewDashboardService(dashboardRepository)
    dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	locRepo := repository.NewLocationRepository(db)
	locService := service.NewLocationService(locRepo)
	locHandler := handler.NewLocationHandler(locService)


    e := echo.New()

    // Setup CORS
    e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    }))

    // ---------------------------------------------------------
    // ROUTE PUBLIC (Bisa diakses tanpa token)
    // ---------------------------------------------------------
    // Biasanya file gambar boleh diakses siapa saja
    e.Static("/uploads", "uploads") 


    // ---------------------------------------------------------
    // ROUTE PROTECTED (Harus punya Token Valid)
    // ---------------------------------------------------------
    // Kita bikin Group. Semua yang masuk group ini dicegat middleware
    protected := e.Group("") 
    protected.Use(middleware.JWTMiddleware) // Pasang "Satpam" di sini

    // Pindahkan route-route sensitif ke variable 'protected'
    protected.GET("/logs", cleaningLogHandler.GetAllLogs)
    protected.POST("/logs", cleaningLogHandler.CreateFullLog)
    protected.GET("/form-options", cleaningLogHandler.GetFormOptionsHandler)
    protected.GET("/dashboard/stats", dashboardHandler.GetCleanerStats)

	protected.GET("/location-types", locHandler.GetTypes)
	protected.POST("/location-types", locHandler.CreateType)
	protected.PUT("/location-types/:id", locHandler.UpdateType)
	protected.DELETE("/location-types/:id", locHandler.DeleteType)

	protected.GET("/locations", locHandler.GetLocations)
	protected.POST("/locations", locHandler.CreateLocation)
	protected.PUT("/locations/:id", locHandler.UpdateLocation)
	protected.DELETE("/locations/:id", locHandler.DeleteLocation)

    e.Logger.Fatal(e.Start(":8081"))
}