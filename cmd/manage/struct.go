package main

import "github.com/pt-main/manage/manager"

type ManageCli struct {
	C         *manager.Manager // Config
	NoSave    bool
	TermWidth int
}
