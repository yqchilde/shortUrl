package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type FormParam struct {
	Url string `json:"url"`
}

var regUrl = `scope.short_url = "([?\s+](.*)) ";`

func main() {
	engine := gin.Default()

	tools := engine.Group("/tools")
	{
		tools.GET("/shorturl", FetchView)
		tools.POST("/api/shortUrl", ShortUrl)
	}

	err := engine.Run(":8080")
	if err != nil {
		log.Fatalf("服务器端口启动失败，原因：%s", err)
		return
	}

}

func FetchView(context *gin.Context) {
	f, _ := os.Open("./views/index.html")
	defer f.Close()
	html, _ := ioutil.ReadAll(f)

	context.Header("Content-Type", "text/html; charset=utf-8")
	context.String(http.StatusOK, string(html))
}

func ShortUrl(context *gin.Context) {
	var formParam FormParam
	context.ShouldBind(&formParam)

	urls := strings.Split(formParam.Url, "\n")
	// make slice
	urlMerge := make([]string, 0)
	for _, v := range urls {
		// 二次验证前缀
		if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
			Failed(context, "The url should start with http:// or https://")
			return
		}

		Api := "https://service.weibo.com/share/share.php?url=" + v + "&pic=pic&appkey=key&title=" + v
		resp, err := http.Get(Api)
		if err != nil {
			Failed(context, err.Error())
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			Failed(context, err.Error())
			return
		}

		compile := regexp.MustCompile(regUrl)
		allString := compile.FindAllStringSubmatch(string(body), -1)
		urlShort := allString[0][1]
		if strings.Contains(urlShort, "http://") {
			urlShort = strings.Replace(urlShort, "http://", "https://", 1)
		}
		urlMerge = append(urlMerge, urlShort)
		_ = resp.Body.Close()
	}
	Success(context, urlMerge)
	return
}

// 普通成功返回
func Success(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": v,
	})
}

// 普通失败返回
func Failed(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 201,
		"data": v,
	})
}
