/*
* @Author: Alang Wang
* @Date:   2016-01-31 18:39:40
* @Last Modified by:   Alang Wang
* @Last Modified time: 2016-01-31 18:39:40
 */
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

func srtTovtt(lines []string) (newlines []string) {
	srtTimePattern := "(\\d){2,}:(\\d){2,}:(\\d){2,},(\\d){3}"
	srtPattern := srtTimePattern + " --> " + srtTimePattern
	vtthead := []string{}
	vtthead = append(vtthead, "", "WEBVTT\n")
	vtt := []string{}
	vtt = append(vtthead, lines...)
	for i := 0; i < len(vtt); i++ {
		if matched, _ := regexp.MatchString(srtPattern, vtt[i]); matched {
			vtt[i] = strings.Replace(vtt[i], ",", ".", -1)
			vtt[i] = strings.Replace(vtt[i], " ", ".", -1)
		}
	}
	return vtt
}

func getsub(token string, moviename string) (subcontent string) {
	client, _ := xmlrpc.NewClient("http://api.opensubtitles.org/xml-rpc", nil)
	defer client.Close()

	request := []interface{}{
		token,
		[]struct {
			MovieName string `xmlrpc:"query"`
			Language  string `xmlrpc:"sublanguageid"`
		}{{moviename, "eng"}}}

	result := struct {
		Status    string `xmlrpc:"status"`
		Subtitles []struct {
			FileName        string `xmlrpc:"SubFileName"`
			Hash            string `xmlrpc:"SubHash"`
			Format          string `xmlrpc:"SubFormat"`
			MovieName       string `xmlrpc:"MovieName"`
			Downloads       string `xmlrpc:"SubDownloadsCnt"`
			URL             string `xmlrpc:"ZipDownloadLink"`
			Page            string `xmlrpc:"SubtitlesLink"`
			SubSumCD        string `xmlrpc:"SubSumCD"`
			SubDownloadLink string `xmlrpc:"SubDownloadLink"`
		} `xmlrpc:"data"`
	}{}
	// dict := []interface{}{map[string]string{"query": moviename, "sublanguageid": "eng"}}
	if err := client.Call("SearchSubtitles", request, &result); err != nil {
		log.Println("打印错误")
		log.Fatal(err)
	}
	basevtturl := "http://dl.opensubtitles.org/en/download/subencoding-utf8/filead/src-api/"
	params := strings.Split(result.Subtitles[0].SubDownloadLink, "/")
	basevtturl = basevtturl + params[len(params)-3] + "/" + params[len(params)-2] + "/" + strings.Split(params[len(params)-1], ".")[0]
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
	r.GET("/searchsubtitle/:moviename/", func(c *gin.Context) {
		// cCp := c.Copy()
		if ostoken == "" {
			ostoken = gettoken()
			log.Println("Token: " + ostoken)
		}
		moviename := c.Param("moviename")
		subdownloadlink := getsub(ostoken, moviename)
		log.Println(subdownloadlink)
		resp, err := http.Get(subdownloadlink)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")
		newlines := srtTovtt(lines)

		// c.JSON(200, gin.H{
		// 	"message": strings.Join(newlines,""),
		// })
		c.String(200, strings.Join(newlines,""))

	})
	r.Run("localhost:9002") // listen and server on 0.0.0.0:8080
}
