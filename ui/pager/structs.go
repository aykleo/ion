package pager

import (
	"github.com/charmbracelet/bubbles/viewport"
)

type Pager struct {
	content     []string
	ready       bool
	viewport    viewport.Model
	currentPath string
}
