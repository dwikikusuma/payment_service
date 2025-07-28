package internalGrpc

import (
	"context"
	"google.golang.org/grpc"
	"payment_service/proto/userpb"
	"time"
)

type UserClient struct {
	Client userpb.UserServiceClient
}

func NewUserClient() *UserClient {
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		panic("Failed to connect to gRPC server: " + err.Error())
	}
	client := userpb.NewUserServiceClient(conn)
	return &UserClient{
		Client: client,
	}
}

func (uc *UserClient) GetUserByUserId(ctx context.Context, userId int64) (*userpb.GetUserInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := &userpb.GetUserInfoRequest{
		UserId: userId,
	}
	response, err := uc.Client.GetUserByUserId(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
