package server

import (
	"fmt"
	"log"
	auth "mutating-webhook/internal/gcp-auth"
	types "mutating-webhook/internal/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

	//log.Printf("Incoming AdmissionReview Request is: %+v\n", admissionReview)
	object := admissionReview.Request.Object.(types.Mapping)
	containers := object["spec"].(types.Mapping)["containers"]
	containersList, _ := containers.(types.Listing)
	for _, container := range containersList {
		for i := range blockImage {
			if strings.Contains(container.(types.Mapping)["image"].(string), blockImage[i]) {
				log.Printf("Container %s has blocked image: %s\n", container.(types.Mapping)["name"], container.(types.Mapping)["image"])
				// You can modify the container image here if needed
				//				container.(Mapping)["image"] = "blocked-image"
				foundBlockedImage = true
				break
			}
		}
	}

	if foundBlockedImage {
		log.Printf("Blocked image found in the request: %s\n", admissionReview.Request.UID)
		admissionResponse := types.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: true,
			Status: &types.Status{
				Code:    200,
				Message: "Request Successfully processed",
			},

			PatchType: func(s string) *string { return &s }("JSONPatch"),
			Patch:     []byte(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "blocked-image"}]`),
		}
		c.JSON(http.StatusOK, types.AdmissionReview{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
			Response:   &admissionResponse,
		})
	}

	tokenResponse, err := auth.MakeSTSRequest()
	if err != nil {
		fmt.Errorf("failed to create request: %w", err)
	}

	fmt.Println("Token Response:", tokenResponse)

}
