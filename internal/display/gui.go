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
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	eng         *engine.Engine
	cancel      context.CancelFunc
	regimeLabel *widget.Label
	eigenLabel  *widget.Label
	matrixText  *widget.Label
	alertText   *widget.Label
}

func NewGUI(eng *engine.Engine, cancel context.CancelFunc) *GUI {
	return &GUI{
		eng:         eng,
		cancel:      cancel,
		regimeLabel: widget.NewLabel("Initializing..."),
		eigenLabel:  widget.NewLabel(""),
		matrixText:  widget.NewLabel("Loading data..."),
		alertText:   widget.NewLabel("No alerts"),
	}
}

func (g *GUI) Run() {
	a := app.New()
	w := a.NewWindow("MatrixPulse - Real-time Financial Correlation Engine")

	g.regimeLabel.TextStyle = fyne.TextStyle{Bold: true}
	g.matrixText.TextStyle = fyne.TextStyle{Monospace: true}
	g.alertText.TextStyle = fyne.TextStyle{Monospace: true}

	header := widget.NewLabelWithStyle("MATRIXPULSE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	modeBox := container.NewVBox(
		widget.NewLabelWithStyle("Market Regime", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		g.regimeLabel,
		g.eigenLabel,
	)

	matrixScroll := container.NewScroll(g.matrixText)
	matrixScroll.SetMinSize(fyne.NewSize(600, 300))

	matrixBox := container.NewVBox(
		widget.NewLabelWithStyle("Correlation Matrix", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		matrixScroll,
	)

	alertScroll := container.NewScroll(g.alertText)
	alertScroll.SetMinSize(fyne.NewSize(600, 200))

	alertBox := container.NewVBox(
		widget.NewLabelWithStyle("Active Alerts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		alertScroll,
	)

	content := container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewVSplit(
			container.NewVBox(modeBox, widget.NewSeparator(), matrixBox),
			alertBox,
		),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(1000, 700))
	w.SetOnClosed(func() {
		g.cancel()
	})

	go g.updateLoop()
	w.ShowAndRun()
}

func (g *GUI) updateLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		g.update()
	}
}

func (g *GUI) update() {
	mode := g.eng.Mode()
	if mode != nil {
		regimeColor := "🟢"
		if mode.Regime == "STRESSED" {
			regimeColor = "🟡"
		} else if mode.Regime == "CRISIS" {
			regimeColor = "🔴"
		}

		g.regimeLabel.SetText(fmt.Sprintf("%s %s", regimeColor, mode.Regime))

		eigenStr := fmt.Sprintf("Max Eigenvalue: %.4f  |  Condition: %.2f\nTop Eigenvalues: ",
			mode.MaxEigen, mode.Condition)
		for i := 0; i < len(mode.Eigenvalues) && i < 6; i++ {
			eigenStr += fmt.Sprintf("%.3f ", mode.Eigenvalues[i])
		}
		g.eigenLabel.SetText(eigenStr)
	}

	mat := g.eng.Matrix()
	if mat != nil && len(mat.Cor) > 0 {
		var sb strings.Builder
		n := len(mat.Cor)
		if n > 10 {
			n = 10
		}

		sb.WriteString("       ")
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%-8s", truncate(mat.Symbols[j], 7)))
		}
		sb.WriteString("\n")

		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%-6s ", truncate(mat.Symbols[i], 6)))
			for j := 0; j < n; j++ {
				val := mat.Cor[i][j]
				if i == j {
					sb.WriteString(" 1.000  ")
				} else {
					sb.WriteString(fmt.Sprintf(" %6.3f ", val))
				}
			}
			sb.WriteString("\n")
		}

		g.matrixText.SetText(sb.String())
	}

	alerts := g.eng.Alerts()
	if len(alerts) > 0 {
		var sb strings.Builder
		start := 0
		if len(alerts) > 15 {
			start = len(alerts) - 15
		}

		for i := len(alerts) - 1; i >= start; i-- {
			a := alerts[i]
			icon := "⚠️ "
			if a.Level == "CRITICAL" {
				icon = "🔴"
			} else if a.Level == "HIGH" {
				icon = "🟠"
			}

			sb.WriteString(fmt.Sprintf("%s [%s] %s: %s (%.4f)\n",
				icon, a.Level, truncate(a.Symbol, 15), a.Message, a.Value))
		}

		g.alertText.SetText(sb.String())
	} else {
		g.alertText.SetText("No alerts - market stable")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
