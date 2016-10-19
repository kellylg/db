//Do not edit this file, which is automatically generated by the generator.
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
	"time"
)

type Link struct {
	trans	*factory.Transaction
	objects []*Link
	
	Id      	uint    	`db:"id,omitempty,pk" bson:"id,omitempty" comment:"主键ID" json:"id" xml:"id"`
	Name    	string  	`db:"name" bson:"name" comment:"名称" json:"name" xml:"name"`
	Url     	string  	`db:"url" bson:"url" comment:"网址" json:"url" xml:"url"`
	Logo    	string  	`db:"logo" bson:"logo" comment:"LOGO" json:"logo" xml:"logo"`
	Show    	string  	`db:"show" bson:"show" comment:"是否显示" json:"show" xml:"show"`
	Verified	uint    	`db:"verified" bson:"verified" comment:"验证时间" json:"verified" xml:"verified"`
	Created 	uint    	`db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated 	uint    	`db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
	Catid   	uint    	`db:"catid" bson:"catid" comment:"分类" json:"catid" xml:"catid"`
	Sort    	int     	`db:"sort" bson:"sort" comment:"排序" json:"sort" xml:"sort"`
}

func (this *Link) Trans() *factory.Transaction {
	return this.trans
}

func (this *Link) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *Link) Objects() []*Link {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *Link) NewObjects() *[]*Link {
	this.objects = []*Link{}
	return &this.objects
}

func (this *Link) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("link").SetModel(this)
}

func (this *Link) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *Link) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Link) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Link) Add() (interface{}, error) {
	this.Created = uint(time.Now().Unix())
	return this.Param().SetSend(this).Insert()
}

func (this *Link) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	this.Updated = uint(time.Now().Unix())
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Update()
}

func (this *Link) Upsert(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		this.Updated = uint(time.Now().Unix())
	},func(){
		this.Created = uint(time.Now().Unix())
	})
}

func (this *Link) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetMiddleware(mw).Delete()
}

