# 🛡️ Go-HIDS & HIPS: Canlı Siber Savunma ve EDR Ajanı

Go-HIDS, Go (Golang) dilinin yüksek performanslı ve eşzamanlı (concurrent) yapısı kullanılarak geliştirilmiş, uç nokta güvenliğini (endpoint security) sağlamayı amaçlayan hafif bir **HIDS (Host-based Intrusion Detection System)** ve **HIPS (Intrusion Prevention System)** projesidir. 

Sistem, çekirdek (kernel) seviyesindeki dosya bütünlüğü ihlallerini izler, işletim sistemindeki aktif süreçleri (process) denetler ve kara listedeki siber tehdit unsurlarını algıladığı anda **Aktif Müdahale (Active Response)** motoruyla saliseler içinde imha eder.

---

## 🚀 Öne Çıkan Özellikler

- **🔍 FIM (File Integrity Monitoring):** `fsnotify` kütüphanesiyle donatılmış dosya bütünlüğü motoru, `./guvenli_bolge` klasöründeki her türlü dosya oluşturma (Create), içerik manipülasyonu (Write) ve silme (Remove) hareketlerini anlık olarak yakalar.
- **⚙️ Süreç Denetimi (Process Monitoring):** `gopsutil` entegrasyonu sayesinde, işletim sisteminde arka planda çalışan tüm süreçleri (PID) saniyede bir tarayarak siber tehdit ve hacker araçlarını avlar.
- **💥 Aktif Müdahale (EDR Mode / Kill Response):** Sadece alarm üretmekle kalmaz; kara listedeki tehlikeli bir süreç (örn: `mimikatz`, `wireshark` veya test senaryosundaki `notepad.exe`) tetiklendiği an `p.Kill()` komutuyla işletim sistemi seviyesinde süreci otomatik olarak sonlandırır.
- **🔒 Thread-Safe Veri Havuzu (Mutex Architecture):** Birden fazla bağımsız motorun (Goroutine) aynı anda web paneline log göndermesini `sync.Mutex` kilit mekanizmasıyla güvenli hale getirir, veri çakışmalarını (Race Condition) engeller.
- **📊 Merkezi Siber Komuta Merkezi (Canlı Dashboard):** Arka plandaki tüm savunma operasyonlarını, alarmları ve aktif müdahale loglarını Tailwind CSS destekli modern ve karanlık temalı bir web arayüzünde (`:8081`) canlı listeler.

---



## 🛠️ Kurulum ve Çalıştırma

### 1. Gereksinimler
Sistem süreçlerini tarayabilmek ve zararlı işlemleri sonlandırabilmek (Kill) için terminalinizin **Yönetici (Administrator / Root)** yetkilerine sahip olması önerilir.

### 2. Bağımlılıkların İndirilmesi
Proje dizininde terminali açın ve gerekli tüm harici kütüphaneleri otomatik olarak indirmek için şu komutu çalıştırın:

    go mod tidy

### 3. Ajanı Başlatma
Siber savunma ajanını ve canlı komuta merkezini ayağa kaldırmak için:

    go run main.go

Sistem başarıyla tetiklendiğinde terminalinizde şu çıktıları göreceksiniz:

    🚀 GO-HIDS Siber Savunma Sistemi Başlatılıyor...
    📊 [DASHBOARD] Canlı web arayüzü başlatıldı: http://localhost:8081
    [+] Go-HIDS Süreç Takip ve Aktif Müdahale Modülü Aktif!
    🛡️ [HIDS MOTORU] Dosya bütünlüğü izleme sistemi aktif! Klasör: ./guvenli_bolge

---

## 🔬 Canlı Siber Müdahale Testi (EDR Senaryosu)

Sistemin aktif koruma yeteneklerini canlı olarak doğrulamak için aşağıdaki adımları izleyebilirsiniz:

1. Tarayıcınızdan **`http://localhost:8081`** adresine giderek Siber Komuta Merkezi'ni açın. (Koruma Seviyesinin **YÜKSEK - Aktif Koruma** modunda olduğunu doğrulayın).
2. Bilgisayarınızın arama çubuğundan normal bir şekilde **Not Defteri (`notepad.exe`)** uygulamasını çalıştırın.
3. **Müdahale:** Go-HIDS ajanı süreci algıladığı salisede Not Defteri uygulamasını otomatik olarak kapatacaktır.
4. Tarayıcı panelinize döndüğünüzde siber akış tablosuna şu logun düştüğünü göreceksiniz:  
   `💥 AKTİF MÜDAHALE: notepad.exe (PID: XXXX) başarıyla sonlandırıldı ve tehdit engellendi!`

---

## 🔮 Gelecek Planları (Roadmap)

Projenin bir sonraki fazlarında eklenmesi planlanan siber güvenlik modülleri:
- [ ] **Ağ Bağlantısı İzleme (Network Sniffer):** C2 (Command and Control) sunucularına yapılan şüpheli IP/Port bağlantılarını yakalama.
- [ ] **Kullanıcı Yetkilendirmesi (Auth):** Komuta merkezine erişim için JWT tabanlı şifreli giriş ekranı.
- [ ] **Log Rotasyonu (Log Rotation):** `alerts.json` dosyasının boyut sınırını kontrol altında tutmak için otomatik arşivleme mekanizması.

---

## 📄 Lisans
Bu proje eğitim ve siber güvenlik portföy amacıyla geliştirilmiştir. MIT Lisansı altında lisanslanmıştır.
