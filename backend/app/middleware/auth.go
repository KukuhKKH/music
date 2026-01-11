package middleware

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

// package-level config pointer; set once at startup via SetConfig to avoid repeated parsing
var cfg *config.Config

// SetConfig sets the package-level configuration pointer.
func SetConfig(c *config.Config) {
	cfg = c
}

func Protected() fiber.Handler {
	conf := cfg
	if conf == nil {
		conf = config.NewConfig()
	}

	if conf.Middleware.Jwt.Secret == "" {
		panic("JWT secret is not set")
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(conf.Middleware.Jwt.Secret),
		ErrorHandler: jwtError,
		Claims:       &JWTClaims{},
		TokenLookup:  "header:Authorization,cookie:" + conf.Cookie.Name,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}

	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

type JWTClaims struct {
	Token  string `json:"token"`
	Type   string `json:"type"`
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokenAccess(userID uint64) (*JWTClaims, error) {
	conf := cfg
	if conf == nil {
		conf = config.NewConfig()
	}

	mySigningKey := []byte(conf.Middleware.Jwt.Secret)

	expiration := conf.Middleware.Jwt.Expiration
	if expiration < 1000000 {
		expiration = expiration * time.Second
	}

	// Create the Claims
	now := time.Now()
	exp := now.Add(expiration)
	claims := &JWTClaims{
		Type:   "Bearer",
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    conf.Middleware.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return nil, err
	}

	claims.Token = ss
	return claims, nil
}
