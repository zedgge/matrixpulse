package display

import (
	"context"
	"fmt"
	"strings"
	"time"

	"matrixpulse/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	eng         *engine.Engine
	cancel      context.CancelFunc
	app         fyne.App
	window      fyne.Window
	regimeLabel *widget.Label
	eigenLabel  *widget.Label
	matrixText  *widget.Label
	alertText   *widget.Label
	statsLabel  *widget.Label
}

func NewGUI(eng *engine.Engine, cancel context.CancelFunc) *GUI {
	return &GUI{
		eng:         eng,
		cancel:      cancel,
		regimeLabel: widget.NewLabel("Initializing..."),
		eigenLabel:  widget.NewLabel("Waiting for data..."),
		matrixText:  widget.NewLabel("Loading..."),
		alertText:   widget.NewLabel("No alerts"),
		statsLabel:  widget.NewLabel("System starting..."),
	}
}

func (g *GUI) Run() {
	g.app = app.NewWithID("com.matrixpulse.app")
	g.window = g.app.NewWindow("MatrixPulse - Real-time Financial Correlation Engine")

	g.setupStyles()
	content := g.buildLayout()

	g.window.SetContent(content)
	g.window.Resize(fyne.NewSize(1200, 800))
	g.window.CenterOnScreen()
	
	g.window.SetOnClosed(func() {
		g.cancel()
	})

	go g.updateLoop()
	g.window.ShowAndRun()
}

func (g *GUI) setupStyles() {
	g.regimeLabel.TextStyle = fyne.TextStyle{Bold: true}
	g.matrixText.TextStyle = fyne.TextStyle{Monospace: true}
	g.alertText.TextStyle = fyne.TextStyle{Monospace: true}
	g.statsLabel.TextStyle = fyne.TextStyle{Monospace: true}
}

func (g *GUI) buildLayout() fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle(
		"MATRIXPULSE",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitle := widget.NewLabelWithStyle(
		"Real-time Correlation Matrix & Market Regime Detection",
		fyne.TextAlignCenter,
		fyne.TextStyle{Italic: true},
	)

	headerBox := container.NewVBox(
		header,
		subtitle,
		widget.NewSeparator(),
	)

	// Status section
	statusBox := container.NewVBox(
		widget.NewLabelWithStyle("System Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		g.statsLabel,
		widget.NewSeparator(),
	)

	// Regime section
	regimeBox := container.NewVBox(
		widget.NewLabelWithStyle("Market Regime", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		g.regimeLabel,
		g.eigenLabel,
	)

	// Matrix section
	matrixScroll := container.NewScroll(g.matrixText)
	matrixScroll.SetMinSize(fyne.NewSize(600, 350))

	matrixBox := container.NewVBox(
		widget.NewLabelWithStyle("Correlation Matrix", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		matrixScroll,
	)

	// Alerts section
	alertScroll := container.NewScroll(g.alertText)
	alertScroll.SetMinSize(fyne.NewSize(600, 250))

	alertBox := container.NewVBox(
		widget.NewLabelWithStyle("Active Alerts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		alertScroll,
	)

	// Combine sections
	topSection := container.NewVBox(
		statusBox,
		regimeBox,
		widget.NewSeparator(),
		matrixBox,
	)

	mainContent := container.NewVSplit(topSection, alertBox)
	mainContent.SetOffset(0.65)

	return container.NewBorder(
		headerBox,
		nil,
		nil,
		nil,
		mainContent,
	)
}

func (g *GUI) updateLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	startTime := time.Now()
	updateCount := 0

	for range ticker.C {
		g.updateDisplay()
		updateCount++

		// Update stats every 5 seconds
		if updateCount%25 == 0 {
			elapsed := time.Since(startTime)
			fps := float64(updateCount) / elapsed.Seconds()
			g.updateStats(fps, elapsed)
		}
	}
}

func (g *GUI) updateDisplay() {
	g.updateRegime()
	g.updateMatrix()
	g.updateAlerts()
}

func (g *GUI) updateRegime() {
	mode := g.eng.Mode()
	if mode == nil {
		g.regimeLabel.SetText("⚪ INITIALIZING")
		return
	}

	var regimeIcon string
	switch mode.Regime {
	case "NORMAL":
		regimeIcon = "🟢"
	case "STRESSED":
		regimeIcon = "🟡"
	case "CRISIS":
		regimeIcon = "🔴"
	default:
		regimeIcon = "⚪"
	}

	regimeText := fmt.Sprintf("%s %s", regimeIcon, mode.Regime)
	g.regimeLabel.SetText(regimeText)

	// Eigenvalue details
	eigenText := fmt.Sprintf(
		"Max Eigenvalue: %.4f  |  Condition Number: %.2f\n"+
			"Dominant Eigenvalues: ",
		mode.MaxEigen,
		mode.Condition,
	)

	topN := 6
	if len(mode.Eigenvalues) < topN {
		topN = len(mode.Eigenvalues)
	}

	for i := 0; i < topN; i++ {
		eigenText += fmt.Sprintf("%.3f ", mode.Eigenvalues[i])
	}

	g.eigenLabel.SetText(eigenText)
}

func (g *GUI) updateMatrix() {
	mat := g.eng.Matrix()
	if mat == nil || len(mat.Cor) == 0 {
		g.matrixText.SetText("Waiting for sufficient data...\n(Need at least 2 data points per symbol)")
		return
	}

	var sb strings.Builder

	// Limit display to 10x10
	n := len(mat.Cor)
	displayN := n
	if displayN > 10 {
		displayN = 10
	}

	// Header row
	sb.WriteString("         ")
	for j := 0; j < displayN; j++ {
		sb.WriteString(fmt.Sprintf("%-8s", truncate(mat.Symbols[j], 7)))
	}
	sb.WriteString("\n")

	// Matrix rows
	for i := 0; i < displayN; i++ {
		sb.WriteString(fmt.Sprintf("%-8s ", truncate(mat.Symbols[i], 7)))
		for j := 0; j < displayN; j++ {
			val := mat.Cor[i][j]
			if i == j {
				sb.WriteString("  1.000  ")
			} else {
				sb.WriteString(fmt.Sprintf(" %6.3f ", val))
			}
		}
		sb.WriteString("\n")
	}

	if n > 10 {
		sb.WriteString(fmt.Sprintf("\n(Showing top 10×10 of %d×%d matrix)", n, n))
	}

	g.matrixText.SetText(sb.String())
}

func (g *GUI) updateAlerts() {
	alerts := g.eng.Alerts()

	if len(alerts) == 0 {
		g.alertText.SetText("✓ No alerts - Market operating within normal parameters")
		return
	}

	var sb strings.Builder

	// Show most recent 20 alerts
	start := 0
	if len(alerts) > 20 {
		start = len(alerts) - 20
	}

	sb.WriteString(fmt.Sprintf("Total Alerts: %d (showing most recent)\n\n", len(alerts)))

	// Display in reverse chronological order
	for i := len(alerts) - 1; i >= start; i-- {
		alert := alerts[i]

		var icon string
		switch alert.Level {
		case "CRITICAL":
			icon = "🔴"
		case "HIGH":
			icon = "🟠"
		default:
			icon = "⚠️ "
		}

		timeStr := alert.Time.Format("15:04:05")

		sb.WriteString(fmt.Sprintf(
			"%s [%s] [%s] %s: %s (%.4f > %.4f)\n",
			icon,
			timeStr,
			alert.Level,
			truncate(alert.Symbol, 20),
			alert.Message,
			alert.Value,
			alert.Threshold,
		))
	}

	g.alertText.SetText(sb.String())
}

func (g *GUI) updateStats(fps float64, uptime time.Duration) {
	mat := g.eng.Matrix()
	mode := g.eng.Mode()
	alerts := g.eng.Alerts()

	var symbolCount, dataPoints int
	if mat != nil {
		symbolCount = len(mat.Symbols)
	}

	statsText := fmt.Sprintf(
		"Uptime: %s  |  Update Rate: %.1f FPS  |  Symbols: %d  |  Alerts: %d",
		formatDuration(uptime),
		fps,
		symbolCount,
		len(alerts),
	)

	if mode != nil {
		statsText += fmt.Sprintf("  |  Regime: %s", mode.Regime)
	}

	g.statsLabel.SetText(statsText)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}