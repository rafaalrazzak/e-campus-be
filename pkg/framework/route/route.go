package route

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Route interface {
	fiber.Router
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
