package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	shimo_openapi "github.com/wuhan-support/shimo-openapi"
	"gopkg.in/go-playground/validator.v9"
)

var (
	config Config
	Log    *log.Logger
	db     *gorm.DB
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

func init() {
	err := configor.Load(&config, "config.yml")
	if err != nil {
		Log.Fatalf("failed to initialize config file: %v", err)
	}

	//db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DBInfo.User, config.DBInfo.Pwd, config.DBInfo.Addr, config.DBInfo.DBName))
	//db, err = gorm.Open("sqlite3", "sqlite.db")
	//if err != nil {
	//	fmt.Println(fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", config.DBInfo.User, config.DBInfo.Pwd, config.DBInfo.Addr, config.DBInfo.DBName))
	//	panic("failed to connect mysql:" + err.Error())
	//}
	//
	//if !db.HasTable(Submission{}) {
	//	d := db.CreateTable(Submission{})
	//	if d.Error != nil {
	//		panic("create table failed:" + d.Error.Error())
	//	}
	//}

	//_, err = os.Stat(config.UploadPath)
	//if os.IsNotExist(err) {
	//	err = os.MkdirAll(config.UploadPath, os.ModePerm)
	//	if err != nil {
	//		panic("failed create upload path:" + err.Error())
	//	}
	//}
}

func main() {
	logFile, err := os.OpenFile("runtime.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	Log = log.New(logFile, "[http] ", log.LstdFlags)

	e := echo.New()
	e.Debug = false
	e.Validator = &CustomValidator{validator: validator.New()}

	//e.Use(SimulateDelay)

	shimoC := shimo_openapi.NewClient(config.Shimoauth.ClientId, config.Shimoauth.ClientSecret, config.Shimoauth.Username, config.Shimoauth.Password, config.Shimoauth.Scope)
	Log.Println(config, shimoC)

	api := e.Group("/api")

	// 返回住宿信息列表
	api.GET("/accommodations", func(c echo.Context) error {
		fileId := "6c6GKvX83hRCVdG8"
		opt := shimo_openapi.Opts{"工作表1", 278, "R", "（", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回心理咨询机构列表
	api.GET("/platforms/psychological", func(c echo.Context) error {
		fileId := "Dpy6Q668cj3Xx8Rq"
		opt := shimo_openapi.Opts{"工作表1", 17, "O", "\n", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回线上医疗平台列表
	api.GET("/platforms/medical", func(c echo.Context) error {
		fileId := "kDQJ6vWgWWwq8r8H"
		opt := shimo_openapi.Opts{"工作表1", 23, "O", " (", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回医院需求列表
	api.GET("/hospital/supplies", func(c echo.Context) error {
		fileId := "zN32MwmPjmCLF0Av"
		opt := shimo_openapi.Opts{"已合成", 160, "AP", " ", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	//api.GET("/hospital/supplies/submissions", func(c echo.Context) error {
	//	var request GetSubmissionsRequest
	//	var submissions []Submission
	//
	//	if c.Bind(&request) != nil && c.Validate(&request) != nil {
	//		fmt.Println(c.Validate(&request))
	//		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	//	}
	//
	//	paginator := pagination.Paging(&pagination.Param{
	//		DB:      db,
	//		Page:    request.Page,
	//		Limit:   request.Limit,
	//		OrderBy: []string{"id desc"},
	//	}, &submissions)
	//
	//	return c.JSON(http.StatusOK, paginator)
	//})
	//
	//api.POST("/hospital/supplies/submissions", func(c echo.Context) error {
	//	var request Submission
	//	if c.Bind(&request) != nil && c.Validate(&request) != nil {
	//		// fmt.Println(c.Bind(request))
	//		fmt.Println(c.Validate(&request))
	//		return echo.NewHTTPError(http.StatusBadRequest, "bad reqeust")
	//	}
	//
	//	d := db.Create(&request)
	//	if d.Error != nil {
	//		Log.Printf("create collect_form failed:%v", d.Error)
	//		return echo.NewHTTPError(http.StatusInternalServerError)
	//	}
	//	return c.NoContent(http.StatusNoContent)
	//})

	api.POST("/report", func(c echo.Context) error {
		var request ReportRequest
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}
		Log.Printf("[report] new report record: %v", spew.Sdump(request))
		return c.NoContent(http.StatusNoContent)
	})

	//api.POST("/upload", func(c echo.Context) error {
	//	fn, err := c.FormFile("file")
	//	if err != nil {
	//		return echo.NewHTTPError(http.StatusBadRequest, "file name not found")
	//	}
	//	fs, err := fn.Open()
	//	if err != nil {
	//		Log.Printf("open upload file %s failed:%v", fn.Filename, err)
	//		return echo.NewHTTPError(http.StatusInternalServerError)
	//	}
	//
	//	suf := ""
	//	index := strings.LastIndex(fn.Filename, ".")
	//	if index != -1 {
	//		suf = fn.Filename[index:]
	//	}
	//
	//	fname := tool.GetGUID() + suf
	//	fd, err := os.Create(fmt.Sprintf("%s/%s", config.UploadPath, fname))
	//	if err != nil {
	//		Log.Printf("create file %s failed:%v", fname, err)
	//		return echo.NewHTTPError(http.StatusInternalServerError)
	//	}
	//	if _, err = io.Copy(fd, fs); err != nil {
	//		Log.Printf("copy file failed:%v", err)
	//		return echo.NewHTTPError(http.StatusInternalServerError)
	//	}
	//
	//	return c.JSON(http.StatusOK, map[string]string{"fcode": fname})
	//})

	Log.Fatal(e.Start(config.Server.Address))
}
