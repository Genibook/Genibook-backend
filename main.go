package main

import (
	"log"
	"net/http"
	api "webscrapper/apis/v2"

	"github.com/gin-gonic/gin"
)

// var defaultLogFormatter = func(param gin.LogFormatterParams) string {
// 	var statusColor, methodColor, resetColor string
// 	if param.IsOutputColor() {
// 		statusColor = param.StatusCodeColor()
// 		methodColor = param.MethodColor()
// 		resetColor = param.ResetColor()
// 	}

// 	if param.Latency > time.Minute {
// 		param.Latency = param.Latency.Truncate(time.Second)
// 	}
// 	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
// 		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
// 		statusColor, param.StatusCode, resetColor,
// 		param.Latency,
// 		param.ClientIP,
// 		methodColor, param.Method, resetColor,
// 		param.Path,
// 		param.ErrorMessage,
// 	)
// }

// func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
// 	formatter := conf.Formatter
// 	if formatter == nil {
// 		formatter = defaultLogFormatter
// 	}

// 	out := conf.Output
// 	if out == nil {
// 		out = gin.DefaultWriter
// 	}

// 	notlogged := conf.SkipPaths

// 	var skip map[string]struct{}

// 	if length := len(notlogged); length > 0 {
// 		skip = make(map[string]struct{}, length)

// 		for _, path := range notlogged {
// 			skip[path] = struct{}{}
// 		}
// 	}

// 	return func(c *gin.Context) {
// 		// Start timer
// 		start := time.Now()
// 		path := c.Request.URL.Path
// 		//raw := c.Request.URL.RawQuery

// 		// Process request
// 		c.Next()

// 		// Log only when path is not being skipped
// 		if _, ok := skip[path]; !ok {
// 			param := gin.LogFormatterParams{
// 				Request: c.Request,
// 				Keys:    c.Keys,
// 			}

// 			// Stop timer
// 			param.TimeStamp = time.Now()
// 			param.Latency = param.TimeStamp.Sub(start)

// 			param.ClientIP = c.ClientIP()
// 			param.Method = c.Request.Method
// 			param.StatusCode = c.Writer.Status()
// 			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

// 			param.BodySize = c.Writer.Size()

// 			// if raw != "" {
// 			// 	path = path + "?" + raw
// 			// }

// 			param.Path = path

// 			fmt.Fprint(out, formatter(param))
// 		}
// 	}
// }
// func Logger() gin.HandlerFunc {
// 	return LoggerWithConfig(gin.LoggerConfig{})
// }

func main() {
	gin.SetMode(gin.ReleaseMode)

	//r := gin.New()
	// r.Use(gin.Recovery())
	// r.Use(Logger())
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

	r.POST("/api/student/", api.MakeHandler(api.StudentHandlerV2))
	r.POST("/api/assignments/", api.MakeHandler(api.AssignmentHandlerV2))
	r.POST("/api/grades/", api.MakeHandler(api.GradesHandlerV2))
	r.POST("/api/schedule/", api.MakeHandler(api.ScheduleAssignmentHandlerV2))
	r.POST("/api/profile/", api.MakeHandler(api.ProfileHandlerV2))
	r.POST("/api/login/", api.MakeHandler(api.LoginHandlerV2))
	r.POST("/api/gpas/", api.MakeHandler(api.GPAshandlerV2))
	r.POST("/api/gpas_his/", api.MakeHandler(api.GPAHistoryHandlerV2))
	r.POST("/api/grade_of_students/", api.MakeHandler(api.StudentGradesHandlerV2))
	r.POST("/api/mps/", api.MakeHandler(api.MpsHandlerV2))
	r.POST("/api/ids/", api.MakeHandler(api.StudentIDHandlerV2))
	r.GET("/api/transcript/", api.MakeHandler(api.TranscriptHandlerV2))

	log.Fatal(r.Run(":6969"))
}
