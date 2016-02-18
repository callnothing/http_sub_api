/*
* @Author: Alang Wang
* @Date:   2016-01-31 18:39:40
* @Last Modified by:   Alang Wang
* @Last Modified time: 2016-01-31 18:39:40
 */
package main

import (
	"log"
	"strings"

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

func srtTovtt(lines string) {
	// srtTimePattern := "(\\d){2,}:(\\d){2,}:(\\d){2,},(\\d){3}"
	// srtPattern := srtTimePattern + " --> " + srtTimePattern
	// counter := 0
}

func getsub(token string, moviename string) (subcontent string) {
	client, _ := xmlrpc.NewClient("http://api.opensubtitles.org/xml-rpc", nil)
	defer client.Close()

	request := []interface{}{
		token,
		[]struct {
			MovieName string `xmlrpc:"query"`
			Language  string `xmlrpc:"sublanguageid"`
		}{{"ant man", "eng"}}}

	result := struct {
		Status    string `xmlrpc:"status"`
		Subtitles []struct {
			FileName  string `xmlrpc:"SubFileName"`
			Hash      string `xmlrpc:"SubHash"`
			Format    string `xmlrpc:"SubFormat"`
			MovieName string `xmlrpc:"MovieName"`
			Downloads string `xmlrpc:"SubDownloadsCnt"`
			URL       string `xmlrpc:"ZipDownloadLink"`
			Page      string `xmlrpc:"SubtitlesLink"`
		} `xmlrpc:"data"`
	}{}
	// dict := []interface{}{map[string]string{"query": moviename, "sublanguageid": "eng"}}
	if err := client.Call("SearchSubtitles", request, &result); err != nil {
		log.Println("打印错误")
		log.Fatal(err)
	}
	basevtturl := "http://dl.opensubtitles.org/en/download/subformat-vtt/filead/src-api/"
	params := strings.Split(result.Subtitles[0].URL, "/")
	basevtturl = basevtturl + params[len(params)-3] + "/" + params[len(params)-2] + "/" + params[len(params)-1]
	log.Println(result.Subtitles[0].URL)
	log.Println(params)
	log.Println(basevtturl)
	return basevtturl
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
