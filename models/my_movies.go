package models

import "time"

type MyMovie struct {
	ID         int       `json:"id" xorm:"id pk autoincr"`
	Name       string    `json:"name" xorm:"name"`
	FirstUrl   string    `json:"first_url" xorm:"first_url"`
	SecondUrl  string    `json:"second_url" xorm:"second_url"`
	Catelog    string    `json:"catelog" xorm:"catelog"`
	UpdateAt   time.Time `json:"update_at" xorm:"update_at"`
	IsDownLoad bool      `json:"is_down_load" xorm:"is_down_load"`
}

func (*MyMovie) TableName() string {
	return "my_movies"
}
