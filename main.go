package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/middleware"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	shimo_openapi "github.com/wuhan-support/shimo-openapi"
	"gopkg.in/go-playground/validator.v9"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	config Config
	Log    *log.Logger
	db     *gorm.DB
	tgbot  *tgbotapi.BotAPI
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

	w := io.MultiWriter(os.Stdout, os.Stderr, logFile)

	Log = log.New(w, "[http] ", log.LstdFlags)

	e := echo.New()
	e.Use(middleware.Gzip())
	//e.Use(middleware.Logger())

	e.Validator = &CustomValidator{validator: validator.New()}

	//e.Use(SimulateDelay)

	shimoC := shimo_openapi.NewClient(config.Shimoauth.ClientId, config.Shimoauth.ClientSecret, config.Shimoauth.Username, config.Shimoauth.Password, config.Shimoauth.Scope)
	Log.Println(config, shimoC)

	tgbot, err = tgbotapi.NewBotAPI(config.Telegram.BotToken)
	if err != nil {
		Log.Printf("failed to initialize telegram bot: %v", err)
	}

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
		fileId := "kDQJ6vWgWWwq8r8H"
		opt := shimo_openapi.Opts{"上线版本", 30, "D", " (", time.Minute * 30}
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
		opt := shimo_openapi.Opts{"实时", 100, "H", " ", time.Minute * 30}
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

	e.POST("/hospital/supplies/submissions", func(c echo.Context) error {
		var request Submission
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			// fmt.Println(c.Bind(request))
			fmt.Println(c.Validate(&request))
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}

		tmpl := template.Must(template.New("suppliesSubmission").Parse(`*新的物资需求提交*

- 医院名称：{{.Name}}
- 医院所在地区：{{.Province}} {{.City}} {{.Suburb}}
- 医院详细地址：{{.Address}}
- 医院现每天接待患者数量：{{.Patients}}
- 医院床位数：{{.Beds}}
- 责任人姓名：{{.ContactName}}
- 责任人所在单位或组织：{{.ContactOrg}}
- 责任人联系方式：{{.ContactPhone}}
- 物资需求列表：{{range .Supplies}}
	- 物资名称：{{.Name}}
	  数量单位：{{.Unit}}
	  需求数量：{{.Need}}
	  每日消耗：{{.Daily}}
	  库存数量：{{.Have}}
	  物资要求：{{.Requirements}}
{{end}}
- 可接受的捐物资渠道：{{.Pathways}}
- 现在的物流状况：{{.LogisticStatus}}
- 需求信息数据来源：{{.Source}}
- 需求的官方证明：{{.Proof}}
- 其他备注：{{.Notes}}`))

		buf := bytes.NewBufferString("")
		err = tmpl.Execute(buf, request)
		if err != nil {
			Log.Printf("failed to execute template. invalid data? %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed to execute template. invalid data?")
		}

		go func() {
			Log.Print(notifyAdmins(buf.String()))
		}()

		//d := db.Create(&request)
		//if d.Error != nil {
		//	Log.Printf("create collect_form failed:%v", d.Error)
		//	return echo.NewHTTPError(http.StatusInternalServerError)
		//}
		return c.NoContent(http.StatusNoContent)
	})

	// 返回社区物资需求列表
	e.GET("/community/supplies", func(c echo.Context) error {
		fileId := "qrpCHCDY8t6wccpD"
		opt := shimo_openapi.Opts{"main", 100, "M", " ", time.Minute * 5}
		message, err := shimoC.GetFileWithOpts(fileId, opt)
		if err != nil {
			Log.Printf("failed to get document: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get document")
		}
		refactMessage, err := RefactCommunitySubmissionFromShimoDoc(message)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to refact")
		}
		return c.JSONBlob(http.StatusOK, refactMessage)
	})

	// curl -X POST 'http://localhost:3166/community/supplies/submissions' -H "Content-Type: application/json" -d '{"name":"xyz","age":11, "medicalSupplies": [{"name":"口罩", "[unit":"个", "need": "10", "daily":"1", "have": "0", "requirement": "急需!"}], "province": "湖北", "city": "武汉", "suburb": "汉口", "address": "xx花园"}'
	e.POST("/community/supplies/submissions", func(c echo.Context) error {
		var request CommunitySubmission
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			// fmt.Println(c.Bind(request))
			fmt.Println(c.Validate(&request))
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}
		fileId := "qrpCHCDY8t6wccpD"
		w := shimo_openapi.NewWriteOpts("main", &request)
		err := shimoC.AppendFileFromAPI(fileId, *w)
		if err != nil {
			Log.Printf("[CommunitySubmission] failed, err: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "append file failed")
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.POST("/report", func(c echo.Context) error {
		var request ReportRequest
		if c.Bind(&request) != nil && c.Validate(&request) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request")
		}
		Log.Printf("[report] new report record: %v", spew.Sdump(request))

		tmpl := template.Must(template.New("suppliesSubmission").Parse(`*新的网站信息纠错请求*

- 来源页面名称：{{.Type}}
- 信息纠错请求原因：{{.Cause}}
- 信息原数据：{{.Content}}`))

		buf := bytes.NewBufferString("")
		err = tmpl.Execute(buf, request)
		if err != nil {
			Log.Printf("failed to execute template. invalid data? %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed to execute template. invalid data?")
		}

		go func() {
			_ = notifyAdmins(buf.String())
		}()
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
