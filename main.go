package main

import (
	"log"
	"net/http"
	api_v1 "webscrapper/apis/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/hi/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	r.POST("/student/", api_v1.MakeHandler(api_v1.StudentHandlerV1))
	r.POST("/schedule/", api_v1.MakeHandler(api_v1.ScheduleAssignmentHandlerV1))
	r.POST("/assignments/", api_v1.MakeHandler(api_v1.AssignmentHandlerV1))
	r.POST("/grades/", api_v1.MakeHandler(api_v1.GradesHandlerV1))
	r.POST("/profile/", api_v1.MakeHandler(api_v1.ProfileHandlerV1))
	r.POST("/login/", api_v1.MakeHandler(api_v1.LoginHandlerV1))

	log.Fatal(r.Run(":6969"))
	//log.Fatal(http.ListenAndServe(":6969", nil))

}
