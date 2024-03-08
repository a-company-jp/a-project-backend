package uc

import (
	"a-project-backend/pkg/config"
	"a-project-backend/pkg/jwt"
	"a-project-backend/svc/pkg/domain/model/user"
	"context"
)

type LoginUseCase struct {
	jwtSecret string
}

type LoginInput struct {
	JWT user.JWT
}

type LoginOutput struct {
	UserID user.ID
}

func NewLoginUseCase() LoginUseCase {
	conf := config.Get()
	if conf.Service.Authentication.JwtSecret == "" {
		panic("jwt secret is not set")
	}
	return LoginUseCase{
		jwtSecret: conf.Service.Authentication.JwtSecret,
	}
}

func (uc LoginUseCase) Do(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	claims, err := jwt.Verify(string(input.JWT), uc.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		UserID: user.ID(claims.Id),
	}, nil
}
