# 🔋 Battery Monitor Observability

A complete, containerized monitoring stack for Linux laptops using **Go**, **Prometheus**, and **Grafana**.

This project reads battery metrics directly from `/sys/class/power_supply/` and exposes them as Prometheus metrics, which are then visualized through a pre-configured Grafana dashboard.

---

## ✨ Features

- 📊 **Real-time battery metrics**
  - Battery percentage
  - Charging/discharging status
  - Power consumption (Watts)
  - Battery wear level
  - Estimated remaining time

- 🐳 **Fully containerized**
  - Runs entirely with Docker Compose
  - No system-wide installation required

- 📈 **Built-in Grafana dashboard**
  - Ready-to-use dashboard with pre-configured panels

- ⚡ **Prometheus integration**
  - Exposes metrics in Prometheus format for easy monitoring

---

## 📸 Dashboard Preview

---<img width="1621" height="757" alt="Screenshot from 2026-07-03 06-30-34" src="https://github.com/user-attachments/assets/7fbd7dfa-5444-4363-b0df-e4df72641292" />

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone -b observability-integration https://github.com/manihooshmand/battery-monitor.git

cd battery-monitor
```

### 2. Build and start the stack

```bash
sudo docker compose up -d --build
```

### 3. Open the services

| Service | URL |
|---------|-----|
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

**Grafana Default Credentials**

- **Username:** `admin`
- **Password:** `admin`


## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| Exporter | Go |
| Metrics Collection | Prometheus |
| Visualization | Grafana |
| Containerization | Docker & Docker Compose |
