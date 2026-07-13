package process

import (
	"fmt"
	"strings"
	"time"

	"go-hids/config"
	"go-hids/dashboard"
	"go-hids/logger"

	"github.com/shirou/gopsutil/v3/process"
)

// StartProcessMonitoring, sistemde yeni başlayan süreçleri periyodik olarak denetler ve zararlıları engeller.
func StartProcessMonitoring() {
	fmt.Println("[+] Go-HIDS Süreç Takip ve Aktif Müdahale Modülü Aktif!")

	trackedProcesses := make(map[int32]bool)

	initialProcesses, _ := process.Processes()
	for _, p := range initialProcesses {
		trackedProcesses[p.Pid] = true
	}

	for {
		time.Sleep(1 * time.Second)

		currentProcesses, err := process.Processes()
		if err != nil {
			continue
		}

		for _, p := range currentProcesses {
			if !trackedProcesses[p.Pid] {
				name, err := p.Name()
				if err != nil {
					continue
				}

				isMalicious := false
				for _, badProc := range config.BlacklistedProcesses {
					loweredName := strings.ToLower(name)
					loweredBadProc := strings.ToLower(badProc)

					if loweredName == loweredBadProc || strings.TrimSuffix(loweredName, ".exe") == loweredBadProc {
						isMalicious = true
						break
					}
				}

				if isMalicious {
					// =======================================================================
					// 🔥 [AKTİF YANIT / HIPS MODÜLÜ DEVREDE]: ZARARLIYI CANLI SONLANDIRMA
					// =======================================================================
					msg := fmt.Sprintf("Zararlı süreç tespit edildi: %s | PID: %d -> İMHA EDİLİYOR...", name, p.Pid)
					fmt.Printf("🚨 [PROCESS ALARM] %s\n", msg)

					// İşletim sistemine kill sinyali göndererek zararlıyı anında kapatıyoruz
					killErr := p.Kill()

					var dashboardMesaj string
					if killErr != nil {
						fmt.Printf("❌ [HIPS HATASI] %s sonlandırılamadı: %v\n", name, killErr)
						dashboardMesaj = fmt.Sprintf("Kritik Hata: Zararlı süreç (%s) sonlandırılamadı! Manuel müdahale gerek.", name)
						logger.LogAlert("PROCESS", dashboardMesaj, "CRITICAL")
					} else {
						fmt.Printf("💥 [ACTIVE RESPONSE] %s başarıyla etkisiz hale getirildi (Killed)!\n", name)
						dashboardMesaj = fmt.Sprintf("💥 AKTİF MÜDAHALE: %s (PID: %d) başarıyla sonlandırıldı ve tehdit engellendi!", name, p.Pid)
						logger.LogAlert("PROCESS", dashboardMesaj, "HIGH")
					}

					// Yakalanan bu aktif müdahale olayını web paneline fırlatıyoruz
					dashboard.VeriEkle("WRITE", dashboardMesaj)
				}

				// Yeni süreci takip listemize ekliyoruz
				trackedProcesses[p.Pid] = true
			}
		}
	}
}
