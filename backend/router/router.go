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
		coin := v1.Group("/coin")
		{
			coin.POST("/transfer", p.ct.TransferCoin)
			coin.POST("/transfer/from", p.ct.TransferCoinFrom)
		}
		token := v1.Group("/token")
		{
			token.GET("/symbol", p.ct.GetTokenSymbol)
			token.GET("/balance/:address", p.ct.GetTokenBalance)
			token.POST("/transfer", p.ct.TransferToken)
			token.POST("/transfer/from", p.ct.TransferTokenFrom)
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
