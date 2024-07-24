package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func authMiddleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Token != "" {
			providedTokenInQuery := c.Query("token")
			providedTokenInHeader := c.GetHeader("Authorization")

			// Compatability with the Bearer token format
			if providedTokenInHeader != "" {
				parts := strings.Split(providedTokenInHeader, " ")
				if len(parts) == 2 {
					if parts[0] == "Bearer" || parts[0] == "DeepL-Auth-Key" {
						providedTokenInHeader = parts[1]
					} else {
						providedTokenInHeader = ""
					}
				} else {
					providedTokenInHeader = ""
				}
			}

			if providedTokenInHeader != cfg.Token && providedTokenInQuery != cfg.Token {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid access token",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func main() {
	cfg := initConfig()

	fmt.Printf("DeepL X has been successfully launched! Listening on 0.0.0.0:%v\n", cfg.Port)
	fmt.Println("Developed by sjlleo <i@leo.moe> and missuo <me@missuo.me>.")

	// Set Proxy
	proxyURL := os.Getenv("PROXY")
	if proxyURL == "" {
		proxyURL = cfg.Proxy
	}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			log.Fatalf("Failed to parse proxy URL: %v", err)
		}
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	if cfg.Token != "" {
		fmt.Println("Access token is set.")
	}
	if cfg.AuthKey != "" {
		fmt.Println("DeepL Official Authentication key is set.")
	}

	// Setting the application to release mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
    r.LoadHTMLGlob("templates/*")
	r.Use(cors.Default())

    // Free API endpoint for GoldenDict, No Pro Account required
    r.GET("/goldendict", authMiddleware(cfg), func(c *gin.Context) {
        sourceLang := ""
        targetLang := "ZH"
		translateText := c.Query("gdword")
		authKey := cfg.AuthKey
		proxyURL := cfg.Proxy

		result, err := translateByDeepLX(sourceLang, targetLang, translateText, authKey, proxyURL)
		if err != nil {
			log.Fatalf("Translation failed: %s", err)
		}

		if result.Code == http.StatusOK {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": result.Data,
            })
		} else {
            c.String(result.Code, result.Message)
		}
	})

	// Catch-all route to handle undefined paths
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})

	envPort, ok := os.LookupEnv("PORT")
	if ok {
		r.Run(":" + envPort)
	} else {
		r.Run(fmt.Sprintf(":%v", cfg.Port))
	}
}
