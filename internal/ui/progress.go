package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type ProgressTracker struct {
	spinner     *spinner.Spinner
	currentStep string
	startTime   time.Time
}

func NewProgressTracker() *ProgressTracker {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	return &ProgressTracker{spinner: s}
}

func (p *ProgressTracker) StartStep(step string) {
	p.currentStep = step
	p.startTime = time.Now()
	p.spinner.Suffix = fmt.Sprintf(" %s...", step)
	p.spinner.Start()
}

func (p *ProgressTracker) EndStep(success bool, message string) {
	p.spinner.Stop()

	duration := time.Since(p.startTime)
	status := color.GreenString("✓")
	if !success {
		status = color.RedString("✗")
	}

	if message != "" {
		fmt.Printf("  %s %s (%s) - %s\n", status, p.currentStep, duration.Round(time.Millisecond), message)
	} else {
		fmt.Printf("  %s %s (%s)\n", status, p.currentStep, duration.Round(time.Millisecond))
	}
}

func (p *ProgressTracker) Info(message string) {
	p.spinner.Stop()
	fmt.Printf("  %s %s\n", color.BlueString("ℹ"), message)
	p.spinner.Start()
}
