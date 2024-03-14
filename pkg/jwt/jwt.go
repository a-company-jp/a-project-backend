package jwt

import (
	"a-project-backend/svc/pkg/domain/model/exception"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"net/http"
)

var keys map[string]*rsa.PublicKey

func init() {
	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		log.Fatalf("failed to get firebase public key: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	var certs map[string]string

	if err = json.Unmarshal(body, &certs); err != nil {
		log.Fatalf("Failed to json unmarshal: %v", err)
	}

	keys = make(map[string]*rsa.PublicKey)
	for k, v := range certs {
		block, _ := pem.Decode([]byte(v))
		if block == nil {
			continue
		}
		pem.Decode([]byte(v))
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatalf("failed to parse certificate: %v", err)
		}
		rsaPub := cert.PublicKey.(*rsa.PublicKey)
		keys[k] = rsaPub
	}
}

func Verify(j string, issuer string) (string, error) {
	token, err := jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("kid not found")
		}
		key, ok := keys[kid.(string)]
		if !ok {
			return nil, exception.ErrInvalidJWT
		}
		return key, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse jwt: %v", err)
	}
	if err = token.Claims.Valid(); err != nil {
		return "", fmt.Errorf("invalid claims: %v", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	if iss, ok := claims["iss"]; !ok || iss != issuer {
		return "", errors.New("invalid issuer")
	}
	if sub, ok := claims["sub"].(string); !ok || sub == "" {
		return "", errors.New("invalid subject")
	} else {
		return sub, nil
	}
}
