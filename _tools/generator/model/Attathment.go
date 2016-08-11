package model

type Attathment struct {
	Id            	int     	`db:"id"`
	Name          	string  	`db:"name"`
	Path          	string  	`db:"path"`
	Extension     	string  	`db:"extension"`
	Type          	string  	`db:"type"`
	Size          	int64   	`db:"size"`
	Uid           	int     	`db:"uid"`
	Deleted       	int     	`db:"deleted"`
	Created       	int     	`db:"created"`
	Audited       	int     	`db:"audited"`
	RcId          	int     	`db:"rc_id"`
	RcType        	string  	`db:"rc_type"`
	Tags          	string  	`db:"tags"`
}
