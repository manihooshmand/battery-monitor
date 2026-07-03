##🔋 Battery Monitor Stack

A complete, containerized monitoring stack for Linux laptops using Go, Prometheus, and Grafana.

This project reads battery metrics directly from /sys/class/power_supply/ and exposes them as Prometheus metrics, which are then visualized in a beautiful Grafana dashboard.

✨ Features
Real-time Metrics: Battery percentage, power flow (Watts), wear level, and estimated remaining time.
Fully Containerized: Runs everything (Exporter, Prometheus, Grafana) via Docker Compose. No system-level installations required.
Out-of-the-box Dashboard: Comes with a pre-configured Grafana dashboard.
📸 Dashboard Preview
Battery Monitor Grafana Dashboard

🚀 Quick Start
Clone the repository:
git clone -b observability-integration https://github.com/manihooshmand/battery-monitor.gitcd battery-monitor
Start the stack with Docker Compose:
bash

sudo docker compose up -d --build
Access the services:
Grafana: http://localhost:3000 (Admin / Admin)
Prometheus: http://localhost:9090
🛠️ Tech Stack
Exporter: Go
Metrics Storage: Prometheus
Visualization: Grafana
Orchestration: Docker & Docker Compose
