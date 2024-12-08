package main

import (
	"crypto/tls"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "go_microservice_backend_api/cmd/swag/docs"
	"go_microservice_backend_api/internal/initialize"
	"log"
	"net/http"
)

// @title           API Document Ecommerce Backend SHOPDEVGO
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   loctodale
// @contact.url    http://www.swagger.io/support
// @contact.email  loctodale.service@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8002
// @BasePath  /v1/2024
// @schema http
func main() {
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := &http.Server{
		Addr:    ":8002",
		Handler: r,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip verification of client certificates
		},
	}
	//r.Run(":8002")
	err := server.ListenAndServeTLS("/certs/cert.crt", "/certs/key.pem")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
