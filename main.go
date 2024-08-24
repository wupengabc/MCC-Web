package main

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"mcc_web/models"
	_ "mcc_web/routers"
	"time"
)

func syncData() {
	servers := models.GetAllServer()
	for _, server := range servers {
		str := models.GetData("http://" + server.Server + "/" + server.Key + "/botlist")
		if str == "error" {
			continue
		} else {
			var data [][]interface{}
			err := json.Unmarshal([]byte(str), &data)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}
			for _, row := range data {
				models.InsertBot(row[1].(string), row[2].(string), row[3].(string), row[4].(string), row[5].(string), server.Server, row[6].(string))
			}
		}
	}
}

func startSync() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				syncData()
			}
		}
	}()
}

func main() {
	startSync()
	beego.Run()
}
