package main

import (
	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	"github.com/wuhan-support/shimo"
	"log"
	"net/http"
	"os"
)

var config Config
var Log *log.Logger
var HotelDocument *shimo.Document

func main() {
	logFile, err := os.OpenFile("runtime.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	Log = log.New(logFile, "[http] ", log.LstdFlags)

	err = configor.Load(&config, "config.yml")
	if err != nil {
		Log.Fatalf("failed to initialize config file: %v", err)
	}

	HotelDocument = shimo.NewDocument(config.Documents.Hotel, config.Cookie)
	HotelDocument.EliminateSuffix = "（"

	e := echo.New()
	e.GET("/hotels", func(c echo.Context) error {
		message, err := HotelDocument.GetJSON()
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		marshalled, err := message.MarshalJSON()
		if err != nil {
			Log.Printf("failed to marshal json: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal json")
		}
		return c.JSONBlob(http.StatusOK, marshalled)
	})

	Log.Fatal(e.Start(config.Server.Address))
}
