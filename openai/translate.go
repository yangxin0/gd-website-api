package openai

import (
	"fmt"
    "net/http"
    "context"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	oai "github.com/sashabaranov/go-openai"
)

var appSecret = ""

func TranslateInit(route *gin.Engine, cfg *ini.File) {
    enabled := cfg.Section("openai").Key("enable").MustBool()
    if enabled == false {
        fmt.Println("Dict: OpenAI Disabled")
        return
    }
    fmt.Println("Dict: OpenAI Enabled")
    appSecret = cfg.Section("openai").Key("app_secret").String()

    route.GET("/openai", func(c *gin.Context){
        text := c.Query("gdword")
        result, err := Translate("zh-CN", text)
        if err != nil {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": err,
            })

        } else {
            c.HTML(http.StatusOK, "goldendict.tmpl", gin.H{
                "Text": result,
            })
        }
    })
}

func Translate(targetLang string, text string) (string, error) {
    client := oai.NewClient(appSecret)
    systemPrompt := fmt.Sprintf("You are a highly skilled translation engine with expertise in the technology sector. Your function is to translate texts accurately into the target %s, maintaining the original format, technical terms, and abbreviations. Do not add any explanations or annotations to the translated text.", targetLang)
    prompt := fmt.Sprintf("Translate the following source text to %s, Output translation directly without any additional text.\nSource Text: %s,\nTranslated Text:", targetLang, text)
    resp, err := client.CreateChatCompletion(
        context.Background(),
		oai.ChatCompletionRequest{
            Model: oai.GPT4oMini,
			Messages: []oai.ChatCompletionMessage{
                {
                    Role: oai.ChatMessageRoleSystem,
                    Content: systemPrompt,
                },
				{
					Role:    oai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
        return "", fmt.Errorf("OpenAI Error: %v\n", err)
	}
	return resp.Choices[0].Message.Content, nil
}
