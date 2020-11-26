package main

import (
	"github.com/diamondburned/handy"
	"github.com/gotk3/gotk3/gtk"
)

type Page struct {
	*gtk.Box
	Header *handy.HeaderBar
}

func NewPage() (p Page) {
	var err error

	p.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		panic(err)
	}

	p.Header = handy.HeaderBarNew()
	p.Add(p.Header)

	return
}
