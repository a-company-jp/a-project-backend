package uc

import (
	"a-project-backend/pkg/jwt"
	"context"
)

type LoginUseCase struct {
	certString map[string]string
}

type LoginInput struct {
	JWT string
}

type LoginOutput struct {
	UserID string
}

func NewLoginUseCase() LoginUseCase {
	return LoginUseCase{}
}

func (uc LoginUseCase) Do(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	issuer := "https://securetoken.google.com/a-projcect-frontend"

	sub, err := jwt.Verify(input.JWT, issuer)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		UserID: sub,
	}, nil
}
