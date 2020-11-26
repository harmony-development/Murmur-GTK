package client

import (
	"context"

	corev1 "murmur/gen/core"
	foundationv1 "murmur/gen/foundation"
	profilev1 "murmur/gen/profile"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn          *grpc.ClientConn
	CoreKit       corev1.CoreServiceClient
	FoundationKit foundationv1.FoundationServiceClient
	Profilekit    profilev1.ProfileServiceClient
	sessionToken  string
	userID        uint64
	onceHandlers  []func(*corev1.Event)
}

func (c Client) Context() context.Context {
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "auth", c.sessionToken)

	return ctx
}
