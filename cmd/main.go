package main

import (
	"log"
	"net/http"
	"timesheet/cmd/handler"
	"timesheet/config"
	"timesheet/internal/repository"
	"timesheet/internal/timesheet"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	databaseConfig, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("Setup config error %s", err.Error())
	}
	databaseConnection, err := sqlx.Connect("mysql", databaseConfig.GetURI())
	if err != nil {
		log.Fatal("Cannot connect database", err.Error())
	}
	defer databaseConnection.Close()
	timesheetRepository := repository.TimesheetRepository{
		DatabaseConnection: databaseConnection,
	}
	api := handler.TimesheetAPI{
		Timesheet: timesheet.Timesheet{
			Repository: timesheetRepository,
		},
		Repository:            timesheetRepository,
		RepositoryToTimesheet: timesheetRepository,
	}

	router := gin.Default()
	router.GET("/login", handler.OauthGoogleLogin)
	router.GET("/callback", api.OauthGoogleCallback)
	router.POST("/logout", api.OauthGoogleLogout)
	router.GET("/deleteOauthState", handler.DeleteOauthStateCookie)
	router.GET("/showProfile", api.GetProfileHandler)
	router.POST("/showSummaryTimesheet", api.GetSummaryHandler)
	router.POST("/createIncome", api.CreateIncomeHandler)
	router.POST("/calculatePayment", api.CalculatePaymentHandler)
	router.POST("/showTimesheetByEmployeeID", api.GetSummaryByEmployeeIDHandler)
	router.POST("/updateStatusCheckingTransfer", api.UpdateStatusCheckingTransferHandler)
	router.POST("/deleteIncomeItem", api.DeleteIncomeHandler)
	router.POST("/showEmployeeDetailsByEmployeeID", api.ShowEmployeeDetailsByEmployeeIDHandler)
	router.POST("/updateEmployeeDetails", api.UpdateEmployeeDetailsHandler)
	router.POST("/showSummaryInYear", api.ShowSummaryInYearHandler)
	router.StaticFS("/home", http.Dir("ui"))
	log.Fatal(router.Run())
}
