package model

type Link struct {
	Id      	int     	`db:"id"`
	Name    	string  	`db:"name"`
	Url     	string  	`db:"url"`
	Logo    	string  	`db:"logo"`
	Show    	string  	`db:"show"`
	Verified	int     	`db:"verified"`
	Created 	int     	`db:"created"`
	Updated 	int     	`db:"updated"`
	Catid   	int     	`db:"catid"`
	Sort    	int     	`db:"sort"`
}
