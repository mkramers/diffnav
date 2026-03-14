package common

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	scrollThumb = "┃"
	scrollTrack = "│"
)

// RenderScrollbar renders a vertical scrollbar of the given height, 1 column wide.
func RenderScrollbar(height, totalLines, visibleLines, offset int) string {
	return RenderScrollbarW(height, totalLines, visibleLines, offset, 1)
}

// RenderScrollbarW renders a vertical scrollbar with the given column width.
func RenderScrollbarW(height, totalLines, visibleLines, offset, width int) string {
	if height <= 0 || width <= 0 {
		return ""
	}

	if totalLines <= visibleLines {
		return ""
	}

	trackChar := strings.Repeat(scrollTrack, width)
	thumbChar := strings.Repeat(scrollThumb, width)

	trackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("237"))
	thumbStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("248"))

	thumbSize := max(1, height*visibleLines/totalLines)
	thumbPos := 0
	if totalLines > visibleLines {
		thumbPos = offset * (height - thumbSize) / (totalLines - visibleLines)
	}
	thumbPos = max(0, min(thumbPos, height-thumbSize))

	var sb strings.Builder
	for i := range height {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if i >= thumbPos && i < thumbPos+thumbSize {
			sb.WriteString(thumbStyle.Render(thumbChar))
		} else {
			sb.WriteString(trackStyle.Render(trackChar))
		}
	}
	return sb.String()
}
