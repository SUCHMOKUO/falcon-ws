package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const apiURL = "https://tool.lu/ip/ajax.html"

type httpRes struct {
	Text struct{
		Location string `json:"ip2region_location"`
	} `json:"text"`
}

var httpResPool = &sync.Pool{
	New: func() interface{} {
		return new(httpRes)
	},
}

func genFormData(host string) url.Values {
	v := url.Values{}
	v.Set("ip", host)
	return v
}

func getLocationFromAPI(host string) (string, error) {
	resp, err := http.PostForm(apiURL, genFormData(host))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	jsonRes := httpResPool.Get().(*httpRes)
	defer httpResPool.Put(jsonRes)
	json.Unmarshal(body, jsonRes)
	return jsonRes.Text.Location, nil
}
