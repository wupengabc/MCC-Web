package models

import (
	"github.com/levigross/grequests"
)

func GetData(url string) string {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return "error"
	}
	return resp.String()
}

func PostData(url string, data string) string {
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		JSON: data,
	})
	if err != nil {
		return "error"
	}
	return resp.String()
}
