package config

// BlacklistedProcesses, sistemde alarm üretilecek şüpheli süreçlerin listesidir.
var BlacklistedProcesses = []string{
	"mimikatz",
	"mimikatz.exe",
	"nc",
	"nc.exe",
	"netcat",
	"powershell.exe",
	"cmd.exe",
	"notepad.exe",
}

// BlacklistedPorts, C2 sunucularının veya reverse shell'lerin sıkça kullandığı şüpheli portlardır.
var BlacklistedPorts = []uint32{
	4444, // Metasploit / Standart reverse shell portu
	9001, // Netcat / Zararlı bağlantı favorisi
	8443, // Bazı popüler C2 framework'leri (Sliver, Havoc vb.)
}
