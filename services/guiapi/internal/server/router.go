package server

import (
	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/fasthttp/router"
)

func (s *Server) initRouter() *router.Router {
	c := cors.DefaultHandler()

	r := router.New()

	r.POST("/auth", c.CorsMiddleware(s.auth))
	r.POST("/reg", c.CorsMiddleware(s.reg))

	r.POST("/addRecipe", c.CorsMiddleware(s.addRecipe))
	r.POST("/editRecipe", c.CorsMiddleware(s.editRecipe))
	r.DELETE("/deleteRecipe", c.CorsMiddleware(s.deleteRecipe))

	r.POST("/getRecipe", c.CorsMiddleware(s.getRecipe))
	r.POST("/getRecipesList", c.CorsMiddleware(s.getRecipesList))

	return r
}
