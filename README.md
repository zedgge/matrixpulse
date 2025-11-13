# MatrixPulse

**Real-time Financial Correlation Engine** — Desktop application for monitoring market correlations and detecting systemic risk through eigenvalue decomposition.

---

## What It Does

MatrixPulse continuously computes rolling correlation matrices across multiple assets, performs eigenvalue analysis to detect market regimes (Normal/Stressed/Crisis), and alerts you to correlation spikes in real-time.

**Built for education and exploration of:**
- High-throughput data processing in Go
- Concurrent analytics with goroutines
- Real-time statistical computation
- Native desktop GUI with Fyne

---

## Quick Start

### Prerequisites

- **Go 1.21+**
- **GCC** (for Fyne compilation)

**Platform-specific:**
- **Linux**: `sudo apt-get install libgl1-mesa-dev xorg-dev`
- **macOS**: `xcode-select --install`
- **Windows**: Install MinGW-w64

### Build & Run

```bash
# Clone
git clone https://github.com/yourusername/matrixpulse.git
cd matrixpulse

# Build
make build

# Run
./matrixpulse
```

A GUI window will open showing real-time correlation matrices, eigenvalue analysis, and market regime detection.

---

## Features

### Real-time Correlation Matrix
- Computes Pearson correlations across all tracked symbols
- Rolling window with configurable size (default: 120 points)
- Sub-millisecond compute latency

### Eigenvalue-Based Regime Detection
- **🟢 NORMAL**: Diversified market, eigenvalues distributed
- **🟡 STRESSED**: High condition number (>50), elevated correlation
- **🔴 CRISIS**: Max eigenvalue exceeds threshold, systemic correlation

### Alert System
- Correlation spike detection (threshold: 0.82)
- Crisis mode alerts (eigenvalue > 2.8)
- Real-time notification panel

### Native Desktop GUI
- Live correlation matrix display (10×10 view)
- Eigenvalue statistics
- Alert feed with severity indicators
- System status (uptime, FPS, symbol count)

### State Persistence
- Automatic snapshots every 60 seconds
- JSON format for easy inspection
- Recovery on restart

---

## Configuration

Edit `config.yaml`:

```yaml
symbols:
  - AAPL
  - GOOGL
  - MSFT
  # Add your symbols

window_size: 120      # Rolling window size
update_hz: 40         # Calculations per second

alerts:
  correlation_threshold: 0.82
  eigenvalue_threshold: 2.8

persistence:
  enabled: true
  path: "matrixpulse_state.json"
  interval_seconds: 60

dashboard:
  refresh_ms: 200
  enabled: true
```

---

## Architecture

### Concurrency Model

MatrixPulse runs coordinated goroutines:

1. **Feed Loop** - Generates simulated market ticks (25ms interval)
2. **Ingestion Loop** - Consumes ticks, updates rolling windows
3. **Compute Loop** - Performs matrix calculations (40 Hz)
4. **Persistence Loop** - Saves snapshots (60s interval)
5. **GUI Loop** - Refreshes display (200ms interval)

All loops use `context.Context` for graceful shutdown coordination.

### Data Flow

```
Simulated Feed → Rolling Windows → Log Returns → Covariance Matrix
                                                       ↓
                                                Correlation Matrix
                                                       ↓
                                              Eigen Decomposition
                                                       ↓
                                               Regime Detection
                                                       ↓
                                             GUI Display + Alerts
```

### Statistical Methods

- **Returns**: Log returns `ln(Pt / Pt-1)`
- **Correlation**: Pearson coefficient with Bessel correction
- **Eigenvalues**: Gonum symmetric decomposition
- **Condition Number**: λ_max / λ_min (matrix stability)

---

## Performance

- **Latency**: <1ms tick-to-compute
- **Throughput**: 1000+ ticks/sec/symbol
- **Memory**: ~50MB + (window_size × symbols × 8 bytes)
- **CPU**: Multi-core with `GOMAXPROCS`

**Benchmarks** (6 symbols, 120-point window, 40 Hz):
- Compute cycle: ~500μs
- GUI update: ~2ms

---

## Educational Goals

This project demonstrates:

✅ **Concurrent Processing**
- Goroutine coordination with channels
- Context-based cancellation
- Race-free data sharing (RWMutex)

✅ **Real-time Analytics**
- Rolling window computation
- Matrix operations with Gonum
- Eigenvalue decomposition

✅ **Native GUI Development**
- Cross-platform desktop with Fyne
- Live data visualization
- Event-driven updates

✅ **Production Patterns**
- Configuration management (YAML)
- Graceful shutdown
- State persistence
- Error handling

---

## Troubleshooting

### GUI Won't Start

**Linux:**
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**macOS:**
```bash
xcode-select --install
```

### High CPU Usage

Reduce computation frequency in `config.yaml`:
```yaml
update_hz: 20  # Down from 40
```

### Alert Flooding

Increase thresholds:
```yaml
alerts:
  correlation_threshold: 0.90
  eigenvalue_threshold: 3.5
```

---

## Building from Source

```bash
# Development build with race detector
make dev

# Production build
make build

# Run tests
make test

# Clean artifacts
make clean
```

---

## Project Structure

```
matrixpulse/
├── cmd/matrixpulse/main.go    # Entry point
├── internal/
│   ├── config/                # YAML configuration
│   ├── display/               # Fyne GUI
│   ├── engine/                # Correlation engine
│   ├── feed/                  # Simulated data feed
│   ├── math/                  # Statistical functions
│   ├── persist/               # State snapshots
│   ├── types/                 # Core data types
│   └── window/                # Rolling window
├── config.yaml                # Configuration
├── Makefile                   # Build system
└── README.md
```

---

## Technical Notes

### Why Go?

- Native concurrency (goroutines + channels)
- Excellent performance for numerical workloads
- Strong standard library
- Cross-platform compilation
- Fast build times

### Why Fyne for GUI?

- True native desktop (no web browser)
- Cross-platform (Windows, macOS, Linux)
- Clean API
- Active development

### Data Source

Currently uses a **simulated feed** with synthetic price movements. Designed for easy replacement with live exchange APIs (Alpaca, Polygon, Binance, etc.).

---

## License

MIT License — Free for personal, academic, and commercial use.

---

## Author's Note

MatrixPulse was built to explore high-performance concurrent analytics in Go. It prioritizes **clarity, performance, and reliability** over unnecessary complexity.

The focus is on demonstrating real-time data processing patterns that can scale to production workloads while remaining understandable and maintainable.

**Key Learning Areas:**
- Goroutine-based concurrent architecture
- Real-time statistical computation
- Matrix operations and eigenvalue analysis
- Native desktop GUI development
- Production-ready error handling and shutdown

---

**Questions?** Open an issue on GitHub.