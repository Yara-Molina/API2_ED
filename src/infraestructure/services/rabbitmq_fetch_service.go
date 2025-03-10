package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RabbitMQService struct {
}

func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{}
}

func (s *RabbitMQService) FetchAlerts() ([]map[string]interface{}, error) {
	resp, err := http.Get("http://3.83.5.8:9090/notifications")
	if err != nil {
		log.Println("Error obteniendo alertas del consumidor:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error leyendo respuesta del consumidor:", err)
		return nil, err
	}

	var alerts []map[string]interface{}
	if err := json.Unmarshal(body, &alerts); err != nil {
		log.Println("Error decodificando JSON:", err)
		return nil, err
	}

	return alerts, nil
}
