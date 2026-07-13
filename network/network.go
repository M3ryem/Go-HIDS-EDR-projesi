package network

import (
	"fmt"
	"go-hids/config"
	"go-hids/logger"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// StartNetworkMonitoring, sistemdeki aktif ağ bağlantılarını periyodik olarak denetler.
func StartNetworkMonitoring() {
	fmt.Println("[+] Go-HIDS Ağ Bağlantısı Takip Modülü Aktif!")

	for {
		time.Sleep(2 * time.Second) // Ağı sürekli tarayıp sistemi yormamak için 2 saniyede bir kontrol et

		// Sistemdeki tüm aktif ağ bağlantılarını (TCP/UDP) alıyoruz
		connections, err := net.Connections("all")
		if err != nil {
			continue
		}

		for _, conn := range connections {
			// Sadece başarılı şekilde kurulmuş (ESTABLISHED) dış bağlantıları inceliyoruz
			if conn.Status == "ESTABLISHED" {

				// Bağlantının gittiği uzak portu kara listemizle karşılaştırıyoruz
				for _, badPort := range config.BlacklistedPorts {
					if conn.Raddr.Port == badPort {
						msg := fmt.Sprintf("Şüpheli port bağlantısı! Hedef: %s:%d | PID: %d", conn.Raddr.IP, conn.Raddr.Port, conn.Pid)
						fmt.Printf("🚨 [NETWORK ALARM] %s\n", msg)

						// YENİ: JSON Log dosyamıza kaydetmek için çağırıyoruz
						logger.LogAlert("NETWORK", msg, "CRITICAL")
					}
				}
			}
		}
	}
}
