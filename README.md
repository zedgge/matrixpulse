# MatrixPulse

**Real-time Financial Correlation Analytics Engine**

MatrixPulse is a concurrent analytics system that computes rolling correlation matrices, eigenvalue decomposition, and market regime detection in real-time. Built for reliability and performance.

---

## Features

- **Real-time Correlation Matrix** - Continuous computation of asset correlations
- **Eigenvalue Analysis** - Market regime detection (Normal/Stressed/Crisis)
- **Alert System** - Configurable thresholds for correlation spikes and eigenvalue anomalies
- **Native GUI Dashboard** - Clean, responsive interface showing live metrics
- **REST API** - JSON endpoints for external integration
- **WebSocket Stream** - Live data feed for clients
- **State Persistence** - Automatic snapshots for recovery

---

## Quick Start

### Prerequisites

- Go 1.21 or later
- GCC (for Fyne GUI compilation)

**Platform-specific requirements:**
- **Linux**: `libgl1-mesa-dev xorg-dev`
- **macOS**: Xcode command line tools
- **Windows**: GCC via MinGW-w64

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/matrixpulse.git
cd matrixpulse

# Install dependencies
go mod download

# Build
make build

# Run
./bin/matrixpulse
```

### Docker (Alternative)

```bash
docker build -t matrixpulse .
docker run -p 8080:8080 -p 8081:8081 matrixpulse
```

---

## Configuration

Edit `config.yaml` to customize behavior:

```yaml
symbols:
  - AAPL
  - GOOGL
  - MSFT
  - AMZN
  - TSLA
  - META

window_size: 120        # Rolling window size (data points)
update_hz: 40          # Computation frequency (per second)

alerts:
  correlation_threshold: 0.82   # Alert when |correlation| exceeds
  eigenvalue_threshold: 2.8     # Crisis mode trigger
  volatility_threshold: 0.04    # Volatility spike alert

persistence:
  enabled: true
  path: "state.json"
  interval_seconds: 60

websocket:
  enabled: true
  port: 8080

rest:
  enabled: true
  port: 8081

dashboard:
  enabled: true
  refresh_ms: 200
```

---

## Architecture

### Concurrency Model

MatrixPulse runs multiple coordinated goroutines:

1. **Feed Loop** - Ingests market data ticks
2. **Compute Loop** - Performs rolling calculations (40 Hz default)
3. **Persistence Loop** - Saves state snapshots (60s default)
4. **Server Loops** - REST and WebSocket interfaces
5. **GUI Loop** - Updates dashboard display (200ms default)

All loops are context-managed for clean shutdown.

### Data Flow

```
Market Feed → Rolling Windows → Statistics → Correlation Matrix
                                              ↓
                                         Eigenvalue Analysis
                                              ↓
                                         Regime Detection
                                              ↓
                                    GUI / REST / WebSocket
```

---

## API Reference

### REST Endpoints

**GET /matrix**
```json
{
  "Cor": [[1.0, 0.85, ...], ...],
  "Symbols": ["AAPL", "GOOGL", ...],
  "Time": "2025-01-15T10:30:00Z"
}
```

**GET /mode**
```json
{
  "Eigenvalues": [2.45, 1.32, ...],
  "MaxEigen": 2.45,
  "Condition": 8.5,
  "Regime": "NORMAL",
  "Time": "2025-01-15T10:30:00Z"
}
```

**GET /alerts**
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

**GET /health**
```
ok
```

### WebSocket

Connect to `ws://localhost:8080/ws` to receive updates every second:

```json
{
  "matrix": { ... },
  "mode": { ... },
  "alerts": [ ... ]
}
```

---

## GUI Dashboard

The native GUI displays:

- **Market Regime** - Current state with color coding (🟢 Normal / 🟡 Stressed / 🔴 Crisis)
- **Eigenvalue Stats** - Max eigenvalue, condition number, top 6 eigenvalues
- **Correlation Matrix** - Live heatmap-style grid (top 10×10 symbols)
- **Active Alerts** - Recent alerts with severity indicators

### Controls

- **Window close** triggers graceful shutdown of all systems
- **Refresh rate** is configurable in `config.yaml`

---

## Performance

- **Latency**: Sub-millisecond tick-to-compute
- **Throughput**: Handles 1000+ ticks/second per symbol
- **Memory**: ~50MB baseline + (window_size × symbols × 8 bytes)
- **CPU**: Scales with `GOMAXPROCS` (defaults to all cores)

### Benchmarks

6 symbols, 120-point window, 40 Hz:
- Compute cycle: ~500μs
- GUI update: ~2ms
- REST response: <1ms

---

## Troubleshooting

### GUI won't start

**Issue**: Missing graphics libraries

**Linux**:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**macOS**:
```bash
xcode-select --install
```

### High CPU usage

**Solution**: Reduce `update_hz` in config.yaml (try 20 Hz)

### Alerts flooding

**Solution**: Increase thresholds in config.yaml:
```yaml
alerts:
  correlation_threshold: 0.90
  eigenvalue_threshold: 3.5
```

### State file corruption

**Solution**: Delete `state.json` and restart (will recompute)

---

## Building from Source

### Development Build

```bash
make dev
```

### Production Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

---

## License

MIT License - See LICENSE file

---

## Technical Details

### Statistical Methods

- **Returns**: Log returns computed as `ln(Pt / Pt-1)`
- **Correlation**: Pearson correlation coefficient
- **Covariance**: Unbiased estimator with Bessel's correction
- **Eigenvalues**: Computed via Gonum's symmetric eigenvalue decomposition

### Market Regimes

- **NORMAL**: MaxEigen < threshold, Condition < 50
- **STRESSED**: Condition > 50
- **CRISIS**: MaxEigen > threshold (default 2.8)

Condition number = λ_max / λ_min measures matrix stability.

---

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/yourusername/matrixpulse/issues
- Documentation: https://matrixpulse.dev/docs

---

## Acknowledgments

Built with:
- [Fyne](https://fyne.io) - Native GUI toolkit
- [Gonum](https://gonum.org) - Numerical computing
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation