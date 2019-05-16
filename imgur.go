package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type IUploader interface {
	Upload(interface{}) (interface{}, error)
	Delete(string) error
	Setup() error
}

type ImgurUploader struct {
	Config ImgurConfig
}

type ImgurConfig struct {
	ClientID         string
	ClientSecret     string
	GenTokenUrl      string
	AuthorizationUrl string
	UploadUrl        string
}

type ImgurUploadBody struct {
	Image       []byte `json:"image"`
	Album       string `json:"album,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ImgurResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Status  int         `json:"status"`
}

type ImgurSuccessData struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    int    `json:"datetime"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Size        int    `json:"size"`
	Views       int    `json:"views"`
	Bandwidth   int    `json:"bandwidth"`
	Vote        string `json:"vote"`
	Favorite    bool   `json:"favorite"`
	Nsfw        bool   `json:"nsfw"`
	Section     string `json:"section"`
	DeleteHash  string `json:"deletehash,omitempty"`
	Link        string `json:"link"`
	InGallery   bool   `json:"in_gallery"`
}

type ImgurErrorResponse struct {
	Error   string `json:"error"`
	Request string `json:"request"`
	Method  string `json:"method"`
}

func (iu *ImgurUploader) Upload(iub ImgurUploadBody) (result *ImgurResponse, err error) {
	result = &ImgurResponse{}

	form := url.Values{"image": {base64.StdEncoding.EncodeToString(iub.Image)}}

	req, err := http.NewRequest(http.MethodPost, iu.Config.UploadUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Client-ID "+iu.Config.ClientID)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "api.imgur.com")
	req.Header.Add("content-length", strconv.FormatInt(req.ContentLength, 10))
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		result.Data = ImgurErrorResponse{}
	} else {
		result.Data = ImgurSuccessData{}
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (iu *ImgurUploader) Setup() (err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("foo")
	viper.AddConfigPath("./config/")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	iu.Config.ClientID = viper.GetString("ClientID")
	iu.Config.ClientSecret = viper.GetString("ClientSecret")
	iu.Config.UploadUrl = viper.GetString("AuthorizationUrl")

	return nil
}
