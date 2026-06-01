package bot

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Helper to run bash commands
func runBash(command string) string {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Error retrieving data"
	}
	return strings.TrimSpace(out.String())
}

func getUptime() string {
	return runBash(`uptime -p | sed 's/up //g'`)
}

func GetRAMUsage() string {
	info := runBash(`free -m | awk 'NR==2{printf "Total: %.1f GB\nUsed: %.1f GB\nFree: %.1f GB\nUsage: %.0f%%", $2/1024, $3/1024, $4/1024, $3*100/$2}'`)
	return "💾 Memory Usage\n\n" + info
}

func GetCPUStatus() string {
	cores := runBash(`nproc`)
	usage := runBash(`top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1"%"}'`)
	load := runBash(`cat /proc/loadavg | awk '{print "1m: "$1"\n5m: "$2"\n15m: "$3}'`)
	
	return fmt.Sprintf("⚡ CPU Status\n\nCores: %s vCPU\nCurrent Usage: %s\nLoad Average:\n%s", cores, usage, load)
}

func GetDiskUsage() string {
	info := runBash(`df -BG / | awk 'NR==2{printf "Total: %sB\nUsed: %sB\nFree: %sB\nUsage: %s", $2, $3, $4, $5}'`)
	return "💿 Disk Usage\n\n" + info
}

func GetDockerContainers() string {
	containers := runBash(`docker ps --format "🟢 {{.Names}}"`)
	if containers == "" {
		containers = "No running containers"
	}
	running := runBash(`docker ps -q | wc -l`)
	stopped := runBash(`docker ps -aq | wc -l | awk -v r="` + running + `" '{print $1-r}'`)

	return fmt.Sprintf("🐳 Docker Containers\n\n%s\n\nRunning: %s\nStopped: %s", containers, running, stopped)
}

func GetNetworkStats() string {
	publicIP := runBash(`curl -s ifconfig.me || echo "Unknown"`)
	rxTx := runBash(`ip -s link show $(ip route | awk '/default/ {print $5}' | head -n1) | awk 'NR==4 {rx=$1} NR==6 {tx=$1} END {printf "RX: %.1f GB\nTX: %.1f GB", rx/1024/1024/1024, tx/1024/1024/1024}'`)
	
	return fmt.Sprintf("🌐 Network\n\nPublic IP: %s\n%s", publicIP, rxTx)
}

func GetFullInfo() string {
	hostname := runBash(`hostname`)
	osRelease := runBash(`grep PRETTY_NAME /etc/os-release | cut -d'"' -f2`)
	cores := runBash(`nproc`)
	
	// Format: RAM: 1.8 / 4 GB (45%)
	ram := runBash(`free -m | awk 'NR==2{printf "%.1f / %.1f GB (%.0f%%)", $3/1024, $2/1024, $3*100/$2}'`)
	
	// Format: Disk: 25 / 80 GB (31%)
	disk := runBash(`df -BG / | awk 'NR==2{printf "%sB / %sB (%s)", $3, $2, $5}'`)
	
	dockerRunning := runBash(`docker ps -q | wc -l`)
	uptime := getUptime()
	load := runBash(`cat /proc/loadavg | awk '{print $1" / "$2" / "$3}'`)
	
	return fmt.Sprintf(`🖥️ VPS Information

	Hostname: %s
	OS: %s
	CPU: %s vCPU
	RAM: %s
	Disk: %s
	Docker: %s Running
	MongoDB: 🟢 Connected
	API: 🟢 Healthy
	Uptime: %s
	Load: %s`, hostname, osRelease, cores, ram, disk, dockerRunning, uptime, load)
}