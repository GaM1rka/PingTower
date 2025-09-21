package internal

import (
	"api_service/configs"
	"net/http"
	"os"
	"path/filepath"
)

// SwaggerHandler handles Swagger UI and specification requests
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new SwaggerHandler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// ServeSwaggerSpec serves the Swagger specification file
func (s *SwaggerHandler) ServeSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		configs.APILogger.Printf("Failed to get working directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Construct the path to swagger.yaml file
	swaggerPath := filepath.Join(pwd, "swagger.yaml")
	
	// Check if file exists
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		// Try alternative path
		swaggerPath = filepath.Join(pwd, "..", "..", "swagger.yaml")
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			configs.APILogger.Printf("Swagger file not found at: %s or %s", filepath.Join(pwd, "swagger.yaml"), filepath.Join(pwd, "..", "..", "swagger.yaml"))
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}
	}

	// Set the appropriate content type for YAML
	w.Header().Set("Content-Type", "application/x-yaml")
	
	// Serve the swagger.yaml file
	http.ServeFile(w, r, swaggerPath)
}

// ServeSwaggerUI serves the Swagger UI HTML page
func (s *SwaggerHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		configs.APILogger.Printf("Failed to get working directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Construct the path to swagger-ui.html file
	uiPath := filepath.Join(pwd, "swagger-ui.html")
	
	// Check if file exists
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		// Try alternative path
		uiPath = filepath.Join(pwd, "..", "..", "swagger-ui.html")
		if _, err := os.Stat(uiPath); os.IsNotExist(err) {
			configs.APILogger.Printf("Swagger UI file not found at: %s or %s", filepath.Join(pwd, "swagger-ui.html"), filepath.Join(pwd, "..", "..", "swagger-ui.html"))
			http.Error(w, "Swagger UI file not found", http.StatusNotFound)
			return
		}
	}

	// Set the appropriate content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Serve the swagger-ui.html file
	http.ServeFile(w, r, uiPath)
}