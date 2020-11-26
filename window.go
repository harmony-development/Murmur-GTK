package main

import (
	"time"

	"github.com/diamondburned/handy"
)

var MainWindow MurmurWindow

// MurmurWindow is the main window
type MurmurWindow struct {
	*handy.Window
	deck      *handy.Deck
	pageStack []Page
}

func (m *MurmurWindow) Push(p Page) {
	m.deck.Add(p)
	m.deck.ShowAll()
	m.deck.SetVisibleChild(p)
	m.pageStack = append(m.pageStack, p)
}

func (m *MurmurWindow) Pop() {
	m.deck.Navigate(handy.NavigationDirectionBack)

	go func() {
		time.Sleep(time.Duration(m.deck.GetTransitionDuration()) * time.Millisecond)

		item := m.pageStack[len(m.pageStack)-1]

		m.pageStack = m.pageStack[:len(m.pageStack)-1]

		m.deck.Remove(item)
		item.Destroy()
	}()
}

func (m *MurmurWindow) Replace(p Page) {
	current := m.pageStack[len(m.pageStack)-1]
	m.pageStack = m.pageStack[:len(m.pageStack)-1]

	m.Push(p)

	go func() {
		time.Sleep(time.Duration(m.deck.GetTransitionDuration()) * time.Millisecond)

		m.deck.Remove(current)
		current.Destroy()
	}()
}

func NewMurmurWindow() MurmurWindow {
	win := MurmurWindow{
		Window: handy.WindowNew(),
		deck:   handy.DeckNew(),
	}

	win.Add(win.deck)

	return win
}
