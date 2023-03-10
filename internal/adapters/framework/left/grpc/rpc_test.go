package rpc

import (
	"context"
	"hex/internal/adapters/app/api"
	"hex/internal/adapters/core/arithmetic"
	"hex/internal/adapters/framework/left/grpc/pb"
	"hex/internal/adapters/framework/right/db"
	"hex/internal/ports"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	var err error
	lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	//ports
	var dbaseAdapter ports.DbPort
	var core ports.ArithmeticPort
	var appAdapter ports.APIPort
	var gRPCAdapter ports.GRPCPort

	dbaseAdapter, err = db.NewAdapter("sqlite3", "./db.db")

	if err != nil {
		log.Fatal(err)
	}

	core = arithmetic.NewAdapter()

	appAdapter = api.NewAdapter(dbaseAdapter, core)

	gRPCAdapter = NewAdapter(appAdapter)

	pb.RegisterArithmeticServiceServer(grpcServer, gRPCAdapter)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("test server start error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getGRPCConnection(ctx context.Context, t *testing.T) *grpc.ClientConn {
	// conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	return conn
}

func TestGetAddition(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewArithmeticServiceClient(conn)

	params := &pb.OperationParameters{
		A: 1,
		B: 1,
	}

	answer, err := client.GetAddition(ctx, params)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}

	require.Equal(t, answer.Value, int32(2))
}

func TestGetSubtraction(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 1,
		B: 1,
	}

	answer, err := client.GetSubtraction(ctx, params)

	if err != nil {
		t.Fatalf("Expected: %v, got: %v", nil, err)
	}

	require.Equal(t, answer.Value, int32(0))
}

func TestGetMultiplication(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 1,
		B: 1,
	}

	answer, err := client.GetMultiplication(ctx, params)

	if err != nil {
		t.Fatalf("Expected: %v, got: %v", nil, err)
	}

	require.Equal(t, answer.Value, int32(1))
}

func TestGetDivision(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewArithmeticServiceClient(conn)
	params := &pb.OperationParameters{
		A: 1,
		B: 1,
	}

	answer, err := client.GetAddition(ctx, params)

	if err != nil {
		t.Fatalf("Expected: %v, got: %v", nil, err)
	}

	require.Equal(t, answer.Value, int32(1))
}
