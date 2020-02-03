package main

import (
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	e.Use(middleware.Gzip())

	e.Validator = &CustomValidator{validator: validator.New()}

	//e.Use(SimulateDelay)

	shimoC := shimo_openapi.NewClient(config.Shimoauth.ClientId, config.Shimoauth.ClientSecret, config.Shimoauth.Username, config.Shimoauth.Password, config.Shimoauth.Scope)
	Log.Println(config, shimoC)

	// 返回住宿信息列表
	e.GET("/accommodations", func(c echo.Context) error {
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
	e.GET("/platforms/psychological", func(c echo.Context) error {
		fileId := "Dpy6Q668cj3Xx8Rq"
		opt := shimo_openapi.Opts{"上线版本", 19, "M", "\n", time.Minute * 30}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回线上医疗平台列表
	e.GET("/platforms/medical", func(c echo.Context) error {
		fileId := "DqpyXVgXCwdvqrht"
		opt := shimo_openapi.Opts{"总表", 263, "BR", "", time.Minute * 30}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回医院需求列表
	e.GET("/hospital/supplies", func(c echo.Context) error {
		fileId := "zN32MwmPjmCLF0Av"
		opt := shimo_openapi.Opts{"已合成", 426, "AP", " ", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回第二版的医院需求列表
	e.GET("/hospital/supplies/v2", func(c echo.Context) error {
		fileId := "DqpyXVgXCwdvqrht"
		opt := shimo_openapi.Opts{"总表", 300, "BR", "----", time.Minute * 30}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})


	// 返回武汉在外人员住宿信息
	e.GET("/people/accommodations", func(c echo.Context) error {
		fileId := "DR3OV8MN9yUxFnAB"
		opt := shimo_openapi.Opts{"工作表1", 934, "L", " ", time.Hour * 1}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	// 返回零散信息
	e.GET("/wiki/stream", func(c echo.Context) error {
		fileId := "XRkgJOMRW0CrFbqM"
		opt := shimo_openapi.Opts{"实时", 100, "H", " ", time.Minute * 3}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		return c.JSONBlob(http.StatusOK, message)
	})

	//e.GET("/hospital/supplies/submissions", func(c echo.Context) error {
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
	//e.POST("/hospital/supplies/submissions", func(c echo.Context) error {
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

	e.POST("/report", func(c echo.Context) error {
		var request ReportRequest
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}
		Log.Printf("[report] new report record: %v", spew.Sdump(request))
		return c.NoContent(http.StatusNoContent)
	})

	//e.POST("/upload", func(c echo.Context) error {
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
