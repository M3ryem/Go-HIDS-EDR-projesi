package fim

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// CalculateSHA256, değişen veya yeni eklenen dosyanın bütünlüğünü kontrol etmek için hash üretir.
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// StartFIM, belirtilen hedef dizini anlık (real-time) olarak izler.
func StartFIM(dirPath string) {
	// Yeni bir dosya izleyici (watcher) oluşturuyoruz
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("[-] Watcher başlatılamadı: %v", err)
	}
	defer watcher.Close()

	// Arka planda olayları dinlemek için bir goroutine başlatıyoruz
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Dosya Oluşturulma Olayı
				if event.Has(fsnotify.Create) {
					// Klasör değilse hash hesapla (HIDS koruması)
					hash, err := CalculateSHA256(event.Name)
					if err == nil {
						fmt.Printf("🟢 [FIM ALARM] YENİ DOSYA OLUŞTURULDU: %s | SHA256: %s\n", event.Name, hash)
					} else {
						fmt.Printf("🟢 [FIM ALARM] YENİ DİZİN/DOSYA: %s\n", event.Name)
					}
				}

				// Dosya Değiştirilme/Yazılma Olayı
				if event.Has(fsnotify.Write) {
					hash, err := CalculateSHA256(event.Name)
					if err == nil {
						fmt.Printf("🟡 [FIM ALARM] DOSYA İÇERİĞİ DEĞİŞTİRİLDİ: %s | NEW SHA256: %s\n", event.Name, hash)
					}
				}

				// Dosya Silinme Olayı
				if event.Has(fsnotify.Remove) {
					fmt.Printf("🔴 [FIM ALARM] KRİTİK DOSYA SİLİNDİ!: %s\n", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[-] İzleme hatası: %v", err)
			}
		}
	}()

	// İzlemek istediğimiz klasörü watcher'a ekliyoruz
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatalf("[-] Dizin izleme listesine eklenemedi: %v", err)
	}

	fmt.Printf("[+] Go-HIDS FIM Modülü Aktif! İzlenen Dizin: %s\n", dirPath)

	// Ana thread'in kapanmaması için boş bir select ile blockluyoruz (sonsuz döngü)
	select {}
}
