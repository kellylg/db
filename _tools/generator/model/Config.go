package model

type Config struct {
	Id      	int     	`db:"id"`
	Key     	string  	`db:"key"`
	Val     	string  	`db:"val"`
	Updated 	int     	`db:"updated"`
}
