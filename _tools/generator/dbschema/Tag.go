//Do not edit this file, which is automatically generated by the generator.
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
	"time"
)

type Tag struct {
	trans	*factory.Transaction
	objects []*Tag
	
	Id     	uint    	`db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name   	string  	`db:"name" bson:"name" comment:"标签名" json:"name" xml:"name"`
	Uid    	uint    	`db:"uid" bson:"uid" comment:"创建者" json:"uid" xml:"uid"`
	Created	uint    	`db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Times  	uint    	`db:"times" bson:"times" comment:"使用次数" json:"times" xml:"times"`
	RcType 	string  	`db:"rc_type" bson:"rc_type" comment:"关联类型" json:"rc_type" xml:"rc_type"`
}

func (this *Tag) Trans() *factory.Transaction {
	return this.trans
}

func (this *Tag) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *Tag) Objects() []*Tag {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *Tag) NewObjects() *[]*Tag {
	this.objects = []*Tag{}
	return &this.objects
}

func (this *Tag) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("tag").SetModel(this)
}

func (this *Tag) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *Tag) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Tag) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Tag) Add() (pk interface{}, err error) {
	this.Created = uint(time.Now().Unix())
	this.Id = 0
	pk, err = this.Param().SetSend(this).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			this.Id = v
		}
	}
	return
}

func (this *Tag) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Update()
}

func (this *Tag) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		
	},func(){
		this.Created = uint(time.Now().Unix())
	this.Id = 0
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			this.Id = v
		}
	}
	return 
}

func (this *Tag) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetMiddleware(mw).Delete()
}

