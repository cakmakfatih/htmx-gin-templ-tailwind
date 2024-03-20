package api

import (
	"context"
	"htmxgo/core"
	"htmxgo/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

func providerAuthCallback(c *gin.Context) {
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

func signInWithProvider(c *gin.Context) {
	provider, providerExists := c.GetQuery("provider")

	if !providerExists {
		c.Status(http.StatusBadRequest)

		return
	}

	if !utils.IsStringIsInSlice(provider, core.Config.Providers) {
		c.Status(http.StatusBadRequest)

		return
	}

	providerResponse, err := core.DbClient.Auth.SignInWithProvider(supabase.ProviderSignInOptions{Provider: provider, RedirectTo: "http://localhost:3000/api/v1/auth/provider/callback", FlowType: supabase.PKCE})

	if err != nil {
		c.Redirect(http.StatusNotAcceptable, "/")

		return
	}

	c.Header("HX-Redirect", providerResponse.URL)
	c.SetCookie("Code-Verifier", providerResponse.CodeVerifier, 3600, "/", "localhost", false, true)

	c.Status(http.StatusAccepted)
}

func registerAuth(apiGroup *gin.RouterGroup) {
	group := apiGroup.Group("/auth")

	group.GET("/with", signInWithProvider)
	group.GET("/provider/callback", providerAuthCallback)
}
