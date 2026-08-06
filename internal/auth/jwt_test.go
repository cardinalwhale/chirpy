package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidToken(t *testing.T) {
	orig_uuid := uuid.New()
	token, err := MakeJWT(orig_uuid, "megasecret", time.Duration(time.Minute))
	if err != nil {
		t.Error("error in making token")
		t.Log(err)
	}
	uuid, err := ValidateJWT(token, "megasecret")
	if err != nil {
		t.Error("error in validating token")
		t.Log(err)
	}
	if uuid != orig_uuid {
		t.Error("uuids do not match")
	}
}

func TestExpiredToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "megasecret", time.Duration(time.Second*3))
	if err != nil {
		t.Error("error in making token")
		t.Log(err)
	}
	time.Sleep((time.Second * 4))
	uuid, err := ValidateJWT(token, "megasecret")
	if err == nil {
		t.Error("token should be expired")
		t.Log(err)
		t.Log(uuid)
	}
}

func TestWrongSecret(t *testing.T) {
	orig_uuid := uuid.New()

	token, err := MakeJWT(orig_uuid, "secretone", time.Duration(time.Minute))
	if err != nil {
		t.Error("error in making token")
		t.Log(err)
	}

	uuid, err := ValidateJWT(token, "secrettwo")
	if err == nil {
		t.Error("it shouldn't validate the token")
		t.Log(err)
		t.Log(uuid)
	}
}

func TestGetBearerToken(t *testing.T) {
	header := make(http.Header)
	token, err := GetBearerToken(header)
	if err.Error() != "Header doesn't exist" {
		t.Error("no error handling for the non existing headers")
	}
	header.Set("Authorization", "Bearer TOKEN_STRING")
	token, _ = GetBearerToken(header)
	if token != "TOKEN_STRING" {
		t.Error("get bearer token is not working correctly")
	}
}
