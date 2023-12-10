package jwt

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type Payload map[string]interface{}

type GenerationData struct {
	Payload Payload
	Expiry  time.Duration
	Key     string
}

func Sign(data GenerationData) (string, time.Time, error) {
	duration := data.Expiry
	if duration == 0 {
		duration = time.Second * 86400
	}
	token := jwt.New()
	timeNow := time.Now()
	expiry := timeNow.Add(time.Second * 86400)
	token.Set(jwt.IssuedAtKey, int32(timeNow.Unix()))
	token.Set(jwt.ExpirationKey, int32(expiry.Unix()))

	for key, value := range data.Payload {
		token.Set(key, value)
	}

	signedToken, err := jwt.Sign(token, jwa.HS256, []byte(data.Key))

	return string(signedToken), expiry, err
}

func Verify(token string, key string) (Payload, error) {
	parsedToken, err := jwt.Parse([]byte(token), jwt.WithVerify(jwa.HS256, []byte(key)), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}

	payload, err := parsedToken.AsMap(context.Background())

	return payload, err
}
