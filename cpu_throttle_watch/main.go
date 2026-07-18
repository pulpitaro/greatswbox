package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	t := flag.Float64("t", 80.0, "Threshold")
	cn := flag.Bool("cn", false, "Disable clear (clear none)")
	flag.Parse()

	zones, _ := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if len(zones) == 0 { fmt.Println("No thermal zones"); return }
	paths, _ := filepath.Glob("/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq")

	if !*cn { fmt.Print("\033[H\033[2J") }
	fmt.Printf("🔥 Watcher (T: %.1f°C) | Ctrl+C to stop\n", *t)

	for range time.Tick(time.Second) {
		d, _ := os.ReadFile(zones[0])
		temp, _ := strconv.ParseFloat(strings.TrimSpace(string(d)), 64)
		temp /= 1000

		var sum float64
		for _, p := range paths {
			d, _ := os.ReadFile(p)
			f, _ := strconv.ParseFloat(strings.TrimSpace(string(d)), 64)
			sum += f / 1000
		}

		if temp >= *t {
			fmt.Printf("\n[%s] 🚨 ALERT! %.1f°C | Avg: %.2f MHz", time.Now().Format("15:04:05"), temp, sum/float64(len(paths)))
		} else {
			fmt.Printf("\rStatus: Normal | Temp: %.1f°C | Avg: %.1f MHz ", temp, sum/float64(len(paths)))
		}
	}
}