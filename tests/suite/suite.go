package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/sy-miller/protos/gen/go/sso"
	"github.com/sy-miller/sso-grpc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	TEST_CONFIG_FILE_PATH = "../config/local_test.json"
	GRPCHOST              = "localhost"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(TEST_CONFIG_FILE_PATH)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout.ToTimeDuration())

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(GRPCHOST, strconv.Itoa(*cfg.GRPC.Port))
}
