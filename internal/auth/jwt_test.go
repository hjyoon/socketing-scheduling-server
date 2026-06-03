package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signed(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestVerify(t *testing.T) {
	token := signed(t, jwt.MapClaims{"sub": "scheduling", "exp": time.Now().Add(time.Minute).Unix()}, "secret")
	claims, err := Verify(token, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if claims["sub"] != "scheduling" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestVerifyRejectsInvalidToken(t *testing.T) {
	if _, err := Verify("bad", "secret"); err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestUUID(t *testing.T) {
	id := UUID()
	if len(id) != 36 || strings.Count(id, "-") != 4 {
		t.Fatalf("unexpected uuid: %s", id)
	}
}

func TestUUIDFallback(t *testing.T) {
	orig := randomBytes
	randomBytes = func([]byte) (int, error) { return 0, errors.New("no entropy") }
	defer func() { randomBytes = orig }()
	if id := UUID(); len(id) != 32 {
		t.Fatalf("unexpected fallback uuid: %s", id)
	}
}
