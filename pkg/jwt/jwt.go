package jwt

import (
	// "encoding/json"
	// "strconv"
	// "time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		// "exp": json.Number(strconv.FormatInt(time.Now().Add(time.Hour*time.Duration(1)).Unix(), 10)),
		// "iat": json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	})
	stringToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return stringToken, nil
}


func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	
	email := t.Claims.(jwt.MapClaims)["email"]
	return true, &JWTData{
		Email: email.(string),
	}
}