package filenode

import (
	"fmt"
	"image/color"
	"path/filepath"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/bluekeyes/go-gitdiff/gitdiff"

	"github.com/dlvhdr/diffnav/pkg/config"
	"github.com/dlvhdr/diffnav/pkg/icons"
	"github.com/dlvhdr/diffnav/pkg/utils"
)

// Icon style constants.
const (
	IconsNerdStatus   = "nerd-fonts-status"
	IconsNerdSimple   = "nerd-fonts-simple"
	IconsNerdFiletype = "nerd-fonts-filetype"
	IconsNerdFull     = "nerd-fonts-full"
	IconsUnicode      = "unicode"
	IconsASCII        = "ascii"
)

type FileNode struct {
	File       *gitdiff.File
	Depth      int
	YOffset    int
	Selected   bool
	PanelWidth int
	Cfg        config.Config
	ShowFullPath bool
}

func (f *FileNode) Path() string {
	return GetFileName(f.File)
}

func (f *FileNode) Value() string {
	name := filepath.Base(f.Path())
	if f.ShowFullPath {
		name = f.Path()
	}

	// filetype and full use: [status letter/icon] [file-type icon] [filename]
	if f.Cfg.UI.Icons == IconsNerdFull || f.Cfg.UI.Icons == IconsNerdFiletype {
		return utils.RemoveReset(f.renderFullLayout(name))
	}

	// All other styles: [icon] [filename] with optional coloring
	return utils.RemoveReset(f.renderStandardLayout(name))
}

// renderStandardLayout renders: [icon colored] [filename]
// Used by status, simple, filetype, unicode, ascii.
func (f *FileNode) renderStandardLayout(name string) string {
	icon := f.getIcon() + " "
	iconWidth := lipgloss.Width(icon) + 1

	stats := ""
	if f.Cfg.UI.ShowDiffStats {
		stats = " " + ViewFileDiffStats(f.File, lipgloss.NewStyle())
	}

	nameMaxWidth := f.PanelWidth - f.Depth - iconWidth - lipgloss.Width(stats)
	truncatedName := utils.TruncateString(name, nameMaxWidth)
	coloredIcon := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(icon)

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	if f.Selected {
		nameStyle = nameStyle.Bold(true)
		if f.PanelWidth > 0 {
			availableWidth := f.PanelWidth - iconWidth - f.Depth
			if availableWidth > 0 {
				nameStyle = nameStyle.Width(availableWidth)
			}
		}
	}
	return coloredIcon + nameStyle.Render(truncatedName) + stats
}

// renderFullLayout renders: [status icon colored] [file-type icon colored] [filename]
// All icons colored by git status.
func (f *FileNode) renderFullLayout(name string) string {
	statusIcon := f.getStatusIcon()
	fileIcon := icons.GetIcon(name, false)
	statusStyle := lipgloss.NewStyle().Foreground(f.StatusColor())
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	stats := ""
	if f.Cfg.UI.ShowDiffStats {
		stats = " " + ViewFileDiffStats(f.File, lipgloss.NewStyle())
	}

	iconsPrefix := statusStyle.Render(statusIcon) + " " + iconStyle.Render(fileIcon) + " "
	iconsWidth := lipgloss.Width(statusIcon) + 1 + lipgloss.Width(fileIcon) + 1

	nameMaxWidth := f.PanelWidth - f.Depth - iconsWidth - lipgloss.Width(stats)
	truncatedName := utils.TruncateString(name, nameMaxWidth)

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	if f.Selected {
		nameStyle = nameStyle.Bold(true)
		if f.PanelWidth > 0 {
			if w := f.PanelWidth - iconsWidth - f.Depth; w > 0 {
				nameStyle = nameStyle.Width(w)
			}
		}
	}
	return iconsPrefix + nameStyle.Render(truncatedName) + stats
}

// getIcon returns the left icon based on the icon style.
func (f *FileNode) getIcon() string {
	name := filepath.Base(f.Path())
	switch f.Cfg.UI.Icons {
	case IconsNerdStatus:
		if f.File.IsNew {
			return ""
		} else if f.File.IsDelete {
			return ""
		}
		return ""
	case IconsNerdSimple:
		return ""
	case IconsNerdFiletype:
		return icons.GetIcon(name, false) // File-type specific icon (colored by status)
	case IconsUnicode:
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "⛌"
		}
		return "●"
	default: // ascii (fallback for unknown values)
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "x"
		}
		return "*"
	}
}

// getStatusIcon returns the git status indicator (used by full layout).
func (f *FileNode) getStatusIcon() string {
	if f.Cfg.UI.Icons == IconsNerdFull {
		if f.File.IsNew {
			return "\uf457"
		} else if f.File.IsDelete {
			return "\ueadf"
		}
		return "\uf459"
	}
	// Colored letters for filetype and other styles
	if f.File.IsNew {
		return "A"
	} else if f.File.IsDelete {
		return "D"
	}
	return "M"
}

// StatusColor returns the color for this file based on its git status.
func (f *FileNode) StatusColor() color.Color {
	if f.File.IsNew {
		return lipgloss.Green
	} else if f.File.IsDelete {
		return lipgloss.Red
	}
	return lipgloss.Yellow
}

func (f *FileNode) String() string {
	return f.Value()
}

func (f *FileNode) Children() tree.Children {
	return tree.NodeChildren(nil)
}

func (f *FileNode) Hidden() bool {
	return false
}

func (f *FileNode) SetHidden(bool) {}

func (f *FileNode) SetValue(any) {}

func DiffStats(file *gitdiff.File) (int64, int64) {
	if file == nil {
		return 0, 0
	}
	var added int64 = 0
	var deleted int64 = 0
	frags := file.TextFragments
	for _, frag := range frags {
		added += frag.LinesAdded
		deleted += frag.LinesDeleted
	}
	return added, deleted
}

func ViewDiffStats(added, deleted int64, base lipgloss.Style) string {
	addedView := ""
	deletedView := ""

	if added > 0 {
		addedView = base.Foreground(lipgloss.Green).Render(fmt.Sprintf("+%d", added))
	}

	if added > 0 && deleted > 0 {
		addedView += base.Render(" ")
	}

	if deleted > 0 {
		deletedView = base.Foreground(lipgloss.Red).Render(fmt.Sprintf("-%d", deleted))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, addedView, deletedView)
}

func ViewFileDiffStats(file *gitdiff.File, base lipgloss.Style) string {
	added, deleted := DiffStats(file)

	return ViewDiffStats(added, deleted, base)
}

func GetFileName(file *gitdiff.File) string {
	if file.NewName != "" {
		return file.NewName
	}
	return file.OldName
}
