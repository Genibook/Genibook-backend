// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	api "webscrapper/apis/v2"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	gin.SetMode(gin.ReleaseMode)

// 	//r := gin.New()
// 	// r.Use(gin.Recovery())
// 	// r.Use(Logger())
// 	r := gin.Default()

// 	r.LoadHTMLFiles("./static/howGPA.html", "./static/pp.html")
// 	//r.Use(api_v1.JsonLoggerMiddleware())

// 	r.POST("/hi/", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})

// 	})

// 	r.GET("/pp/", func(ctx *gin.Context) {
// 		ctx.HTML(http.StatusOK, "pp.html", gin.H{})
// 	})

// 	r.GET("/howGPA/", func(ctx *gin.Context) {
// 		ctx.HTML(http.StatusOK, "howGPA.html", gin.H{})
// 	})

// 	r.POST("/api/student/", api.MakeHandler(api.StudentHandlerV2))
// 	r.POST("/api/assignments/", api.MakeHandler(api.AssignmentHandlerV2))
// 	r.POST("/api/grades/", api.MakeHandler(api.GradesHandlerV2))
// 	r.POST("/api/schedule/", api.MakeHandler(api.ScheduleAssignmentHandlerV2))
// 	r.POST("/api/profile/", api.MakeHandler(api.ProfileHandlerV2))
// 	r.POST("/api/login/", api.MakeHandler(api.LoginHandlerV2))
// 	r.POST("/api/gpas/", api.MakeHandler(api.GPAshandlerV2))
// 	r.POST("/api/gpas_his/", api.MakeHandler(api.GPAHistoryHandlerV2))
// 	r.POST("/api/grade_of_students/", api.MakeHandler(api.StudentGradesHandlerV2))
// 	r.POST("/api/mps/", api.MakeHandler(api.MpsHandlerV2))
// 	r.POST("/api/ids/", api.MakeHandler(api.StudentIDHandlerV2))
// 	r.GET("/api/transcript/", api.MakeHandler(api.TranscriptHandlerV2))

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "6969"
// 	}

//		log.Fatal(r.Run(":" + port))
//	}
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": 1})
	})

	router.Run(":" + port)
}
