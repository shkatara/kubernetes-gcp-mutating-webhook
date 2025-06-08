package server

import (
	"fmt"
	"log"
	"net/http"

	types "mutating-webhook/internal/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Mapping = map[string]interface{}
type Listing = []interface{}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)
	r.POST("/inject", s.injectHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) injectHandler(c *gin.Context) {

	var admissionReview types.AdmissionReview
	if err := c.BindJSON(&admissionReview); err != nil {
		log.Printf("JSON Parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON parse failed"})
		return
	}
	//log.Printf("Incoming AdmissionReview Request is: %+v\n", admissionReview.Request.Object)
	object := admissionReview.Request.Object.(Mapping)
	containers := object["spec"].(Mapping)["containers"]
	containersList, _ := containers.(Listing)
	for _, container := range containersList {
		fmt.Printf("Container %+v has image -> %+v\n", container.(Mapping)["name"], container.(Mapping)["image"])
	}
}
