package main

import (
	"log"
	infrastructure "notifications/src/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	deps, err := infrastructure.NewDependencies()
	if err != nil {
		log.Fatal("Error al inicializar dependencias:", err)
	}

	r := gin.Default()

	infrastructure.RegisterRoutes(r, deps.ProcessLoanUseCase)

	r.Run(":8082")

}
