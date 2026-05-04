package application

import "context"

type AuthService struct{}

func (a *AuthService) Register(ctx context.Context, username, password string) (RegisterResponse, error) {
	return RegisterResponse{}, nil
}

func (a *AuthService) Login(ctx context.Context, username, password string) (LoginResponse, error) {
	return LoginResponse{}, nil
}

type UserService struct{}

func (u *UserService) GetAll(ctx context.Context) ([]UserDTO, error) {
	return nil, nil
}
