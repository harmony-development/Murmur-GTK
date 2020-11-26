package main

import (
	"murmur/client"

	"github.com/diamondburned/handy"
	"github.com/gotk3/gotk3/gtk"
)

func DecideLogin() {
	loginPage := NewPage()

	clamp := handy.ClampNew()
	clamp.SetHExpand(true)
	clamp.SetMaximumSize(600)
	clamp.SetTighteningThreshold(500)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)
	if err != nil {
		panic(err)
	}
	box.Add(NewLabel("Login to Harmony"))

	homeserverEntry := NewEntry()
	homeserverEntry.SetPlaceholderText("Homeserver URL")
	box.Add(homeserverEntry)

	emailEntry := NewEntry()
	emailEntry.SetPlaceholderText("Email")
	box.Add(emailEntry)

	passwordEntry := NewEntry()
	passwordEntry.SetPlaceholderText("Password")
	box.Add(passwordEntry)

	loginButton := NewButton()
	loginButton.SetLabel("Login")
	loginButton.Connect("clicked", func() {
		cli, err := client.NewClient(TextOf(homeserverEntry), TextOf(emailEntry), TextOf(passwordEntry))
		if err != nil {
			panic(err)
		}

		guilds, err := cli.GuildList()
		if err != nil {
			panic(err)
		}

		guildPage := NewPage()
		for _, guild := range guilds {
			guildPage.Add(NewLabel(guild.Name))
		}
		MainWindow.Replace(guildPage)
	})
	box.Add(loginButton)

	clamp.Add(box)
	clamp.SetVAlign(gtk.ALIGN_CENTER)
	clamp.SetVExpand(true)
	loginPage.Add(clamp)

	MainWindow.Push(loginPage)
	MainWindow.ShowAll()
}
