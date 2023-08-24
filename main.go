package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/asters1/tools"
	"github.com/gin-gonic/gin"
)

var (
	m3u8_base_url string
	m3u8_path_url string
	m3u8_url_name string
)

func clean() {
	os.RemoveAll("./m3u8_cache/")
	os.MkdirAll("./m3u8_cache", 0666)
}

func main() {
	clean()
	r := gin.Default()
	r.GET("/:a", func(c *gin.Context) {
		fmt.Println("zsa")
	})
	r.GET("/m3u8", func(c *gin.Context) {
		Url := strings.TrimSpace(c.Query("url"))

		if !(strings.HasPrefix(Url, "http://") || strings.HasPrefix(Url, "https://")) {
			c.JSON(404, gin.H{
				"code":    404,
				"message": "页面不存在!",
			})
			return
		} else {
			m3u8_url_name = Url[strings.LastIndex(Url, "/"):]
		}
		if strings.Index(Url, ".m3u") != -1 {
			clean()
			l, _ := url.Parse(Url)
			if Url[:8] == "https://" {
				m3u8_base_url = "https://" + l.Hostname() + "/"
			} else {
				m3u8_base_url = "http://" + l.Hostname() + "/"
			}
			m3u8_path_url = Url[:strings.LastIndex(Url, "/")+1]
			fmt.Println(m3u8_base_url)
			fmt.Println(m3u8_path_url)
		}

		Method := c.Query("method")
		Header := c.Query("header")
		Data := c.Query("data")
		if Method == "" {
			Method = "get"
		}

		if Url != "" {
			body := tools.RequestClient(Url, Method, Header, Data)
			// fmt.Println(Method, Header, Data)
			ioutil.WriteFile("./m3u8_cache/"+m3u8_url_name, []byte(body), 0666)
			c.Writer.Header().Set("content-type", "application/vnd.apple.mpegurl")
			c.File("./m3u8_cache/" + m3u8_url_name)

		} else {
		}
	})
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(":9979")
}
