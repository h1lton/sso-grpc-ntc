package auth

import (
	"context"
	"google.golang.org/grpc"
	ssov1 "sso-grpc-ntc/pkg/api"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(server *grpc.Server) {
	ssov1.RegisterAuthServer(server, &ServerAPI{})
}

func (s ServerAPI) Login(c context.Context, r *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
func (s ServerAPI) Login(
	c context.Context,
	r *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s ServerAPI) Register(c context.Context, r *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
func (s ServerAPI) Register(
	c context.Context,
	r *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s ServerAPI) IsAdmin(c context.Context, r *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
func (s ServerAPI) IsAdmin(
	c context.Context,
	r *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
