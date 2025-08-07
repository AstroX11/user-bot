package utils

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var startTime = time.Now()

func PortServe() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		uptime := time.Since(startTime)

		fmt.Fprintf(w, "PID: %d\n", os.Getpid())
		fmt.Fprintf(w, "Uptime: %s\n", uptime)
		fmt.Fprintf(w, "Go Version: %s\n", runtime.Version())
		fmt.Fprintf(w, "Num CPU: %d\n", runtime.NumCPU())
		fmt.Fprintf(w, "Num Goroutines: %d\n", runtime.NumGoroutine())
		fmt.Fprintf(w, "Compiler: %s\n", runtime.Compiler)
		fmt.Fprintf(w, "OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		fmt.Fprintf(w, "Alloc: %d KB\n", memStats.Alloc/1024)
		fmt.Fprintf(w, "TotalAlloc: %d KB\n", memStats.TotalAlloc/1024)
		fmt.Fprintf(w, "Sys: %d KB\n", memStats.Sys/1024)
		fmt.Fprintf(w, "NumGC: %d\n", memStats.NumGC)

		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Fprintf(w, "Error fetching interfaces: %v\n", err)
		} else {
			for _, iface := range ifaces {
				fmt.Fprintf(w, "- Name: %s\n", iface.Name)
				fmt.Fprintf(w, "  HardwareAddr: %s\n", iface.HardwareAddr.String())
				fmt.Fprintf(w, "  Flags: %s\n", iface.Flags.String())
				addrs, _ := iface.Addrs()
				for _, addr := range addrs {
					fmt.Fprintf(w, "  Addr: %s\n", addr.String())
				}
			}
		}

		env := os.Environ()
		for i, e := range env {
			if i >= 10 {
				fmt.Fprintf(w, "...and %d more\n", len(env)-10)
				break
			}
			pair := strings.SplitN(e, "=", 2)
			fmt.Fprintf(w, "%s = %s\n", pair[0], pair[1])
		}
	})

	now := time.Now().Format("15:04:05.000")
	fmt.Printf("\033[34m%s [Server INFO] SERVER HOST: 8000\033[0m\n", now)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
