package model

type Ocontent struct {
	Id      	int     	`db:"id"`
	RcId    	int     	`db:"rc_id"`
	RcType  	string  	`db:"rc_type"`
	Content 	string  	`db:"content"`
	Etype   	string  	`db:"etype"`
}
