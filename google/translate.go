package google

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/translate"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"gopkg.in/ini.v1"
)

var appSecret = ""

func TranslateInit(route *gin.Engine, cfg *ini.File) {
    enabled := cfg.Section("google").Key("enable").MustBool()
    if enabled == false {
        fmt.Println("Dict: Google Disabled")
        return
    }
    fmt.Println("Dict: Google Enabled")
    appSecret = cfg.Section("google").Key("app_secret").String()
    route.GET("/google", func(c *gin.Context){
        text := c.Query("gdword")
        result := Translate("zh-CN", text)
        if result != "" {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": result,
            })
        } else {
            c.String(http.StatusNotFound, "No Translation")
        }
    })
}

func Translate(targetLang string, text string) string {
    ctx := context.Background()

    lang, _ := language.Parse(targetLang)
    opts := option.WithAPIKey(appSecret)
    client, err := translate.NewClient(ctx, opts)
    if err != nil {
        fmt.Printf("xx%v", err)
        return ""
    }
    defer client.Close()

    resp, err := client.Translate(ctx, []string{text}, lang, nil)
    if err != nil  || len(resp) == 0 {
        fmt.Printf("sss%v", err)
        return ""
    }
    return resp[0].Text
}
