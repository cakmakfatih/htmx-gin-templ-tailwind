package main

import (
	"context"
	"fmt"
	"htmxgo/core"
	pages "htmxgo/views/pages"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "public, max-age=86400")
		}

		c.Next()
	}
}

func userMiddleware(c *gin.Context) {
	sessionStr, err := c.Cookie("Session")

	if err != nil {
		c.Next()

		return
	}

	sessionArr := strings.Split(sessionStr, "|")

	accessToken := sessionArr[0]

	user, err := core.DbClient.Auth.User(context.Background(), accessToken)

	if err != nil {
		c.Next()
	}

	c.Set("user", user)
}

func main() {
	core.LoadEnv()
	core.InitDb()

	router := gin.Default()

	router.Use(staticCacheMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	ginHtmlRenderer := router.HTMLRender
	router.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	err := router.SetTrustedProxies(nil)

	if err != nil {
		log.Fatal(err)
	}

	router.Static("/assets", "./public")

	router.GET("/api/auth/sign-in", func(c *gin.Context) {
		discordResponse, err := core.DbClient.Auth.SignInWithProvider(supabase.ProviderSignInOptions{Provider: "discord", RedirectTo: "http://localhost:3000/api/auth/callback", FlowType: supabase.PKCE})

		if err != nil {
			c.Redirect(http.StatusNotAcceptable, "/")

			return
		}

		c.Header("HX-Redirect", discordResponse.URL)
		c.SetCookie("Code-Verifier", discordResponse.CodeVerifier, 3600, "/", "localhost", false, true)

		c.Status(http.StatusAccepted)
	})

	router.GET("/api/auth/callback", func(c *gin.Context) {
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
	})

	router.GET("/", func(c *gin.Context) {
		r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, pages.Home())
		c.Render(http.StatusOK, r)
	})

	router.GET("/user-info", userMiddleware, func(c *gin.Context) {
		_, isAuthenticated := c.Get("user")

		if !isAuthenticated {
			c.HTML(http.StatusOK, "index", gin.H{
				"title": "not authenticated",
			})

			return
		}

		c.HTML(http.StatusOK, "index", gin.H{
			"title": "congratz",
		})
	})

	router.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT")))
}
