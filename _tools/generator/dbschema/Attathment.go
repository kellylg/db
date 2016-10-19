//Do not edit this file, which is automatically generated by the generator.
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
	"time"
)

type Attathment struct {
	trans	*factory.Transaction
	objects []*Attathment
	
	Id       	uint    	`db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name     	string  	`db:"name" bson:"name" comment:"文件名" json:"name" xml:"name"`
	Path     	string  	`db:"path" bson:"path" comment:"保存路径" json:"path" xml:"path"`
	Extension	string  	`db:"extension" bson:"extension" comment:"扩展名" json:"extension" xml:"extension"`
	Type     	string  	`db:"type" bson:"type" comment:"文件类型" json:"type" xml:"type"`
	Size     	uint64  	`db:"size" bson:"size" comment:"文件尺寸" json:"size" xml:"size"`
	Uid      	uint    	`db:"uid" bson:"uid" comment:"UID" json:"uid" xml:"uid"`
	Deleted  	uint    	`db:"deleted" bson:"deleted" comment:"被删除时间" json:"deleted" xml:"deleted"`
	Created  	uint    	`db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Audited  	uint    	`db:"audited" bson:"audited" comment:"审核时间" json:"audited" xml:"audited"`
	RcId     	uint    	`db:"rc_id" bson:"rc_id" comment:"关联id" json:"rc_id" xml:"rc_id"`
	RcType   	string  	`db:"rc_type" bson:"rc_type" comment:"关联类型" json:"rc_type" xml:"rc_type"`
	Tags     	string  	`db:"tags" bson:"tags" comment:"标签" json:"tags" xml:"tags"`
}

func (this *Attathment) Trans() *factory.Transaction {
	return this.trans
}

func (this *Attathment) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *Attathment) Objects() []*Attathment {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *Attathment) NewObjects() *[]*Attathment {
	this.objects = []*Attathment{}
	return &this.objects
}

func (this *Attathment) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("attathment").SetModel(this)
}

func (this *Attathment) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *Attathment) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Attathment) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *Attathment) Add() (interface{}, error) {
	this.Created = uint(time.Now().Unix())
	return this.Param().SetSend(this).Insert()
}

func (this *Attathment) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Update()
}

func (this *Attathment) Upsert(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		
	},func(){
		this.Created = uint(time.Now().Unix())
	})
}

func (this *Attathment) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	
	return this.Param().SetMiddleware(mw).Delete()
}

