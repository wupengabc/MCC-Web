package models

type Notices struct {
	ID      int64  `xorm:"'id' pk autoincr"`
	Message string `xorm:"'message' TEXT"`
}

func GetNotices() []Notices {
	var notices []Notices
	err := Db.Find(&notices)
	if err != nil {
		return nil
	}
	return notices
}
