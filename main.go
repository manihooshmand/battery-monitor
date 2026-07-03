package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Defining Prometheus metrics
var (
	batteryStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "battery_status_info",
		Help: "Battery status (1=Charging, 2=Discharging, 3=Full, 0=Unknown).",
	}, []string{"status"})

	batteryCapacity = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_capacity_percent",
		Help: "Battery capacity percentage reported by system.",
	})

	batteryWear = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_wear_percent",
		Help: "Battery wear level calculated from full/design capacity.",
	})

	batteryPowerWatts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_power_watts",
		Help: "Current power flow in watts.",
	})

	batteryEnergyNow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_energy_now_micro_wh",
		Help: "Current energy in micro Watt-hours.",
	})

	batteryRemainingTimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_remaining_time_seconds",
		Help: "Estimated remaining time in seconds.",
	})
)

// Helper functions
func readFirstExisting(base string, names []string) (string, error) {
	for _, n := range names {
		p := filepath.Join(base, n)
		if _, err := os.Stat(p); err == nil {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(string(b)), nil
		}
	}
	return "", fmt.Errorf("none exist")
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return v
	}
	fields := strings.Fields(s)
	if len(fields) > 0 {
		v, _ = strconv.ParseFloat(fields[0], 64)
		return v
	}
	return 0
}

func parseInt(s string) int64 {
	s = strings.TrimSpace(s)
	v, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return v
	}
	fields := strings.Fields(s)
	if len(fields) > 0 {
		v, _ = strconv.ParseInt(fields[0], 10, 64)
		return v
	}
	return 0
}

// Main logic to read sysfs and update metrics
func collectMetrics(base string) {
	status, _ := readFirstExisting(base, []string{"status"})
	capStr, _ := readFirstExisting(base, []string{"capacity"})
	efStr, _ := readFirstExisting(base, []string{"energy_full", "charge_full"})
	edStr, _ := readFirstExisting(base, []string{"energy_full_design", "charge_full_design"})
	enStr, _ := readFirstExisting(base, []string{"energy_now", "charge_now"})
	pwStr, _ := readFirstExisting(base, []string{"power_now"})

	ef := parseFloat(efStr)
	ed := parseFloat(edStr)
	en := parseFloat(enStr)
	pw := parseFloat(pwStr)
	pwW := pw / 1e6

	st := strings.ToLower(status)
	batteryStatus.Reset()
	if st == "charging" {
		batteryStatus.WithLabelValues("Charging").Set(1)
	} else if st == "discharging" {
		batteryStatus.WithLabelValues("Discharging").Set(2)
	} else if st == "full" {
		batteryStatus.WithLabelValues("Full").Set(3)
	} else {
		batteryStatus.WithLabelValues(status).Set(0)
	}

	batteryCapacity.Set(float64(parseInt(capStr)))
	batteryPowerWatts.Set(pwW)
	batteryEnergyNow.Set(en)

	if ef > 0 && ed > 0 {
		w := (1 - (ef / ed)) * 100
		if w < 0 {
			w = 0
		}
		batteryWear.Set(w)
	}

	if pwW > 0 {
		var sec float64
		if st == "discharging" {
			sec = (en / 1e6) / pwW * 3600
		} else if st == "charging" {
			sec = ((ef - en) / 1e6) / pwW * 3600
		}
		batteryRemainingTimeSeconds.Set(sec)
	} else {
		batteryRemainingTimeSeconds.Set(0)
	}
}

func main() {
	// READ PATH FROM ENVIRONMENT VARIABLE, DEFAULT TO /sys/class/power_supply/BAT1
	base := os.Getenv("BATTERY_PATH")
	if base == "" {
		base = "/sys/class/power_supply/BAT1"
	}

	prometheus.MustRegister(batteryStatus, batteryCapacity, batteryWear, batteryPowerWatts, batteryEnergyNow, batteryRemainingTimeSeconds)

	go func() {
		for {
			collectMetrics(base)
			time.Sleep(10 * time.Second)
		}
	}()

	fmt.Println("Starting Prometheus Battery Exporter on :9191/metrics")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9191", nil)
}
