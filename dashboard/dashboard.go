package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// PanelVerisi web arayüzüne göndereceğimiz istatistik yapısı
type PanelVerisi struct {
	TotalAlerts int      `json:"total_alerts"`
	CreateCount int      `json:"create_count"`
	WriteCount  int      `json:"write_count"`
	DeleteCount int      `json:"delete_count"`
	RecentLogs  []string `json:"recent_logs"`
}

// Güvenli veri havuzumuz
var (
	SystemStats PanelVerisi
	statsMutex  sync.Mutex
)

// Alert yapısı
type Alert struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Details   string `json:"details"`
}

// VeriEkle fonksiyonunu monitor paketi çağıracak
func VeriEkle(islemTipi string, mesaj string) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	SystemStats.TotalAlerts++
	zaman := time.Now().Format("2006-01-02 15:04:05")
	logMesaji := "[" + zaman + "] " + mesaj

	switch islemTipi {
	case "CREATE":
		SystemStats.CreateCount++
	case "WRITE":
		SystemStats.WriteCount++
	case "REMOVE":
		SystemStats.DeleteCount++
	}

	if len(SystemStats.RecentLogs) >= 10 {
		SystemStats.RecentLogs = SystemStats.RecentLogs[1:]
	}
	SystemStats.RecentLogs = append(SystemStats.RecentLogs, logMesaji)

	// Alarmları alerts.json dosyasına yazıyoruz
	yeniAlarm := Alert{
		Timestamp: zaman,
		Type:      "PROCESS",
		Severity:  "CRITICAL", // Aktif engelleme alarmları kritiktir
		Message:   mesaj,
		Details:   "Go-HIDS EDR Aktif Yanıt Motoru tarafından işletim sistemi seviyesinde engelleme tetiklendi.",
	}

	jsonBytes, err := json.Marshal(yeniAlarm)
	if err == nil {
		f, err := os.OpenFile("alerts.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			_, _ = f.WriteString(string(jsonBytes) + "\n")
		}
	}
}

// StartWebDashboard Web sunucusunu başlatır
func StartWebDashboard(port string) {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/alerts", handleAlertsAPI)

	fmt.Printf("📊 [DASHBOARD] Canlı web arayüzü başlatıldı: http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("❌ [DASHBOARD] Web sunucu başlatılamadı: %v\n", err)
	}
}

func handleAlertsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := os.Stat("alerts.json"); os.IsNotExist(err) {
		w.Write([]byte("[]"))
		return
	}

	fileBytes, err := os.ReadFile("alerts.json")
	if err != nil {
		w.Write([]byte("[]"))
		return
	}

	var alerts []Alert
	content := string(fileBytes)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, ",") {
			line = strings.TrimSuffix(line, ",")
		}

		var alert Alert
		if err := json.Unmarshal([]byte(line), &alert); err == nil {
			alerts = append(alerts, alert)
		}
	}

	if len(alerts) == 0 && len(fileBytes) > 0 {
		_ = json.Unmarshal(fileBytes, &alerts)
	}

	jsonBytes, err := json.Marshal(alerts)
	if err != nil {
		w.Write([]byte("[]"))
		return
	}

	w.Write(jsonBytes)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go-HIDS | Siber Güvenlik Komuta Merkezi</title>
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
    <style>
        body { background-color: #0b0f19; font-family: 'Courier New', Courier, monospace; }
        .cyber-card { background: rgba(17, 24, 39, 0.7); border: 1px solid #1f2937; backdrop-filter: blur(10px); }
        .neon-text-red { text-shadow: 0 0 10px rgba(239, 68, 68, 0.5); }
        .neon-border-blue { box-shadow: 0 0 15px rgba(59, 130, 246, 0.2); border-color: #3b82f6; }
        .neon-border-red { box-shadow: 0 0 15px rgba(239, 68, 68, 0.2); border-color: #ef4444; }
    </style>
</head>
<body class="text-gray-100 min-h-screen p-6">
    <div class="max-w-7xl mx-auto">
        
        <header class="flex justify-between items-center mb-8 border-b border-gray-800 pb-4">
            <div>
                <h1 class="text-3xl font-bold text-blue-500 tracking-wider">⚡ GO-HIDS MONITOR v1.5</h1>
                <p class="text-xs text-gray-400 mt-1">Host-based Intrusion Prevention & Endpoint Response Agent</p>
            </div>
            <div class="flex items-center space-x-3">
                <span class="relative flex h-3 w-3">
                    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                    <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                </span>
                <span class="text-sm font-semibold text-green-400 tracking-widest">SİSTEM CANLI</span>
            </div>
        </header>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="cyber-card p-6 rounded-lg neon-border-blue">
                <h3 class="text-sm text-gray-400 uppercase tracking-wider">Toplam Tehdit Alarmı</h3>
                <p id="total-alerts-count" class="text-4xl font-bold text-red-500 mt-2 neon-text-red">0</p>
            </div>
            <div class="cyber-card p-6 rounded-lg border-gray-800">
                <h3 class="text-sm text-gray-400 uppercase tracking-wider">Aktif Modüller</h3>
                <div class="flex space-x-2 mt-3">
                    <span class="bg-blue-950 text-blue-400 px-2 py-1 text-xs rounded border border-blue-800">FILE_HIDS</span>
                    <span class="bg-red-950 text-red-400 px-2 py-1 text-xs rounded border border-red-800 animate-pulse">ACTIVE_RESPONSE_EDR</span>
                </div>
            </div>
            <div class="cyber-card p-6 rounded-lg neon-border-red">
                <h3 class="text-sm text-gray-400 uppercase tracking-wider">Koruma Seviyesi</h3>
                <p class="text-xl font-bold text-red-500 mt-2">YÜKSEK (Aktif Koruma)</p>
                <span class="text-[10px] text-green-400 block mt-1">✓ Tehditleri Otomatik Engelleme Aktif</span>
            </div>
        </div>

        <div class="cyber-card rounded-lg overflow-hidden border-gray-800">
            <div class="bg-gray-900 px-6 py-4 border-b border-gray-800 flex justify-between items-center">
                <h2 class="text-lg font-bold text-gray-300 tracking-wide">🚨 CANLI TEHDİT VE ALARM AKIŞI</h2>
                <span class="text-xs text-gray-500">Her 3 saniyede bir güncellenir</span>
            </div>
            <div class="p-6">
                <div class="overflow-x-auto">
                    <table class="w-full text-left border-collapse">
                        <thead>
                            <tr class="border-b border-gray-800 text-gray-400 text-xs uppercase tracking-wider">
                                <th class="pb-3 w-1/5">Zaman Damgası</th>
                                <th class="pb-3 w-1/6">Modül</th>
                                <th class="pb-3 w-1/12">Seviye</th>
                                <th class="pb-3">Detay ve Açıklama</th>
                            </tr>
                        </thead>
                        <tbody id="alerts-table-body" class="text-sm divide-y divide-gray-900">
                        </tbody>
                    </table>
                    <div id="no-alerts-msg" class="text-center py-12 text-gray-500 hidden">
                        Sistemde henüz şüpheli bir aktivite algılanmadı. Temiz!
                    </div>
				</div>
			</div>
		</div>

	</div>

	<script>
		function fetchAlerts() {
			fetch('/api/alerts')
				.then(response => response.json())
				.then(data => {
					if (!Array.isArray(data)) {
						data = [];
					}
					data.reverse();

					const tableBody = document.getElementById('alerts-table-body');
					const totalCountEl = document.getElementById('total-alerts-count');
					const noAlertsMsg = document.getElementById('no-alerts-msg');

					totalCountEl.innerText = data.length;

					if(data.length === 0) {
						tableBody.innerHTML = '';
						noAlertsMsg.classList.remove('hidden');
						return;
					}
					noAlertsMsg.classList.add('hidden');

					let rowsHtml = '';
					data.forEach(alert => {
						if(!alert.timestamp) return;
						
						let severityBadge = '';
						if(alert.severity === 'CRITICAL' || alert.severity === 'HIGH') {
							severityBadge = '<span class="bg-red-950 text-red-400 border border-red-800 px-2 py-0.5 rounded text-xs font-bold animate-pulse">CRITICAL</span>';
						} else {
							severityBadge = '<span class="bg-yellow-950 text-yellow-400 border border-yellow-800 px-2 py-0.5 rounded text-xs">' + alert.severity + '</span>';
						}

						let typeBadge = '';
						if(alert.type === 'PROCESS') {
							typeBadge = '<span class="text-purple-400 font-bold">[PROCESS]</span>';
						} else if(alert.type === 'NETWORK') {
							typeBadge = '<span class="text-cyan-400 font-bold">[NETWORK]</span>';
						} else {
							typeBadge = '<span class="text-gray-400">[' + alert.type + ']</span>';
						}

						rowsHtml += '<tr class="hover:bg-gray-950/50 transition">' +
							'<td class="py-4 text-xs text-gray-400">' + alert.timestamp + '</td>' +
							'<td class="py-4 text-xs">' + typeBadge + '</td>' +
							'<td class="py-4 text-xs">' + severityBadge + '</td>' +
							'<td class="py-4">' +
								'<div class="text-gray-200 font-semibold">' + alert.message + '</div>' +
								'<div class="text-xs text-gray-500 mt-1 bg-black/30 p-2 rounded border border-gray-900">' + alert.details + '</div>' +
							'</td>' +
						'</tr>';
					});
					tableBody.innerHTML = rowsHtml;
				})
				.catch(err => console.error("Arayüz veri çekme hatası:", err));
		}

		fetchAlerts();
		setInterval(fetchAlerts, 3000);
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
