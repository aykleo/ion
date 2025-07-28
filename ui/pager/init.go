package pager

import "github.com/charmbracelet/bubbles/viewport"

func InitPager() Pager {
	return Pager{
		content: content,
		ready:   false,
		viewport: viewport.Model{
			Width:  0,
			Height: 0,
		},
	}
}
