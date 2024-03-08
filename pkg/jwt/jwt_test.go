package jwt

import (
	"a-project-backend/svc/pkg/domain/model/exception"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	t.Run("JWT_INVALID", func(t *testing.T) {
		claims := jwt.StandardClaims{
			Audience:  "TestAudience",
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			Id:        "testID1234",
			IssuedAt:  time.Now().Add(-5 * time.Hour).Unix(),
			Issuer:    "TestIssuer",
			Subject:   "TestSubject",
		}
		issueJWT, err := IssueJWT(claims, "testSecret")
		assert.NoError(t, err)
		newClaims, err := Verify(issueJWT, "testSecret")
		assert.EqualError(t, err, exception.ErrInvalidJWT.Error())
		assert.Nil(t, newClaims)
	})
	t.Run("JWT_VALID", func(t *testing.T) {
		claims := jwt.StandardClaims{
			Audience:  "TestAudience",
			ExpiresAt: time.Now().Add(10 * time.Hour * 24).Unix(),
			Id:        "testID1234",
			IssuedAt:  time.Now().Add(-5 * time.Hour).Unix(),
			Issuer:    "TestIssuer",
			Subject:   "TestSubject",
		}
		issueJWT, err := IssueJWT(claims, "testSecret")
		assert.NoError(t, err)
		newClaims, err := Verify(issueJWT, "testSecret")
		assert.NoError(t, err)
		assert.Equal(t, claims.Audience, newClaims.Audience)
		assert.Equal(t, claims.ExpiresAt, newClaims.ExpiresAt)
		assert.Equal(t, claims.Id, newClaims.Id)
		assert.Equal(t, claims.IssuedAt, newClaims.IssuedAt)
		assert.Equal(t, claims.Issuer, newClaims.Issuer)
		assert.Equal(t, claims.Subject, newClaims.Subject)
	})
}

func TestCreateClaims(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		secret string
	}{
		{
			name:   "test1",
			id:     "testing123",
			secret: "ThisIsSecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := CreateClaims(tt.id, 24*time.Hour, "testIssuer")
			got, err := IssueJWT(claims, tt.secret)
			if err != nil {
				t.Errorf("IssueJWT() error = %v", err)
				return
			}
			verify, err := Verify(got, tt.secret)
			assert.NoError(t, err)
			assert.Equal(t, claims.Id, verify.Id)
		})
	}
}
