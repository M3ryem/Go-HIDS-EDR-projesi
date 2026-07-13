package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AlertPayload, JSON formatında kaydedilecek alarm yapısıdır.
type AlertPayload struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`     // PROCESS, NETWORK, FIM
	Message   string `json:"message"`  // Alarm detayı
	Severity  string `json:"severity"` // HIGH, MEDIUM, CRITICAL
}

// LogAlert, gelen alarmı hem ekrana basar hem de alerts.json dosyasına yazar.
func LogAlert(alertType, message, severity string) {
	payload := AlertPayload{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Type:      alertType,
		Message:   message,
		Severity:  severity,
	}

	// JSON formatına dönüştür
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[-] Log dönüştürme hatası: %v\n", err)
		return
	}

	// alerts.json dosyasını aç (yoksa oluştur, varsa sonuna ekle)
	file, err := os.OpenFile("alerts.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[-] Log dosyası açılamadı: %v\n", err)
		return
	}
	defer file.Close()

	// Dosyaya yaz ve sonuna yeni satır ekle
	if _, err := file.WriteString(string(jsonData) + "\n"); err != nil {
		fmt.Printf("[-] Log yazma hatası: %v\n", err)
	}
}
