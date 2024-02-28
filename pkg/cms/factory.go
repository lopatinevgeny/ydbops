package cms

import (
	"context"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-genproto/protos/Ydb_Operations"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/ydb-platform/ydb-ops/pkg/options"
)

const (
	BufferSize = 32 << 20
)

type Factory struct {
	grpc     options.GRPC
	auth     options.Auth
	rootOpts options.RootOptions
	user     string
}

func NewConnectionFactory(
	grpc options.GRPC,
	rootOpts options.RootOptions,
	user string,
) *Factory {
	return &Factory{
		grpc:     grpc,
		rootOpts: rootOpts,
		user:     user,
	}
}

func (f Factory) Context() (context.Context, context.CancelFunc) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*time.Duration(f.grpc.TimeoutSeconds))

	t, err := f.auth.Token()
	if err != nil {
		zap.S().Warnf("Failed to load auth token: %v", err)
		return ctx, cf
	}

	return metadata.AppendToOutgoingContext(ctx,
		"x-ydb-auth-ticket", t.Secret,
		"authorization", t.Token()), cf
}

func (f Factory) OperationParams() *Ydb_Operations.OperationParams {
	return &Ydb_Operations.OperationParams{
		OperationMode:    Ydb_Operations.OperationParams_SYNC,
		OperationTimeout: durationpb.New(time.Duration(f.grpc.TimeoutSeconds) * time.Second),
		CancelAfter:      durationpb.New(time.Duration(f.grpc.TimeoutSeconds) * time.Second),
	}
}

func (f Factory) Connection() (*grpc.ClientConn, error) {
  // TODO somewhere here the rootOpts.Auth needs to be used to 
  // supply the necessary credentials and headers to a grpc call
	cr, err := f.Credentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	return grpc.Dial(f.Endpoint(),
		grpc.WithTransportCredentials(cr),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(BufferSize),
			grpc.MaxCallRecvMsgSize(BufferSize)))
}

func (f Factory) Credentials() (credentials.TransportCredentials, error) {
	if !f.grpc.Secure {
		return insecure.NewCredentials(), nil
	}

	return credentials.NewClientTLSFromFile(f.rootOpts.CaFile, "")
}

func (f Factory) Endpoint() string {
  // TODO decide if we want to support multiple endpoints or just one
  // Endpoint in rootOpts will turn from string -> []string in this case
	return fmt.Sprintf("%s:%d", f.rootOpts.Endpoint, f.grpc.Port)
}

func (f Factory) User() string {
	return f.user
}