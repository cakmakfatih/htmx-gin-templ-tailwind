package middlewares

import (
	"context"
	"htmxgo/core"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

func UserMiddleware(c *gin.Context) {
	sessionStr, err := c.Cookie("Session")

	if err != nil {
		c.Next()

		return
	}

	sessionArr := strings.Split(sessionStr, "|")

	accessToken := sessionArr[0]
	refreshToken := sessionArr[1]

	user, err := core.DbClient.Auth.User(context.Background(), accessToken)

	if err == nil {
		c.Set("user", user)
		return
	}

	user, err = renewToken(c, accessToken, refreshToken)

	if err != nil {
		c.Next()

		return
	}

	c.Set("user", user)
}

func renewToken(c *gin.Context, accessToken string, refreshToken string) (*supabase.User, error) {
	authDetails, err := core.DbClient.Auth.RefreshUser(context.Background(), accessToken, refreshToken)

	if err != nil {
		return nil, err
	}

	c.SetCookie("Session", authDetails.AccessToken+"|"+authDetails.RefreshToken, 3600, "/", "localhost", false, true)
	user, err := core.DbClient.Auth.User(context.Background(), authDetails.AccessToken)

	if err != nil {
		return nil, err
	}

	return user, nil
}
