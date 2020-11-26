package main

import (
	"context"
	corev1 "murmur/gen/core"
	foundationv1 "murmur/gen/foundation"
	"net"

	"github.com/diamondburned/handy"
	"github.com/gotk3/gotk3/gtk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
		conn, err := grpc.Dial("localhost:2289", grpc.WithInsecure(), grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			netConn, err := net.Dial("tcp", s)
			if err != nil {
				return nil, err
			}

			return netConn, nil
		}))
		if err != nil {
			panic(err)
		}

		c := corev1.NewCoreServiceClient(conn)
		f := foundationv1.NewFoundationServiceClient(conn)

		println("Logging in...")
		session, err := f.Login(context.Background(), &foundationv1.LoginRequest{
			Login: &foundationv1.LoginRequest_Local_{
				Local: &foundationv1.LoginRequest_Local{
					Email:    TextOf(emailEntry),
					Password: []byte(TextOf(passwordEntry)),
				},
			},
		})
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		ctx = metadata.AppendToOutgoingContext(ctx, "auth", session.SessionToken)

		println("Obtaining guild list...")
		list, err := c.GetGuildList(ctx, &corev1.GetGuildListRequest{})
		if err != nil {
			panic(err)
		}

		for _, guild := range list.Guilds {
			println(guild.GuildId)
		}
	})
	box.Add(loginButton)

	clamp.Add(box)
	clamp.SetVAlign(gtk.ALIGN_CENTER)
	clamp.SetVExpand(true)
	loginPage.Add(clamp)

	MainWindow.Push(loginPage)
	MainWindow.ShowAll()
}
