package models

import (
	"fmt"
	"github.com/levigross/grequests"
)

func GetData(url string) string {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		fmt.Println("请求失败:", err)
		return "error"
	}
	return resp.String()
}
