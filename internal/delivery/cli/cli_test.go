package cli

import (
	"bytes"
	"context"
	"github.com/Dendyator/AntiBF/internal/delivery/grpc/proto/pb"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) CheckAuthorization(ctx context.Context, in *pb.AuthRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error) {
	panic("implement me")
}

func (m *MockClient) AddToBlacklist(ctx context.Context, in *pb.ListRequest, opts ...grpc.CallOption) (*pb.ListResponse, error) {
	panic("implement me")
}

func (m *MockClient) RemoveFromBlacklist(ctx context.Context, in *pb.ListRequest, opts ...grpc.CallOption) (*pb.ListResponse, error) {
	panic("implement me")
}

func (m *MockClient) AddToWhitelist(ctx context.Context, in *pb.ListRequest, opts ...grpc.CallOption) (*pb.ListResponse, error) {
	panic("implement me")
}

func (m *MockClient) RemoveFromWhitelist(ctx context.Context, in *pb.ListRequest, opts ...grpc.CallOption) (*pb.ListResponse, error) {
	panic("implement me")
}

func (m *MockClient) ResetBucket(ctx context.Context, req *pb.ResetRequest, opts ...grpc.CallOption) (*pb.ResetResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.ResetResponse), args.Error(1)
}

var createClientFunc = func() (*grpc.ClientConn, pb.AntiBruteForceClient) {
	return nil, new(MockClient)
}

func createTestCmd(cmd *cobra.Command, args []string) (*bytes.Buffer, *bytes.Buffer) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)
	return outBuf, errBuf
}

func TestResetBucketCmd(t *testing.T) {
	mockClient := new(MockClient)
	appLogger = logger.New("test")

	createClientFunc = func() (*grpc.ClientConn, pb.AntiBruteForceClient) {
		return nil, mockClient
	}

	mockClient.On("ResetBucket", mock.Anything, &pb.ResetRequest{Login: "test-login", Ip: "127.0.0.1"}).
		Return(&pb.ResetResponse{Success: true}, nil)

	cmd := &cobra.Command{
		Use: "reset",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client := createClientFunc()
			res, err := client.ResetBucket(context.Background(), &pb.ResetRequest{Login: args[0], Ip: args[1]})
			if err != nil {
				return err
			}
			cmd.Println("Reset bucket succeeded:", res.Success)
			return nil
		},
	}
	outBuf, errBuf := createTestCmd(cmd, []string{"test-login", "127.0.0.1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}

	mockClient.AssertExpectations(t)
	assert.Contains(t, outBuf.String(), "Reset bucket succeeded: true")
	assert.Empty(t, errBuf.String())
}
