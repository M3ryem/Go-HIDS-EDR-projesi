package main

import (
	"fmt"
	"go-hids/dashboard" // ✅ Doğru modül adı olan go-hids yapıldı
	"go-hids/monitor"   // ✅ Doğru modül adı olan go-hids yapıldı
	"go-hids/process"   // ✅ Doğru modül adı olan go-hids yapıldı
)

func main() {
	fmt.Println("🚀 GO-HIDS Siber Savunma Sistemi Başlatılıyor...")

	// 1. Web Dashboard'u arka planda 8081 portu ile başlatıyoruz
	go func() {
		dashboard.StartWebDashboard("8081")
	}()

	// 2. Gerçek Süreç İzleme Motorunu arka planda başlatıyoruz
	go func() {
		process.StartProcessMonitoring()
	}()

	// 3. Dosya İzleme Motorunu ana akışta başlatıyoruz (Klasörü dinler)
	monitor.StartFileWatcher("guvenli_bolge")
}
