# 🔋 Battery Monitor (Go)

A lightweight command-line utility for monitoring battery status and health on Linux systems.

## Features
- Real-time battery status monitoring
- Charge percentage and health (wear level)
- Power consumption/charge rate (Watts)
- Estimated time remaining
- Raw sysfs data access
- Auto-refresh every 10 seconds

## Installation
### Prerequisites
- Linux system with sysfs battery interface
- Go 1.16+

### Build from Source
```bash
go build -o battery-monitor battery_monitor.go
