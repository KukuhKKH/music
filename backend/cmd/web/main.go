package main

import (
	"go.uber.org/fx"

	"git.dev.siap.id/kukuhkkh/app-music/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track"
	"git.dev.siap.id/kukuhkkh/app-music/app/router"
	_ "git.dev.siap.id/kukuhkkh/app-music/docs"
	"git.dev.siap.id/kukuhkkh/app-music/internal/bootstrap"
	"git.dev.siap.id/kukuhkkh/app-music/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
	"git.dev.siap.id/kukuhkkh/app-music/utils/session"
	fxzerolog "github.com/efectn/fx-zerolog"
	_ "go.uber.org/automaxprocs"
)

// @title                       Aplikasi Music API
// @version                     1.0
// @description                 This is a sample API documentation.
// @termsOfService              http://swagger.io/terms/
// @contact.name                Kukuh Rahmadani
// @contact.email               krahmadani1@gmail.com
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @host                        localhost:8080
// @schemes                     http https
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 "Type 'Bearer {TOKEN}' to correctly set the API Key"
// @BasePath                    /
func main() {
	fx.New(
		// config
		fx.Provide(config.NewConfig),
		// logging
		fx.Provide(bootstrap.NewLogger),
		// fiber
		fx.Provide(bootstrap.NewFiber),
		// database
		fx.Provide(database.NewDatabase),
		// session
		fx.Provide(session.NewStore),
		// middleware
		fx.Provide(middleware.NewMiddleware),
		// router
		fx.Provide(router.NewRouter),

		// provide modules
		auth.NewAuthModule,
		track.NewTrackModule,
		dashboard.NewDashboardModule,

		// start aplication
		fx.Invoke(bootstrap.Start),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
