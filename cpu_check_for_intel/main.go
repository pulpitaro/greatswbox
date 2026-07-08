package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type CPUCore struct {
	ID       int
	MaxFreq  int64
	Temp     float64
	CoreType string
}

// Clears the terminal screen using the system 'clear' command
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Reads core temperature from the coretemp hwmon driver
func getCoreTemp(coreID int) float64 {
	matches, _ := filepath.Glob("/sys/class/hwmon/hwmon*/name")
	for _, match := range matches {
		nameBytes, _ := ioutil.ReadFile(match)
		if strings.TrimSpace(string(nameBytes)) == "coretemp" {
			dir := filepath.Dir(match)
			labels, _ := filepath.Glob(dir + "/input*_label")
			for _, labelPath := range labels {
				labelBytes, _ := ioutil.ReadFile(labelPath)
				labelText := strings.TrimSpace(string(labelBytes))
				
				// Core labels usually match "Core 0", "Core 1", etc.
				if labelText == fmt.Sprintf("Core %d", coreID) {
					// Swap '_label' with '_input' to read the actual temperature value
					inputPath := strings.Replace(labelPath, "_label", "_input", 1)
					tempBytes, _ := ioutil.ReadFile(inputPath)
					tempStr := strings.TrimSpace(string(tempBytes))
					tempMilli, _ := strconv.ParseFloat(tempStr, 64)
					return tempMilli / 1000.0 // Convert milli-Celsius to Celsius
				}
			}
		}
	}
	return 0.0
}

// Scans sysfs to detect CPU cores, frequencies, and prints the mapped architecture
func printCPUMap() {
	cpuDirs, _ := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*")
	var cores []CPUCore
	var highestFreq int64 = 0
	var lowestFreq int64 = 9999999999

	for _, dir := range cpuDirs {
		base := filepath.Base(dir)
		id, _ := strconv.Atoi(base[3:]) // Extract ID from "cpuX"

		freqBytes, err := ioutil.ReadFile(dir + "/cpufreq/cpuinfo_max_freq")
		if err != nil {
			continue // Skip hyper-threading logical pairs if they share sysfs attributes
		}
		freqStr := strings.TrimSpace(string(freqBytes))
		freq, _ := strconv.ParseInt(freqStr, 10, 64)

		if freq > highestFreq {
			highestFreq = freq
		}
		if freq < lowestFreq {
			lowestFreq = freq
		}

		cores = append(cores, CPUCore{
			ID:      id,
			MaxFreq: freq,
		})
	}

	fmt.Printf("💻 INTEL HYBRID ARCHITECTURE DETECTOR & MONITOR 💻\n")
	fmt.Printf("==================================================\n")
	fmt.Printf("ID\t| Type\t\t| Max Clock\t| Temperature\n")
	fmt.Printf("--------------------------------------------------\n")

	for i := range cores {
		// Intel P-Cores always expose higher maximum frequencies than E-Cores
		if cores[i].MaxFreq == lowestFreq && highestFreq != lowestFreq {
			cores[i].CoreType = "E-Core (Efficient)"
		} else {
			cores[i].CoreType = "P-Core (Performance)"
		}

		cores[i].Temp = getCoreTemp(cores[i].ID)

		tempStr := fmt.Sprintf("%.1f°C", cores[i].Temp)
		if cores[i].Temp == 0.0 {
			tempStr = "N/A (HT Link)" // HT threads share thermal sensor with the physical core
		}

		fmt.Printf("CPU %d\t| %s\t| %.2f GHz\t| %s\n", 
			cores[i].ID, 
			cores[i].CoreType, 
			float64(cores[i].MaxFreq)/1000000.0, 
			tempStr,
		)
	}
	fmt.Printf("--------------------------------------------------\n")
}

func main() {
	// Custom English help and usage menu
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This tool detects Intel Hybrid Architecture (P/E Cores) and tracks temperatures.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c\t\tRun once and exit immediately (snapshot mode)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\tShow this help message\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  sudo %s         # Runs as a live monitor (10s refresh)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  sudo %s -c      # Prints the CPU map once and exits\n\n", os.Args[0])
	}

	onceFlag := flag.Bool("c", false, "Run once and exit immediately")
	flag.Parse()

	// If -c flag is passed, execute once and exit
	if *onceFlag {
		printCPUMap()
		return
	}

	// Fallback to continuous monitoring mode with a 10-second ticker
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clearScreen()
	printCPUMap()
	fmt.Printf("[Press Ctrl+C to exit]\n")

	for {
		select {
		case <-ticker.C:
			clearScreen()
			printCPUMap()
			fmt.Printf("[Press Ctrl+C to exit]\n")

		case <-sigChan:
			clearScreen()
			fmt.Println("🛑 Monitor stopped.")
			return
		}
	}
}