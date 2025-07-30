package main

import (
	"database/sql"
	"os"

	"github.com/aykleo/ion/config"
	"github.com/aykleo/ion/data"
	pager "github.com/aykleo/ion/ui/pager"
	textinput "github.com/aykleo/ion/ui/text-input"
)

type terminal struct {
	width         int
	height        int
	err           error
	currentFolder string
	config        config.IConfig
	data          data.IData
	db            *sql.DB
	input         textinput.ITextInput
	pager         pager.IPager
}

func getFolderFromOs() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
