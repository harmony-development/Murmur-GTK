package main

import (
	"github.com/diamondburned/handy"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	gtk.Init(nil)
	handy.Init()

	MainWindow = NewMurmurWindow()
	MainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})
}

func main() {
	DecideLogin()
	gtk.Main()
}
