package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yangxin0/gd-website-api/deepl"
	"github.com/yangxin0/gd-website-api/youdao"
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
    if proxyURL != "" {
        fmt.Printf("Proxy: %v\n", proxyURL)
    } else {
        fmt.Println("Proxy: Disabled")
    }

	// Setting the application to release mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
    r.LoadHTMLGlob("templates/*")
	r.Use(cors.Default())

    deepl.TranslateInit(r, cfg)
    youdao.TranslateInit(r, cfg)

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
