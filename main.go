package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OwO-Network/DeepLX/deepl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

func getProxy() string {
    https_proxy := os.Getenv("https_proxy")
    http_proxy := os.Getenv("http_proxy")
    all_proxy := os.Getenv("all_proxy")
    // socks5
    if all_proxy != "" {
        return all_proxy
    }
    if https_proxy != "" {
        return https_proxy
    }
    if http_proxy != "" {
        return http_proxy
    }

    return ""
}

func main() {
    var config string
    flag.StringVar(&config, "c", "config.ini", "config path")

    cfg, err := ini.Load(config)
    if err != nil {
        fmt.Printf("Fail to load config file: %v\n", err)
        os.Exit(1)
    }
    port := cfg.Section("default").Key("port").MustInt()
    proxyURL := cfg.Section("default").Key("proxy").String()
    proxyEnv := getProxy()
    // overwrite by proxy env
    if proxyEnv !=  "" {
        proxyURL = proxyEnv
    }
	fmt.Printf("Goldendict Website API. Listening on 0.0.0.0:%v\n", port)
	fmt.Printf("    - Listening on 0.0.0.0:%v\n", port)
    if proxyURL != "" {
        fmt.Printf("    - proxy: %v\n", proxyURL)
    } else {
        fmt.Println("    - proxy: disabled")
    }

	// if proxyURL != "" {
	//     proxy, err := url.Parse(proxyURL)
	//     if err != nil {
	//         log.Fatalf("Failed to parse proxy URL: %v", err)
	//     }
	//     http.DefaultTransport = &http.Transport{
	//         Proxy: http.ProxyURL(proxy),
	//     }
	// }

	// Setting the application to release mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
    r.LoadHTMLGlob("templates/*")
	r.Use(cors.Default())

    // DeepL free account API
    if cfg.Section("deepl").Key("enable").MustBool() {
        fmt.Println("    - deepl: enabled")
        r.GET("/deepl", func(c *gin.Context) {
            sourceLang := ""
            targetLang := "ZH"
            translateText := c.Query("gdword")
            // No auth key for free account
            authKey := ""

            result, err := deepl.Translate(sourceLang, targetLang, translateText, authKey, proxyURL)
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
    } else {
        fmt.Println("    - deepl: disabled")
    }

    if cfg.Section("youdao").Key("enable").MustBool() {
        fmt.Println("    - youdao: enabled")
        r.GET("/youdao", func(c *gin.Context) {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": "Not implemented",
            })
        })
    } else {
        fmt.Println("    - youdao: disabled")
    }

    if cfg.Section("openai").Key("enable").MustBool() {
        fmt.Println("    - openai: enabled")
        r.GET("/openai", func(c *gin.Context) {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": "Not implemented",
            })
        })
    } else {
        fmt.Println("    - openai: disabled")
    }

    if cfg.Section("google").Key("enable").MustBool() {
        fmt.Println("    - google: enabled")
        r.GET("/google", func(c *gin.Context) {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": "Not implemented",
            })
        })
    } else {
        fmt.Println("    - google: disabled")
    }

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
		r.Run(fmt.Sprintf(":%v", port))
	}
}
