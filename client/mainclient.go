package client

import (
	"context"
	"strings"

	corev1 "murmur/gen/core"
	foundationv1 "murmur/gen/foundation"
	profilev1 "murmur/gen/profile"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Event struct {
	*corev1.Event
	Client *Client
}

type MainClient struct {
	Client
	homeserver string
	subclients map[string]*Client
	streams    map[chan *corev1.Event]*Client
}

func NewClient(homeserver, email, password string) (client *MainClient, err error) {
	client = &MainClient{
		Client: Client{},
	}

	client.homeserver = homeserver
	client.subclients = make(map[string]*Client)
	client.streams = make(map[chan *corev1.Event]*Client)

	client.conn, err = grpc.Dial(homeserver, grpc.WithInsecure())
	if err != nil {
		client = nil
		err = errors.Wrap(err, "NewClient: failed to dial grpc")
		return
	}

	client.CoreKit = corev1.NewCoreServiceClient(client.conn)
	client.FoundationKit = foundationv1.NewFoundationServiceClient(client.conn)
	client.Profilekit = profilev1.NewProfileServiceClient(client.conn)

	session, err := client.FoundationKit.Login(context.Background(), &foundationv1.LoginRequest{
		Login: &foundationv1.LoginRequest_Local_{
			Local: &foundationv1.LoginRequest_Local{
				Email:    email,
				Password: []byte(password),
			},
		},
	})
	if err != nil {
		client = nil
		err = errors.Wrap(err, "NewClient: failed to login")
		return
	}

	client.sessionToken = session.SessionToken
	client.userID = session.UserId

	return
}

func (m *MainClient) ClientFor(homeserver string) (*Client, error) {
	if m.homeserver == homeserver || strings.Split(homeserver, ":")[0] == "localhost" {
		return &m.Client, nil
	}

	if val, ok := m.subclients[homeserver]; ok {
		return val, nil
	}

	federatedSession, err := m.FoundationKit.Federate(m.Context(), &foundationv1.FederateRequest{
		Target: homeserver,
	})

	if err != nil {
		return nil, err
	}

	client := new(Client)
	client.conn, err = grpc.Dial(homeserver, grpc.WithInsecure())
	if err != nil {
		err = errors.Wrap(err, "ClientFor: failed to dial grpc")
		return nil, err
	}

	client.CoreKit = corev1.NewCoreServiceClient(client.conn)
	client.FoundationKit = foundationv1.NewFoundationServiceClient(client.conn)
	client.Profilekit = profilev1.NewProfileServiceClient(client.conn)

	session, err := client.FoundationKit.Login(context.Background(), &foundationv1.LoginRequest{
		Login: &foundationv1.LoginRequest_Federated_{
			Federated: &foundationv1.LoginRequest_Federated{
				AuthToken: federatedSession.Token,
				Domain:    m.homeserver,
			},
		},
	})
	if err != nil {
		err = errors.Wrap(err, "ClientFor: failed to login")
		return nil, err
	}

	client.sessionToken = session.SessionToken
	client.userID = session.UserId
	return client, nil
}

type Guild struct {
	Name       string
	ID         uint64
	Homeserver string
	AvatarURL  string
}

func (m *MainClient) GuildList() (ret []Guild, err error) {
	list, err := m.CoreKit.GetGuildList(m.Context(), &corev1.GetGuildListRequest{})
	if err != nil {
		return
	}
	for _, guild := range list.Guilds {
		var client *Client
		client, err = m.ClientFor(guild.Host)
		if err != nil {
			return
		}
		var resp *corev1.GetGuildResponse
		resp, err = client.CoreKit.GetGuild(client.Context(), &corev1.GetGuildRequest{GuildId: guild.GuildId})
		if err != nil {
			return
		}
		ret = append(ret, Guild{
			Name:       resp.GuildName,
			ID:         guild.GuildId,
			Homeserver: guild.Host,
			AvatarURL:  resp.GuildPicture,
		})
	}
	return
}
