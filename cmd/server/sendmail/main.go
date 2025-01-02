package main

import (
	"github.com/gin-gonic/gin"
	"go_microservice_backend_api/internal/initialize"
	"go_microservice_backend_api/internal/service_sendmail"
)

func main() {
	r := *gin.Default()
	initialize.Run()
	service_sendmail.ServiceSendMailMain()

	r.Run(":8003")
}
