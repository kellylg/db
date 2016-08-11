package model

type Post struct {
	Id              	int     	`db:"id"`
	Title           	string  	`db:"title"`
	Description     	string  	`db:"description"`
	Content         	string  	`db:"content"`
	Etype           	string  	`db:"etype"`
	Created         	int     	`db:"created"`
	Updated         	int     	`db:"updated"`
	Display         	string  	`db:"display"`
	Uid             	int     	`db:"uid"`
	Uname           	string  	`db:"uname"`
	Passwd          	string  	`db:"passwd"`
	Views           	int     	`db:"views"`
	Comments        	int     	`db:"comments"`
	Likes           	int     	`db:"likes"`
	Deleted         	int     	`db:"deleted"`
	Year            	int     	`db:"year"`
	Month           	string  	`db:"month"`
	AllowComment    	string  	`db:"allow_comment"`
	Tags            	string  	`db:"tags"`
	Catid           	int     	`db:"catid"`
}
