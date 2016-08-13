//Generated by webx-top/db
package dbschema

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	
)

type Ocontent struct {
	trans	*factory.Transaction
	
	Id             	int     	`db:"id,omitempty" comment:"ID"`
	RcId           	int     	`db:"rc_id" comment:"关联ID"`
	RcType         	string  	`db:"rc_type" comment:"关联类型"`
	Content        	string  	`db:"content" comment:"博客原始内容"`
	Etype          	string  	`db:"etype" comment:"编辑器类型"`
}

func (this *Ocontent) SetTrans(trans *factory.Transaction) *Ocontent {
	this.trans = trans
	return this
}

func (this *Ocontent) Param() *factory.Param {
	return factory.NewParam(Factory).SetTrans(this.trans).SetCollection("ocontent")
}

func (this *Ocontent) Get(mw func(db.Result) db.Result) error {
	return this.Param().SetResult(this).SetMiddleware(mw).One()
}

func (this *Ocontent) List(mw func(db.Result) db.Result, page, size int) ([]*Ocontent, func() int64, error) {
	r := []*Ocontent{}
	counter, err := this.Param().SetPage(page).SetSize(size).SetResult(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Ocontent) ListByOffset(mw func(db.Result) db.Result, offset, size int) ([]*Ocontent, func() int64, error) {
	r := []*Ocontent{}
	counter, err := this.Param().SetOffset(offset).SetSize(size).SetResult(&r).SetMiddleware(mw).List()
	return r, counter, err
}

func (this *Ocontent) Add(args ...*Ocontent) (interface{}, error) {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSave(data).Insert()
}

func (this *Ocontent) Edit(mw func(db.Result) db.Result, args ...*Ocontent) error {
	var data = this
	if len(args)>0 {
		data = args[0]
	}
	return this.Param().SetSave(data).SetMiddleware(mw).Update()
}

func (this *Ocontent) Delete(mw func(db.Result) db.Result) error {
	return this.Param().SetMiddleware(mw).Delete()
}

