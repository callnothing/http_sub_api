/*
* @Author: Alang Wang
* @Date:   2016-01-31 18:39:40
* @Last Modified by:   Alang Wang
* @Last Modified time: 2016-01-31 18:39:40
 */
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kolo/xmlrpc"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func gettoken() (token string) {
	client, _ := xmlrpc.NewClient("http://api.opensubtitles.org/xml-rpc", nil)
	defer client.Close()
	result := struct {
		Token string `xmlrpc:"token"`
	}{}
	client.Call("LogIn", []interface{}{"", "", "", "subtitlesforyoutube"}, &result)
	return result.Token
}

func getsub(token string, moviename string) (subcontent string) {
	client, _ := xmlrpc.NewClient("http://api.opensubtitles.org/xml-rpc", nil)
	defer client.Close()
	// // result := struct {
	// // 	SubDownloadLink string `xmlrpc:"data"`
	// // }{}
	var result string
	dict := []interface{}{map[string]string{"query": moviename, "sublanguageid": "en"}}
	if err := client.Call("SearchSubtitles", []interface{}{token, dict}, &result); err == nil {
		log.Println("打印错误")
		log.Fatal(err)
	}
	log.Println("sub url: " + result)
	return result
}

func main() {
	ostoken := ""
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/searchsubtitle", func(c *gin.Context) {
		// cCp := c.Copy()
		if ostoken == "" {
			ostoken = gettoken()
			log.Println("Token: " + ostoken)
		}
		subdownloadlink := getsub(ostoken, "ant man")
		c.JSON(200, gin.H{
			"message": subdownloadlink,
		})
	})
	r.Run("localhost:9002") // listen and server on 0.0.0.0:8080
}
