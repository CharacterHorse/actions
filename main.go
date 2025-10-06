package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type StatusResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type SystemMetrics struct {
	CPUUsage    string `json:"cpu_usage"`
	CPUCount    int    `json:"cpu_count"`
	MemoryUsage string `json:"memory_usage"`
	MemoryUsed  string `json:"memory_used"`
	MemoryFree  string `json:"memory_free"`
	MemoryTotal string `json:"total_memory"`
	SwapUsage   string `json:"swap_usage"`
	SwapUsed    string `json:"swap_used"`
	SwapFree    string `json:"swap_free"`
	SwapTotal   string `json:"swap_total"`
	DiskUsage   string `json:"disk_usage"`
	DiskUsed    string `json:"disk_used"`
	DiskFree    string `json:"disk_free"`
	DiskTotal   string `json:"disk_total"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status: "OK",
		Time:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func getCpuMetrics(w http.ResponseWriter, r *http.Request) {
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		http.Error(w, "Could not get CPU usage", http.StatusInternalServerError)
		log.Fatal("Error getting CPU usage:", err)
		return
	}

	cpuUsage[0] = math.Round(cpuUsage[0]*100) / 100

	cpuUsageStr := strconv.FormatFloat(cpuUsage[0], 'f', 2, 64) + "%"

	cpuCount, err := cpu.Counts(false)
	if err != nil {
		http.Error(w, "Could not get CPU count", http.StatusInternalServerError)
		log.Fatal("Error getting CPU count:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"cpu_usage": cpuUsageStr,
		"cpu_count": strconv.Itoa(cpuCount),
	})
}

func getMemoryMetrics(w http.ResponseWriter, r *http.Request) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		http.Error(w, "Could not get memory usage", http.StatusInternalServerError)
		log.Fatal("Error getting memory usage:", err)
		return
	}

	totalMemory := memInfo.Total / (1024 * 1024 * 1024)
	memoryUsed := memInfo.Used / (1024 * 1024 * 1024)
	memoryFree := memInfo.Free / (1024 * 1024 * 1024)
	swapTotal := memInfo.SwapTotal / (1024 * 1024 * 1024)
	swapUsed := memInfo.SwapCached / (1024 * 1024 * 1024)
	swapFree := memInfo.SwapFree / (1024 * 1024 * 1024)

	memoryUsageStr := strconv.FormatFloat(memInfo.UsedPercent, 'f', 2, 64) + "%"
	memoryTotalStr := strconv.FormatFloat(float64(totalMemory), 'f', 2, 64) + " GB"
	memoryUsedStr := strconv.FormatFloat(float64(memoryUsed), 'f', 2, 64) + " GB"
	memoryFreeStr := strconv.FormatFloat(float64(memoryFree), 'f', 2, 64) + " GB"
	swapUsedStr := strconv.FormatFloat(float64(swapUsed), 'f', 2, 64) + " GB"
	swapFreeStr := strconv.FormatFloat(float64(swapFree), 'f', 2, 64) + " GB"
	swapTotalStr := strconv.FormatFloat(float64(swapTotal), 'f', 2, 64) + " GB"

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"memory_usage": memoryUsageStr,
		"memory_used":  memoryUsedStr,
		"memory_free":  memoryFreeStr,
		"memory_total": memoryTotalStr,
		"swap_used":    swapUsedStr,
		"swap_free":    swapFreeStr,
		"swap_total":   swapTotalStr,
	})
}

func getDiskMetrics(w http.ResponseWriter, r *http.Request) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		http.Error(w, "Could not get disk usage", http.StatusInternalServerError)
		log.Fatal("Error getting disk usage:", err)
		return
	}

	diskUsageStr := strconv.FormatFloat(math.Round(diskUsage.UsedPercent*100)/100, 'f', 2, 64) + "%"
	diskTotal := diskUsage.Total / (1024 * 1024 * 1024)
	diskUsed := diskUsage.Used / (1024 * 1024 * 1024)
	diskFree := diskUsage.Free / (1024 * 1024 * 1024)
	diskTotalStr := strconv.FormatFloat(float64(diskTotal), 'f', 2, 64) + " GB"
	diskUsedStr := strconv.FormatFloat(float64(diskUsed), 'f', 2, 64) + " GB"
	diskFreeStr := strconv.FormatFloat(float64(diskFree), 'f', 2, 64) + " GB"

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"disk_usage": diskUsageStr,
		"disk_used":  diskUsedStr,
		"disk_free":  diskFreeStr,
		"disk_total": diskTotalStr,
	})
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/metrics", getCpuMetrics)
	http.HandleFunc("/metrics/mem", getMemoryMetrics)
	http.HandleFunc("/metrics/disk", getDiskMetrics)

	http.HandleFunc("/metrics/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusSeeOther)
	})

	http.HandleFunc("/metrics/mem/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics/mem", http.StatusSeeOther)
	})

	http.HandleFunc("/metrics/disk/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics/disk", http.StatusSeeOther)
	})

	log.Printf("Starting server on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
