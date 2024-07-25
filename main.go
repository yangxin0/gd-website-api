package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

    // DeepL free account API
    r.GET("/deepl", func(c *gin.Context) {
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

    r.GET("/youdao", func(c *gin.Context) {
        c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
            "Text": "Not implemented",
        })
    })

    r.GET("/openai", func(c *gin.Context) {
        c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
            "Text": "Not implemented",
        })
    })

    r.GET("/google", func(c *gin.Context) {
        c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
            "Text": "Not implemented",
        })
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
