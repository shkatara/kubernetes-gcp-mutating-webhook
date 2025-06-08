package server

import (
	"log"
	"net/http"
	"strings"

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

	blockImage := []string{"docker.io"}
	foundBlockedImage := false

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
		for i := range blockImage {
			if strings.Contains(container.(Mapping)["image"].(string), blockImage[i]) {
				log.Printf("Container %s has blocked image: %s\n", container.(Mapping)["name"], container.(Mapping)["image"])
				// You can modify the container image here if needed
				//				container.(Mapping)["image"] = "blocked-image"
				foundBlockedImage = true
				break
			}
		}
	}

	if foundBlockedImage {
		log.Printf("Blocked image found in the request: %s\n", admissionReview.Request.Object)
		admissionResponse := types.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: true,
			Status: &types.Status{
				Code:    200,
				Message: "Request Successfully processed",
			},

			PatchType: nil,
			Patch:     []byte(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "blocked-image"}]`),
		}
		c.JSON(http.StatusOK, types.AdmissionReview{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
			Response:   &admissionResponse,
		})
	}

}
