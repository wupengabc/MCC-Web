package models

import (
	"fmt"
	"github.com/levigross/grequests"
)

func GetData(url string) string {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		fmt.Println(err)
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
