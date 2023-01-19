package router

import (
	ctl "wba-bc-project-05/backend/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ctl *ctl.Controller) (*Router, error) {
	r := &Router{ct: ctl}
	return r, nil
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (p *Router) Idx() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(CORS())

	v1 := e.Group("/v1")
	{
		token := v1.Group("/token")
		{
			token.GET("/balance", p.ct.GetBalance)
			token.GET("/welcome", p.ct.Welcome)
		}
		game := v1.Group("/game")
		{
			game.POST("", p.ct.CreateGame)
			game.GET("", p.ct.GetGames)
		}
		bet := v1.Group("/bet")
		{
			bet.POST("/away", p.ct.BetAway)
			bet.POST("/home", p.ct.BetHome)
		}
		vote := v1.Group("/vote")
		{
			vote.POST("/away", p.ct.BetAway)
			vote.POST("/home", p.ct.BetHome)
		}
		login := v1.Group("/login")
		{
			login.POST("/signup", p.ct.SignUp)
			login.POST("/signin", p.ct.SignIn)
			login.POST("/logout", p.ct.LogOut)
		}
	}

	return e
}
