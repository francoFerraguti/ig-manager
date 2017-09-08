package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var router *gin.Engine

var clientID string = "f095653f856647acb4a72b28df0805b1" //https://api.instagram.com/oauth/authorize/?client_id=f095653f856647acb4a72b28df0805b1&redirect_uri=&response_type=code
var clientSecret string = "8d1fc344546d416eac61e37e88e98314"
var redirectURI string = "http://4dadf95c.ngrok.io/instagram/redirect"
var scope string = "likes+comments+relationships+follower_list+public_content"
var grantType string = "authorization_code"
var instagramURL string = "https://api.instagram.com/oauth/authorize/?client_id=" + clientID + "&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=" + scope

func main() {
	router = gin.New()

	router.GET("/instagram/getURL", getURL)
	router.GET("/instagram/redirect", getAuthToken)

	router.Run(":80")
}

func loggetURLin(c *gin.Context) {
	c.String(200, instagramURL)
}

func getAuthToken(c *gin.Context) {
	code := c.Query("code")

	oAuthToken, err := getOAuthToken(code)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	oAuthTokenJSON, _ := json.Marshal(oAuthToken)
	c.JSON(200, json.RawMessage(string(oAuthTokenJSON)))
}

func getOAuthToken(code string) (OAuthToken, error) {
	var oAuthToken OAuthToken

	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	form.Add("grant_type", grantType)
	form.Add("redirect_uri", redirectURI)
	form.Add("code", code)

	response, err := http.Post("https://api.instagram.com/oauth/access_token", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return oAuthToken, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return oAuthToken, err
	}
	if response.StatusCode > 299 {
		return oAuthToken, errors.New(string(body[:]))
	}

	err = json.Unmarshal(body, &oAuthToken)

	return oAuthToken, nil
}

type OAuthToken struct {
	Access_Token string
	User         User
}

type User struct {
	Id              string
	Username        string
	Full_Name       string
	Profile_Picture string
}
