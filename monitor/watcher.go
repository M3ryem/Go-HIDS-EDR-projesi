package monitor

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"

	"go-hids/dashboard" 
)

// StartFileWatcher belirtilen klasörü siber tehditlere karşı canlı izler
func StartFileWatcher(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		_ = os.Mkdir(dirPath, 0755)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("❌ [HATA] İzleme motoru başlatılamadı: %v", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("➕ [ALARM] Yeni Dosya: %s", event.Name)
					dashboard.VeriEkle("CREATE", "Şüpheli dosya oluşturuldu: "+event.Name)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("📝 [ALARM] Değiştirildi: %s", event.Name)
					dashboard.VeriEkle("WRITE", "Dosya içeriği manipüle edildi: "+event.Name)
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("❌ [ALARM] Silindi: %s", event.Name)
					dashboard.VeriEkle("REMOVE", "Kritik dosya silindi: "+event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("❌ [HIDS HATASI]: %v", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatalf("❌ [HATA] Klasör izleme listesine eklenemedi: %v", err)
	}

	log.Printf("🛡️  [HIDS MOTORU] Dosya bütünlüğü izleme sistemi aktif! Klasör: ./%s\n", dirPath)
	select {}
}
