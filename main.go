package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "mcc_web/routers"
)

//func syncData() {
//	servers := models.GetAllServer()
//	for _, server := range servers {
//		str := models.GetData("http://" + server.Server + "/" + server.Key + "/botlist")
//		id := server.Id
//		if str == "error" {
//			continue
//		} else {
//			var data [][]interface{}
//			err := json.Unmarshal([]byte(str), &data)
//			if err != nil {
//				fmt.Println("Error parsing JSON:", err)
//				return
//			}
//			for _, row := range data {
//				models.InsertBot(row[1].(string), row[2].(string), row[3].(string), row[4].(string), row[5].(string), strconv.FormatInt(id, 10), row[6].(string))
//			}
//		}
//	}
//}
//
//func startSync() {
//	go func() {
//		ticker := time.NewTicker(4 * time.Second)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case <-ticker.C:
//				syncData()
//			}
//		}
//	}()
//}

func main() {
	//startSync()
	beego.Run()
}
