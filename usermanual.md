# MatrixPulse User Manual

**Version 1.0**

---

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [First Run](#first-run)
4. [Understanding the Dashboard](#understanding-the-dashboard)
5. [Configuration Guide](#configuration-guide)
6. [Interpreting Results](#interpreting-results)
7. [API Integration](#api-integration)
8. [Performance Tuning](#performance-tuning)
9. [Troubleshooting](#troubleshooting)
10. [Advanced Usage](#advanced-usage)

---

## Introduction

MatrixPulse is a real-time financial correlation analysis tool that monitors relationships between multiple assets simultaneously. It uses eigenvalue decomposition to detect market regimes and identifies anomalous correlation spikes that may indicate systemic risk.

### What MatrixPulse Does

- **Computes rolling correlation matrices** between all tracked symbols
- **Detects market regimes** (Normal, Stressed, Crisis) via eigenvalue analysis
- **Generates alerts** when correlations or eigenvalues exceed thresholds
- **Provides real-time visualization** through a native desktop application
- **Exposes APIs** for integration with external systems

### Who Should Use This

- Quantitative analysts monitoring portfolio risk
- Traders tracking inter-market relationships
- Risk managers assessing systemic exposure
- Researchers studying financial network dynamics

---

## Installation

### System Requirements

- **Operating System**: Windows 10+, macOS 10.14+, or Linux (Ubuntu 18.04+)
- **CPU**: Dual-core processor (quad-core recommended)
- **RAM**: 2GB minimum (4GB recommended)
- **Display**: 1024×768 minimum resolution

### Prerequisites

**All Platforms:**
- Go 1.21 or later ([download](https://go.dev/dl/))

**Linux:**
```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev gcc
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
- Install MinGW-w64 from [mingw-w64.org](https://www.mingw-w64.org/)
- Add MinGW bin directory to PATH

### Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/matrixpulse.git
cd matrixpulse

# Download dependencies
make deps

# Build application
make build

# Binary will be in bin/matrixpulse
```

### Pre-built Binaries

Download the latest release for your platform:
- [Windows .exe](https://github.com/yourusername/matrixpulse/releases)
- [macOS .app](https://github.com/yourusername/matrixpulse/releases)
- [Linux binary](https://github.com/yourusername/matrixpulse/releases)

---

## First Run

### Quick Start

1. Navigate to the MatrixPulse directory
2. Run `./bin/matrixpulse` (or `matrixpulse.exe` on Windows)
3. A GUI window will appear showing the dashboard
4. Data will start flowing immediately (simulated feed)

### Initial Configuration

On first run, MatrixPulse uses default settings:
- **Symbols**: AAPL, GOOGL, MSFT, AMZN, TSLA, META
- **Window Size**: 120 data points
- **Update Frequency**: 40 Hz
- **Data Source**: Simulated market feed

To customize, create a `config.yaml` file in the application directory (see Configuration Guide).

---

## Understanding the Dashboard

### Window Layout

The dashboard is divided into four main sections:

#### 1. System Status Bar (Top)
- **Uptime**: How long the system has been running
- **Update Rate**: Actual GUI refresh rate (FPS)
- **Symbol Count**: Number of tracked assets
- **Alert Count**: Total alerts generated
- **Current Regime**: Market state summary

#### 2. Market Regime Panel
- **Regime Indicator**: 
  - 🟢 **NORMAL**: Market operating normally
  - 🟡 **STRESSED**: Elevated condition number (>50)
  - 🔴 **CRISIS**: Max eigenvalue exceeds threshold
- **Max Eigenvalue**: Largest eigenvalue of correlation matrix
- **Condition Number**: Matrix stability measure (λ_max / λ_min)
- **Dominant Eigenvalues**: Top 6 eigenvalues in descending order

#### 3. Correlation Matrix Display
- Shows real-time correlation coefficients between all tracked symbols
- Values range from -1 (perfect negative correlation) to +1 (perfect positive)
- Diagonal is always 1.000 (asset correlated with itself)
- Display limited to 10×10 for readability (configurable)

#### 4. Active Alerts Panel
- Lists recent correlation spikes and regime changes
- Each alert shows:
  - **Timestamp**: When alert was generated
  - **Severity**: CRITICAL, HIGH, or WARNING
  - **Symbol(s)**: Assets involved
  - **Metric**: What threshold was exceeded
  - **Value**: Actual measured value

### Color Coding

- 🔴 **Red (Critical)**: Immediate attention required
- 🟠 **Orange (High)**: Significant event detected
- ⚠️  **Yellow (Warning)**: Noteworthy but not urgent
- 🟢 **Green (Normal)**: System operating normally

---

## Configuration Guide

### Configuration File Location

Create `config.yaml` in the same directory as the MatrixPulse executable.

### Complete Configuration Example

```yaml
# List of symbols to track
symbols:
  - AAPL
  - GOOGL
  - MSFT
  - AMZN
  - TSLA
  - META
  - NVDA
  - AMD

# Rolling window size (number of price points)
# Larger = smoother but slower to react
# Smaller = more responsive but noisier
window_size: 120

# Computation frequency (calculations per second)
# Higher = more real-time but more CPU intensive
update_hz: 40

# Alert thresholds
alerts:
  # Trigger when |correlation| exceeds this value
  correlation_threshold: 0.82
  
  # Trigger crisis mode when max eigenvalue exceeds this
  eigenvalue_threshold: 2.8
  
  # (Future use) Volatility spike detection
  volatility_threshold: 0.04

# State persistence
persistence:
  enabled: true
  path: "matrixpulse_state.json"
  interval_seconds: 60  # Save every 60 seconds

# WebSocket server for live streaming
websocket:
  enabled: true
  port: 8080

# REST API server
rest:
  enabled: true
  port: 8081

# GUI dashboard
dashboard:
  enabled: true
  refresh_ms: 200  # Update every 200ms
```

### Configuration Parameters Explained

#### Symbols
- **Purpose**: Define which assets to track
- **Min**: 1 symbol
- **Max**: 100 symbols (performance decreases with more)
- **Format**: Standard ticker symbols

#### Window Size
- **Purpose**: Number of data points in rolling calculation
- **Min**: 10 (not recommended, too unstable)
- **Recommended**: 60-240
- **Max**: 10,000
- **Effect**: Larger windows smooth out noise but reduce responsiveness

#### Update Hz
- **Purpose**: How often calculations are performed per second
- **Min**: 1 Hz
- **Recommended**: 20-50 Hz
- **Max**: 1000 Hz (extreme, not useful)
- **Effect**: Higher frequencies increase CPU usage

#### Correlation Threshold
- **Purpose**: Define what constitutes an "abnormal" correlation
- **Range**: 0.0 to 1.0
- **Recommended**: 0.75-0.90
- **Interpretation**: |r| > threshold triggers alert

#### Eigenvalue Threshold
- **Purpose**: Crisis mode detection
- **Typical Range**: 2.0-4.0
- **Recommended**: 2.5-3.0
- **Interpretation**: Max eigenvalue > threshold = crisis mode

---

## Interpreting Results

### Correlation Matrix

**What it shows:** Linear relationship strength between asset pairs.

**Reading values:**
- **+1.00**: Perfect positive correlation (move together)
- **+0.80**: Strong positive correlation
- **+0.50**: Moderate positive correlation
- **+0.00**: No linear relationship
- **-0.50**: Moderate negative correlation
- **-0.80**: Strong negative correlation
- **-1.00**: Perfect negative correlation (mirror opposites)

**Practical interpretation:**
- High positive correlations (>0.8) between many pairs indicates **market-wide risk**
- Sudden correlation changes (0.2 → 0.9) signals **regime shift**
- Negative correlations between related assets may indicate **arbitrage opportunities**

### Eigenvalues

**What they show:** Dimensionality and risk concentration in the market.

**Max eigenvalue interpretation:**
- **< 2.0**: Well-diversified, independent assets
- **2.0-3.0**: Moderate systemic correlation (normal)
- **> 3.0**: Excessive correlation, systemic risk elevated
- **> 5.0**: Crisis conditions, assets moving as one

**Condition number interpretation:**
- **< 10**: Stable correlation matrix
- **10-50**: Acceptable conditioning
- **> 50**: Stressed state, eigenvalue concentration
- **> 100**: Numerical instability likely

### Market Regimes

**NORMAL (🟢)**
- Max eigenvalue < threshold
- Condition number < 50
- Assets exhibiting independent behavior
- **Action**: Normal operations

**STRESSED (🟡)**
- Condition number > 50 but eigenvalue < threshold
- Elevated correlation but not critical
- **Action**: Monitor closely, reduce leverage

**CRISIS (🔴)**
- Max eigenvalue > threshold
- Extreme systemic correlation
- Assets losing independence
- **Action**: Risk-off, defensive positioning

---

## API Integration

### REST API Endpoints

**Base URL**: `http://localhost:8081`

#### GET /matrix
Returns current correlation matrix.

**Response:**
```json
{
  "Cor": [[1.0, 0.85, ...], [0.85, 1.0, ...], ...],
  "Symbols": ["AAPL", "GOOGL", "MSFT", ...],
  "Time": "2025-01-15T10:30:00Z"
}
```

#### GET /mode
Returns eigenvalue analysis and regime.

**Response:**
```json
{
  "Eigenvalues": [2.45, 1.32, 1.18, 0.95, ...],
  "MaxEigen": 2.45,
  "Condition": 8.5,
  "Regime": "NORMAL",
  "Time": "2025-01-15T10:30:00Z"
}
```

#### GET /alerts
Returns all generated alerts.

**Response:**
```json
[
  {
    "Level": "HIGH",
    "Symbol": "AAPL-GOOGL",
    "Message": "correlation spike",
    "Value": 0.89,
    "Threshold": 0.82,
    "Time": "2025-01-15T10:30:00Z"
  }
]
```

#### GET /health
Health check endpoint.

**Response:** `ok`

### WebSocket Stream

**URL**: `ws://localhost:8080/ws`

**Protocol**: JSON messages sent every second

**Message format:**
```json
{
  "matrix": { /* correlation matrix */ },
  "mode": { /* eigenvalue analysis */ },
  "alerts": [ /* alert array */ ]
}
```

**Client example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Regime:', data.mode.Regime);
  console.log('Max Eigen:', data.mode.MaxEigen);
};
```

---

## Performance Tuning

### For Low-End Systems

```yaml
window_size: 60
update_hz: 20
dashboard:
  refresh_ms: 500
```

**Effect**: Reduces CPU usage by ~60%, still responsive.

### For High-Performance Systems

```yaml
window_size: 240
update_hz: 60
dashboard:
  refresh_ms: 100
```

**Effect**: Maximum responsiveness and smoothness.

### Memory Optimization

Memory usage ≈ `window_size × symbols × 8 bytes`

Example: 120 window × 10 symbols × 8 bytes = 9.6 KB per rolling window

Total memory scales linearly with symbol count.

---

## Troubleshooting

### GUI Won't Start

**Symptom**: Application crashes immediately or shows blank window

**Solution (Linux)**:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**Solution (macOS)**:
```bash
xcode-select --install
```

### High CPU Usage

**Symptom**: CPU usage > 80% on one or more cores

**Solution**: Reduce `update_hz` in config:
```yaml
update_hz: 20  # Down from 40
```

### Alerts Flooding

**Symptom**: Hundreds of alerts, unreadable dashboard

**Solution**: Increase thresholds:
```yaml
alerts:
  correlation_threshold: 0.90  # Up from 0.82
  eigenvalue_threshold: 3.5    # Up from 2.8
```

### Matrix Shows "Waiting for data"

**Symptom**: Matrix never populates

**Cause**: Not enough data points collected yet

**Solution**: Wait for `window_size` data points to accumulate (~3 seconds at default feed rate)

### Port Already in Use

**Symptom**: Error message about port binding

**Solution**: Change ports in config:
```yaml
websocket:
  port: 9080  # Different port
rest:
  port: 9081  # Different port
```

---

## Advanced Usage

### Headless Mode

Run without GUI for server deployment:

```yaml
dashboard:
  enabled: false
```

Application will run until receiving SIGINT (Ctrl+C).

### Custom Data Sources

Replace `feed.NewSimulated()` in main.go with your own feed implementation:

```go
type CustomFeed struct {
    // Your implementation
}

func (c *CustomFeed) Start(ctx context.Context) <-chan Tick {
    // Return channel of market ticks
}
```

### Correlation Heatmap Export

Extract matrix via REST API and visualize:

```python
import requests
import seaborn as sns
import matplotlib.pyplot as plt

response = requests.get('http://localhost:8081/matrix')
data = response.json()

sns.heatmap(data['Cor'], 
            xticklabels=data['Symbols'],
            yticklabels=data['Symbols'],
            cmap='RdYlGn', center=0)
plt.show()
```

### Alert Webhook Integration

Modify engine.go to send alerts to external systems:

```go
func (e *Engine) addAlert(a types.Alert) {
    // Existing code...
    
    if a.Level == "CRITICAL" {
        sendWebhook("https://your-webhook-url", a)
    }
}
```

---

## Appendix A: Statistical Background

### Log Returns

MatrixPulse uses logarithmic returns:

```
r_t = ln(P_t / P_{t-1})
```

**Advantages:**
- Time-additive
- Symmetric treatment of gains/losses
- Approximates percentage changes for small moves

### Pearson Correlation

Correlation coefficient between X and Y:

```
ρ = Cov(X,Y) / (σ_X × σ_Y)
```

Where:
- Cov(X,Y) = covariance between X and Y
- σ_X, σ_Y = standard deviations

### Eigenvalue Decomposition

Correlation matrix C can be decomposed:

```
C = Q Λ Q^T
```

Where:
- Q = eigenvector matrix
- Λ = diagonal matrix of eigenvalues

**Interpretation:**
- Each eigenvalue represents variance explained by corresponding principal component
- Sum of eigenvalues = number of assets (for correlation matrices)
- Large first eigenvalue = market-wide common factor

---

## Appendix B: Mathematical Definitions

**Condition Number**:
```
κ(C) = λ_max / λ_min
```

High condition number indicates near-singularity (eigenvalue concentration).

**Rolling Window**:

Computation over last N observations:
```
Cor_t = Corr(R_{t-N+1:t})
```

Where R is the matrix of returns.

---

**Document Version**: 1.0  
**Last Updated**: January 2025  
**For Support**: https://github.com/yourusername/matrixpulse/issues