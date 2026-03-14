package common

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
)

type Key int

// Available colors.
const (
	Selected Key = iota
	DarkerSelected
)

var Colors = map[Key]color.RGBA{
	Selected:       {R: 0x28, G: 0x5e, B: 0x28, A: 0xFF}, // "#285e28"
	DarkerSelected: {R: 0x1a, G: 0x3a, B: 0x1a, A: 0xFF}, // "#1a3a1a"
}

var BgStyles = map[Key]lipgloss.Style{
	Selected:       lipgloss.NewStyle().Background(Colors[Selected]),
	DarkerSelected: lipgloss.NewStyle().Background(Colors[DarkerSelected]),
}

// lipglossColorToHex converts a color.Color to hex string
func LipglossColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
