//Generated by webx-top/db
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
)

type Link struct {
	trans	*factory.Transaction
	
	Id             	int     	`db:"id,omitempty" comment:"主键ID"`
	Name           	string  	`db:"name" comment:"名称"`
	Url            	string  	`db:"url" comment:"网址"`
	Logo           	string  	`db:"logo" comment:"LOGO"`
	Show           	string  	`db:"show" comment:"是否显示"`
	Verified       	int     	`db:"verified" comment:"验证时间"`
	Created        	int     	`db:"created" comment:"创建时间"`
	Updated        	int     	`db:"updated" comment:"更新时间"`
	Catid          	int     	`db:"catid" comment:"分类"`
	Sort           	int     	`db:"sort" comment:"排序"`
}

func (this *Link) SetTrans(trans *factory.Transaction) *Link {
	this.trans = trans
	return this
}

func (this *Link) Param() *factory.Param {
	return factory.NewParam(Factory).SetTrans(this.trans).SetCollection("link")
}

func (this *Link) Get(mw func(db.Result) db.Result) error {
	return this.Param().SetRecv(this).SetMiddleware(mw).One()
}

func (this *Link) List(mw func(db.Result) db.Result, page, size int) ([]*Link, func() int64, error) {
	r := []*Link{}
	counter, err := this.Param().SetPage(page).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Link) ListByOffset(mw func(db.Result) db.Result, offset, size int) ([]*Link, func() int64, error) {
	r := []*Link{}
	counter, err := this.Param().SetOffset(offset).SetSize(size).SetRecv(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Link) Add(args ...*Link) (interface{}, error) {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).Insert()
}

func (this *Link) Edit(mw func(db.Result) db.Result, args ...*Link) error {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSend(data).SetMiddleware(mw).Update()
}

func (this *Link) Delete(mw func(db.Result) db.Result) error {
	return this.Param().SetMiddleware(mw).Delete()
}

