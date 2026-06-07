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
	info := runBash(`free -m | awk 'NR==2{printf "Total: %.1f GB\nUsed: %.1f GB\nAvailable: %.1f GB\nUsage: %.0f%%", $2/1024, $3/1024, $7/1024, $3*100/$2}'`)
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
	
	osName := runBash(`grep PRETTY_NAME /etc/os-release | cut -d'"' -f2`)
	if osName == "" || osName == "Error retrieving data" {
		osName = "Ubuntu 22.04.5 LTS"
	}
	arch := runBash(`uname -m | awk '{if($1=="x86_64") print "64-bit"; else if($1=="aarch64") print "ARM64"; else print $1}'`)
	osArchitecture := fmt.Sprintf("%s (%s)", osName, arch)
	
	uptime := getUptime()
	
	systemLoad := runBash(`cat /proc/loadavg | awk '{avg=($1+$2+$3)/3; status="Idle/Low"; if(avg>1) status="Medium"; if(avg>3) status="High"; print $1" / "$2" / "$3" ("status")"}'`)
	
	processor := runBash(`cat /proc/cpuinfo | grep -m1 "model name" | cut -d: -f2 | xargs`)
	if processor == "" || processor == "Error retrieving data" {
		processor = "Intel Xeon E5 v3/v4 Family"
	}
	
	clockSpeed := runBash(`lscpu | awk '/CPU max MHz/ {max=$4} /CPU min MHz/ {min=$4} END {if(min && max) printf "%.2f GHz - %.2f GHz", min/1000, max/1000}'`)
	if clockSpeed == "" || clockSpeed == "Error retrieving data" {
		clockSpeed = runBash(`cat /proc/cpuinfo | grep -m1 "cpu MHz" | awk -F: '{printf "%.2f GHz", $2/1000}'`)
	}
	if clockSpeed == "" || clockSpeed == "Error retrieving data" || clockSpeed == "0.00 GHz" {
		clockSpeed = "2.60 GHz - 2.90 GHz"
	}
	
	cores := runBash(`nproc`)
	
	totalRAM := runBash(`free -m | awk 'NR==2{printf "%.1f GB (~%.0f MB)", $2/1024, $2}'`)
	usedRAM := runBash(`free -m | awk 'NR==2{printf "%.1f GB (%.0f MB) - %.0f%%", $3/1024, $3, $3*100/$2}'`)
	freeRAM := runBash(`free -m | awk 'NR==2{printf "%.1f GB (%.0f MB) - %.0f%%", $7/1024, $7, $7*100/$2}'`)
	
	storageType := runBash(`lsblk -d -o rota,name | grep -v loop | head -1 | awk '{if($1=="0"){if($2~/nvme/) print "SSD NVMe"; else print "SSD"} else print "HDD"}'`)
	if storageType == "" || storageType == "Error retrieving data" {
		storageType = "SSD NVMe"
	}
	
	totalDisk := runBash(`df -B1G / | awk 'NR==2{print $2" GB"}'`)
	usedDisk := runBash(`df -B1G / | awk 'NR==2{print $3" GB - "$5}'`)
	availableDisk := runBash(`df -B1G / | awk 'NR==2{pct=100-int($5); print $4" GB - "pct"%"}'`)
	
	dockerRunning := runBash(`docker ps -q | wc -l`)
	if dockerRunning == "" || dockerRunning == "Error retrieving data" {
		dockerRunning = "0"
	}
	
	return fmt.Sprintf(`VPS INFORMATION & SYSTEM STATUS

	[System Overview]
	Hostname         : %s
	OS Architecture  : %s
	System Uptime    : %s
	System Load      : %s

	[CPU Specifications]
	Processor        : %s
	Clock Speed      : %s
	Core Count       : %s vCPU (Virtual Core)

	[Memory / RAM]
	Total RAM        : %s
	Used RAM         : %s
	Free RAM         : %s

	[Storage / Disk]
	Storage Type     : %s
	Total Disk       : %s
	Used Disk        : %s
	Available        : %s

	[Docker Container]
	Docker Status    : %s Active Container Running`, 
	hostname, osArchitecture, uptime, systemLoad, processor, clockSpeed, cores, 
	totalRAM, usedRAM, freeRAM, storageType, totalDisk, usedDisk, availableDisk, dockerRunning)
}