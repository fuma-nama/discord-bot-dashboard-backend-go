package controllers

import (
	"discord-bot-dashboard-backend-go/discord"
	"discord-bot-dashboard-backend-go/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func AuthController(jwtConfig jwt.Config, auth discord.OAuth2Config, router *gin.RouterGroup) {
	router.GET("/auth", jwt.AuthMiddleware(jwtConfig), func(c *gin.Context) {
		principal := jwt.Principal(c)

		if discord.CheckToken(principal.AccessToken) {
			c.JSON(http.StatusOK, principal.AccessToken)
		} else {
			jwt.InvalidateSession(c)
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
	})

	router.GET("/login", func(c *gin.Context) {
		params := url.Values{
			"client_id":     {auth.ClientId},
			"scope":         {auth.Scope},
			"response_type": {"code"},
			"redirect_uri":  {auth.RedirectUrl},
		}

		c.Redirect(http.StatusFound, "https://discord.com/oauth2/authorize?"+params.Encode())
		c.Abort()
	})

	router.GET("/callback", func(c *gin.Context) {

		if code := c.Query("code"); code != "" {
			tokenData, err := discord.GetToken(auth, auth.RedirectUrl, code)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}

			user, err := discord.GetUser(tokenData.AccessToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}

			jwtToken, err := jwt.GenerateToken(jwtConfig, tokenData.AccessToken, user.Id, tokenData.ExpiresIn)

			jwt.SetSession(c, jwtToken, tokenData.ExpiresIn)
		}

		c.Redirect(http.StatusFound, auth.ClientUrl)
		c.Abort()
	})

	router.POST("/auth/signout", jwt.AuthMiddleware(jwtConfig), func(c *gin.Context) {
		principal := jwt.Principal(c)
		_ = discord.RevokeToken(auth, principal.AccessToken)

		jwt.InvalidateSession(c)
		c.AbortWithStatus(http.StatusOK)
	})
}
