package model

type Tag struct {
	Id      	int     	`db:"id"`
	Name    	string  	`db:"name"`
	Uid     	int     	`db:"uid"`
	Created 	int     	`db:"created"`
	Times   	int     	`db:"times"`
	RcType  	string  	`db:"rc_type"`
}
