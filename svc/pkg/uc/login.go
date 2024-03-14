package uc

import (
	"a-project-backend/pkg/jwt"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		log.Fatalf("Failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	var result map[string]string

	if err = json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to json unmarshal: %v", err)
	}
	return LoginUseCase{
		certString: result,
	}
}

func (uc LoginUseCase) Do(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// JWTのヘッダを解析し署名に用いられている鍵を取得
	parts := strings.Split(input.JWT, ".")

	// decode the header
	headerJson, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT header: %v", err)
	}

	var header map[string]interface{}
	err = json.Unmarshal(headerJson, &header)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JWT header: %v", err)
	}

	kid := header["kid"].(string)
	cert, ok := uc.certString[kid]
	if !ok {
		return nil, fmt.Errorf("kid not found: %v", kid)
	}

	claims, err := jwt.Verify(input.JWT, cert)
	if err != nil {
		return nil, err
	}
	if claims.Issuer != "https://securetoken.google.com/a-projcect-frontend" {
		return nil, fmt.Errorf("invalid issuer: %v", claims.Issuer)
	}
	return &LoginOutput{
		UserID: claims.Subject,
	}, nil
}
