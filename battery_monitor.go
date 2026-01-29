package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

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

func formatDuration(seconds float64) string {
    if seconds <= 0 {
        return "unknown"
    }
    d := time.Duration(seconds) * time.Second
    h := int(d.Hours())
    m := int(d.Minutes()) % 60
    return fmt.Sprintf("%02dh %02dm", h, m)
}

func main() {
    base := "/sys/class/power_supply/BAT1"
    if len(os.Args) > 1 {
        base = os.Args[1]
    }

    interval := 10 * time.Second
    fmt.Printf("Reading battery info from %s every %s\n", base, interval)

    for {
        info := make(map[string]string)

        status, _ := readFirstExisting(base, []string{"status"})
        info["status"] = status

        capStr, _ := readFirstExisting(base, []string{"capacity"})
        info["capacity_pct"] = fmt.Sprintf("%d", parseInt(capStr))

        efStr, _ := readFirstExisting(base, []string{"energy_full", "charge_full"})
        edStr, _ := readFirstExisting(base, []string{"energy_full_design", "charge_full_design"})
        enStr, _ := readFirstExisting(base, []string{"energy_now", "charge_now"})
        pwStr, _ := readFirstExisting(base, []string{"power_now"})

        ef := parseFloat(efStr)
        ed := parseFloat(edStr)
        en := parseFloat(enStr)
        pw := parseFloat(pwStr)

        info["energy_full_raw"] = fmt.Sprintf("%.0f", ef)
        info["energy_design_raw"] = fmt.Sprintf("%.0f", ed)
        info["energy_now_raw"] = fmt.Sprintf("%.0f", en)
        info["power_raw"] = fmt.Sprintf("%.0f", pw)

        if ef > 0 && ed > 0 {
            w := (1 - (ef / ed)) * 100
            if w < 0 {
                w = 0
            }
            info["wear_pct"] = fmt.Sprintf("%.2f", w)
        }

        pwW := pw / 1e6
        info["power_w"] = fmt.Sprintf("%.4f", pwW)

        if ef > 0 {
            rp := (en / ef) * 100
            if rp > 100 {
                rp = 100
            }
            info["remaining_of_full_pct"] = fmt.Sprintf("%.2f", rp)
        }

        var rt string
        if pwW > 0 {
            st := strings.ToLower(status)
            if st == "discharging" {
                sec := (en / 1e6) / pwW * 3600
                rt = formatDuration(sec)
            } else if st == "charging" {
                sec := ((ef - en) / 1e6) / pwW * 3600
                rt = formatDuration(sec)
            } else {
                rt = "n/a"
            }
        } else {
            rt = "n/a"
        }
        info["remaining_time"] = rt

        now := time.Now().Format(time.RFC3339)
        fmt.Printf("=== %s ===\n", now)

        ordered := []string{
            "status", "capacity_pct",
            "power_w", "wear_pct",
            "energy_full_raw", "energy_design_raw", "energy_now_raw", "power_raw",
			"remaining_time", "remaining_of_full_pct",
        }

        for _, k := range ordered {
            if v, ok := info[k]; ok {
                fmt.Printf("%-24s : %s\n", k, v)
            }
        }

        fmt.Println()
        time.Sleep(interval)
    }
}