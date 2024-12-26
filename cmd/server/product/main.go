package main

import (
	"crypto/tls"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_microservice_backend_api/internal/initialize"
	"go_microservice_backend_api/internal/service_product"
	"log"
	"net/http"
)

// @title           API Document Ecommerce Backend SHOPDEVGO
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   loctodale
// @contact.url    http://www.swagger.io/support
// @contact.email  loctodale.server@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8005
// @BasePath  /api
// @schema http
func main() {
	initialize.Run()
	r := service_product.ProductServiceMain()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := &http.Server{
		Addr:    ":8005",
		Handler: r,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip verification of client certificates
		},
	}
	//r.Run(":8002")
	err := server.ListenAndServeTLS("/certs/cert.crt", "/certs/key.pem")
	if err != nil {
		log.Fatalf("Failed to startss server: %v", err)
	}
}
