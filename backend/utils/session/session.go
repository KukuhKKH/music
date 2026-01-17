package session

import (
	"fmt"
	"time"

	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/google/uuid"
)

func NewStore(cfg *config.Config) *session.Store {
	store := session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s", cfg.Middleware.Session.Name),
		CookieHTTPOnly: true,
		CookieSecure:   false,
		CookieSameSite: "Lax",
		CookiePath:     "/",
		Expiration:     24 * time.Hour,
		Storage:        memory.New(),
		KeyGenerator: func() string {
			return uuid.New().String()
		},
	})

	return store
}

func Get(c *fiber.Ctx, store *session.Store) (*session.Session, error) {
	return store.Get(c)
}
