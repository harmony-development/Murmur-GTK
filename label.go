package main

import (
	"github.com/gotk3/gotk3/gtk"
)

func NewLabel(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	if err != nil {
		panic(err)
	}

	return label
}

func NewEntry() *gtk.Entry {
	entry, err := gtk.EntryNew()
	if err != nil {
		panic(err)
	}
	return entry
}

func NewButton() *gtk.Button {
	btn, err := gtk.ButtonNew()
	if err != nil {
		panic(err)
	}
	return btn
}

func TextOf(e *gtk.Entry) string {
	data, err := e.GetText()
	if err != nil {
		panic(err)
	}
	return data
}
