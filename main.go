package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	"github.com/wuhan-support/shimo"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
)

var (
	config                 Config
	Log                    *log.Logger
	AccommodationsDocument *shimo.Document
	PlatformDocument       *shimo.Document
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

//func SimulateDelay(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		time.Sleep(time.Millisecond * 500)
//		return next(c)
//	}
//}

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

	AccommodationsDocument = shimo.NewDocument(config.Documents.Accommodations, config.Cookie)
	AccommodationsDocument.Suffix = "（"

	PlatformDocument = shimo.NewDocument(config.Documents.Platforms, config.Cookie)
	PlatformDocument.Suffix = " ("

	e := echo.New()
	e.Debug = true
	e.Validator = &CustomValidator{validator: validator.New()}

	//e.Use(SimulateDelay)

	e.GET("/accommodations/json", func(c echo.Context) error {
		message, err := AccommodationsDocument.GetJSON()
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

	e.GET("/accommodations/csv", func(c echo.Context) error {
		csv, err := AccommodationsDocument.GetCSV()
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.Blob(http.StatusOK, "text/csv", csv)
	})

	e.GET("/platforms/json", func(c echo.Context) error {
		message, err := PlatformDocument.GetJSON()
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

	e.GET("/platforms/csv", func(c echo.Context) error {
		csv, err := PlatformDocument.GetCSV()
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.Blob(http.StatusOK, "text/csv", csv)
	})

	e.POST("/report", func(c echo.Context) error {
		var request ReportRequest
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}
		Log.Printf("[report] new report record: %v", spew.Sdump(request))
		return c.NoContent(http.StatusNoContent)
	})

	Log.Fatal(e.Start(config.Server.Address))
}
