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

type Core struct { id int; f int64; t, p string }

func main() {
	once := flag.Bool("c", false, "Once")
	flag.Parse()

	// 1. Scan CPU topology and detect core types
	dirs, _ := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*")
	var cores []Core
	var max, min int64 = 0, 9999999999

	for _, d := range dirs {
		id, _ := strconv.Atoi(filepath.Base(d)[3:])
		b, err := os.ReadFile(d + "/cpufreq/cpuinfo_max_freq")
		if err != nil { continue } // Skip HT threads
		f, _ := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64)
		if f > max { max = f }; if f < min { min = f }
		cores = append(cores, Core{id: id, f: f})
	}

	// 2. Map hwmon thermal sensors (Core temp)
	hw, _ := filepath.Glob("/sys/class/hwmon/hwmon*/name")
	var hwDir string
	for _, h := range hw {
		n, _ := os.ReadFile(h)
		if strings.TrimSpace(string(n)) == "coretemp" { hwDir = filepath.Dir(h); break }
	}
	lbls, _ := filepath.Glob(hwDir + "/input*_label")
	lblMap := make(map[string]string)
	for _, l := range lbls {
		b, _ := os.ReadFile(l)
		lblMap[strings.TrimSpace(string(b))] = strings.Replace(l, "_label", "_input", 1)
	}

	for i := range cores {
		cores[i].t = "P-Core"
		if cores[i].f == min && max != min { cores[i].t = "E-Core" }
		cores[i].p = lblMap[fmt.Sprintf("Core %d", cores[i].id)]
	}

	// 3. UI rendering loop
	render := func() {
		if !*once { fmt.Print("\033[H\033[2J") }
		fmt.Printf("💻 HYBRID CPU MONITOR | ID | Type   | Max Clock | Temp\n%s\n", strings.Repeat("=", 55))
		for _, c := range cores {
			tStr := "N/A"
			if c.p != "" {
				b, err := os.ReadFile(c.p)
				v, _ := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
				if err == nil { tStr = fmt.Sprintf("%.1f°C", v/1000) }
			}
			fmt.Printf("CPU %d\t| %s\t| %.2f GHz\t| %s\n", c.id, c.t, float64(c.f)/1000000, tStr)
		}
	}

	render()
	if *once { return }

	for range time.Tick(10 * time.Second) { render() }
}