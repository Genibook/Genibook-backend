package main

import (
	"log"
	"net/http"
	api "webscrapper/apis/v2"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLFiles("static/howGPA.html", "static/pp.html")
	//r.Use(api_v1.JsonLoggerMiddleware())

	r.POST("/hi/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	r.GET("/pp/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pp.html", gin.H{})
	})

	r.GET("/howGPA/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "howGPA.html", gin.H{})
	})

	//Key
	/*
		✅ - all functions/code runs with new genesis
		⚠️ - not all functions/code runs with new genesis (checking and not checking)
		❌ - not got to yet.
		✅💥 - code runs with existing genesis pages availble at commit date, however, there are some "cells"/data that don't exist yet.
	*/

	r.POST("/api/student/", api.MakeHandler(api.StudentHandlerV2))                 //✅
	r.POST("/api/assignments/", api.MakeHandler(api.AssignmentHandlerV2))          //✅
	r.POST("/api/grades/", api.MakeHandler(api.GradesHandlerV2))                   //✅
	r.POST("/api/schedule/", api.MakeHandler(api.ScheduleAssignmentHandlerV2))     //✅
	r.POST("/api/profile/", api.MakeHandler(api.ProfileHandlerV2))                 //✅
	r.POST("/api/login/", api.MakeHandler(api.LoginHandlerV2))                     //✅
	r.POST("/api/gpas/", api.MakeHandler(api.GPAshandlerV2))                       //✅
	r.POST("/api/gpas_his/", api.MakeHandler(api.GPAHistoryHandlerV2))             //✅
	r.POST("/api/grade_of_students/", api.MakeHandler(api.StudentGradesHandlerV2)) //❌
	r.POST("/api/mps/", api.MakeHandler(api.MpsHandlerV2))                         //✅
	r.POST("/api/ids/", api.MakeHandler(api.StudentIDHandlerV2))                   //❌
	r.GET("/api/transcript/", api.MakeHandler(api.TranscriptHandlerV2))            //✅

	log.Fatal(r.Run(":6969"))

	//todo: only detect 1 student for now

}

/*


// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = DefaultWriter
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}

*/
