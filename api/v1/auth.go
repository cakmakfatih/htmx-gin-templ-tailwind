package api

import (
	"context"
	"htmxgo/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

func discordSignInCallback(c *gin.Context) {
	cookie, err := c.Cookie("Code-Verifier")

	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")

		return
	}

	code := c.Request.URL.Query().Get("code")

	user, err := core.DbClient.Auth.ExchangeCode(context.Background(), supabase.ExchangeCodeOpts{AuthCode: code, CodeVerifier: cookie})

	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")

		return
	}

	c.SetCookie("Code-Verifier", "", -1, "/", "localhost", false, true)
	c.SetCookie("Session", user.AccessToken+"|"+user.RefreshToken, 3600, "/", "localhost", false, true)

	c.Redirect(http.StatusPermanentRedirect, "/")
}

func signInWithDiscord(c *gin.Context) {
	discordResponse, err := core.DbClient.Auth.SignInWithProvider(supabase.ProviderSignInOptions{Provider: "discord", RedirectTo: "http://localhost:3000/api/v1/auth/discord-callback", FlowType: supabase.PKCE})

	if err != nil {
		c.Redirect(http.StatusNotAcceptable, "/")

		return
	}

	c.Header("HX-Redirect", discordResponse.URL)
	c.SetCookie("Code-Verifier", discordResponse.CodeVerifier, 3600, "/", "localhost", false, true)

	c.Status(http.StatusAccepted)
}

func RegisterAuth(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/auth")

	group.GET("/sign-in-with-discord", signInWithDiscord)
	group.GET("/discord-callback", discordSignInCallback)
}
