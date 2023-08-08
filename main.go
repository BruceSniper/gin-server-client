package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	viper.SetConfigFile("conf/config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error while loading config file [conf/config.yaml]: %s", err.Error()))
	}

	if viper.GetString("role") == "server" {
		// 启动gin框架，采用默认配置
		router := gin.Default()

		// 编写匿名的handler函数
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "OK",
				"message": "router: test",
			})
		})

		router.GET("/liveness", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "OK",
				"message": "router: liveness",
			})
		})

		router.GET("/readiness", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "OK",
				"message": "router: readiness",
			})
		})
		router.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	} else if viper.GetString("role") == "client" {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic: %v", r)
				}
			}()
		}()

		Get(fmt.Sprintf("%s/test", viper.GetString("server_address")))
		Get(fmt.Sprintf("%s/readiness", viper.GetString("server_address")))
		Get(fmt.Sprintf("%s/liveness", viper.GetString("server_address")))

	}

}

func Get(url string) {
	// 请求数据
	for {
		time.Sleep(2 * time.Second)
		// 创建连接池
		//transport := &http.Transport{
		//	DialContext: (&net.Dialer{
		//		Timeout:   30 * time.Second,
		//		KeepAlive: 30 * time.Second,
		//	}).DialContext,
		//	MaxIdleConns:          100,              // 最大空闲连接数
		//	IdleConnTimeout:       90 * time.Second, // 空闲超时时间
		//	TLSHandshakeTimeout:   10 * time.Second, // tls 握手超时时间
		//	ExpectContinueTimeout: 1 * time.Second,  // 100-continue状态码超时时间
		//}

		// 创建客户端
		//client := &http.Client{
		//	Transport: transport,
		//	Timeout:   30 * time.Second, // 没饿
		//}

		client := http.Client{}
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		start := time.Now() // 记录开始时间
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Sleeping for 60 seconds...")
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			select {
			case <-time.After(60 * time.Second):
				fmt.Println("Awake!")
			case <-c:
				fmt.Println("Interrupted!")
			}
			continue
		}

		// 读取数据
		bds, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Sleeping for 60 seconds...")
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			select {
			case <-time.After(60 * time.Second):
				fmt.Println("Awake!")
			case <-c:
				fmt.Println("Interrupted!")
			}
			continue
		}
		elapsed := time.Since(start) // 计算耗时
		fmt.Printf("访问 %s 返回数据 %s 耗时 %d 毫秒\n", url, string(bds), elapsed.Microseconds())
		resp.Body.Close()
	}
}
