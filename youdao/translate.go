package youdao

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yangxin0/gd-website-api/youdao/authv3"
	"gopkg.in/ini.v1"
)

var appKey = ""
var appSecret = ""

type Translation struct {
    Texts []string `json:"translation"`
}

func TranslateInit(route *gin.Engine, cfg *ini.File) {
    enabled := cfg.Section("youdao").Key("enable").MustBool()
    if enabled == false {
        fmt.Println("Dict: Youdao Disabled")
        return
    }
    fmt.Println("Dict: Youdao Enabled")
    appKey = cfg.Section("youdao").Key("app_key").String()
    appSecret = cfg.Section("youdao").Key("app_secret").String()
    route.GET("/youdao", func(c *gin.Context){
        text := c.Query("gdword")
        result := Translate("", "zh-CHS", text)
        if result != "" {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": result,
            })
        } else {
            c.String(http.StatusNotFound, "No Translation")
        }
    })
}

func Translate(sourceLang string, targetLang string, Text string) string {
    if sourceLang == "" {
        sourceLang = "auto"
    }
    params := map[string][]string{
        "from": { sourceLang },
        "to": { targetLang },
        "q": { Text },
    }

    header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	authv3.AddAuthParams(appKey, appSecret, params)
	body := DoPost("https://openapi.youdao.com/api", header, params, "application/json")
    if body == nil {
        return ""
    }
    var result Translation
    if json.Unmarshal(body, &result) != nil {
        return ""
    }
    return result.Texts[0]
}
