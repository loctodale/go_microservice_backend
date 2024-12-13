package main

import (
	"crypto/tls"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_microservice_backend_api/internal/initialize"
	"go_microservice_backend_api/internal/service_shop"
	"log"
	"net/http"
)

func main() {
	initialize.Run()
	r := service_shop.ServiceShopMain()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := &http.Server{
		Addr:    ":8004",
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
