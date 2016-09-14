package factory

import (
	"log"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type Transaction struct {
	sqlbuilder.Tx
	*Cluster
	*Factory
}

func (t *Transaction) Database(param *Param) db.Database {
	if t.Tx != nil {
		return t.Tx
	}
	if t.Cluster == nil {
		t.Cluster = t.Factory.Cluster(param.Index)
	}
	if param.ReadOrWrite == R {
		return t.Cluster.R()
	}
	return t.Cluster.W()
}

func (t *Transaction) Backend(param *Param) sqlbuilder.Backend {
	return t.Database(param).(sqlbuilder.Backend)
}

func (t *Transaction) Result(param *Param) db.Result {
	return t.C(param).Find(param.Args...)
}

func (t *Transaction) C(param *Param) db.Collection {
	return t.Database(param).Collection(t.Cluster.Table(param.Collection))
}

// ================================
// API
// ================================

// Read ==========================

func (t *Transaction) SelectAll(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param)
	if param.Size > 0 {
		selector = selector.Limit(param.Size).Offset(param.GetOffset())
	}
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.All(param.ResultData)
}

func (t *Transaction) SelectOne(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param).Limit(1)
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.One(param.ResultData)
}

func (t *Transaction) SelectList(param *Param) (func() int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param).Limit(param.Size).Offset(param.GetOffset())
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	countFn := func() int64 {
		cnt, err := t.SelectCount(param)
		if err != nil {
			log.Println(err)
		}
		return cnt
	}
	return countFn, t.joinSelect(param, selector).All(param.ResultData)
}

func (t *Transaction) SelectCount(param *Param) (int64, error) {
	if param.Total > 0 {
		return param.Total, nil
	}

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return param.Total, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	counter := struct {
		Count int64 `db:"_t"`
	}{}
	selector := t.Backend(param).Select(db.Raw("count(1) AS _t")).From(t.Table(param.Collection))
	selector = t.joinSelect(param, selector)
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	selector = selector.Offset(0).Limit(1).OrderBy()
	if err := selector.Iterator().One(&counter); err != nil {
		if err == db.ErrNoMoreRows {
			return 0, nil
		}
		return 0, err
	}
	param.Total = counter.Count
	return counter.Count, nil
}

func (t *Transaction) joinSelect(param *Param, selector sqlbuilder.Selector) sqlbuilder.Selector {
	if param.Joins == nil {
		return selector
	}
	for _, join := range param.Joins {
		coll := t.Table(join.Collection)
		if len(join.Alias) > 0 {
			coll += ` AS ` + join.Alias
		}
		switch strings.ToUpper(join.Type) {
		case "LEFT":
			selector = selector.LeftJoin(coll)
		case "RIGHT":
			selector = selector.RightJoin(coll)
		case "CROSS":
			selector = selector.CrossJoin(coll)
		case "INNER":
			selector = selector.FullJoin(coll)
		default:
			selector = selector.FullJoin(coll)
		}
		if len(join.Condition) > 0 {
			selector = selector.On(join.Condition)
		}
	}
	return selector
}

func (t *Transaction) Select(param *Param) sqlbuilder.Selector {
	selector := t.Backend(param).Select(param.Cols...).From(t.Table(param.Collection))
	return t.joinSelect(param, selector)
}

func (t *Transaction) CheckCached(param *Param) bool {
	if t.Factory.cacher != nil {
		if param.MaxAge > 0 {
			return true
		}
		if param.MaxAge < 0 {
			err := t.Factory.cacher.Del(param.CachedKey())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return false
}

func (t *Transaction) Cached(param *Param, fn func(*Param) error) error {
	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	return fn(param)
}

func (t *Transaction) All(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	res := t.Result(param)
	if param.Size > 0 {
		res = res.Limit(param.Size).Offset(param.GetOffset())
	}
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.All(param.ResultData)
}

func (t *Transaction) List(param *Param) (func() int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	var res db.Result
	if param.Middleware == nil {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := t.Result(param)
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		res = t.Result(param).Limit(param.Size).Offset(param.GetOffset())
	} else {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := param.Middleware(t.Result(param)).OrderBy()
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		res = param.Middleware(t.Result(param).Limit(param.Size).Offset(param.GetOffset()))
	}
	return param.CountFunc, res.All(param.ResultData)
}

func (t *Transaction) One(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.One(param.ResultData)
}

func (t *Transaction) Count(param *Param) (int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return param.Total, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	var cnt uint64
	var err error

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	cnt, err = res.Count()
	param.Total = int64(cnt)
	return param.Total, err
}

// Write ==========================

func (t *Transaction) Insert(param *Param) (interface{}, error) {
	param.ReadOrWrite = W
	return t.C(param).Insert(param.SaveData)
}

func (t *Transaction) Update(param *Param) error {
	param.ReadOrWrite = W
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Update(param.SaveData)
}

func (t *Transaction) Delete(param *Param) error {
	param.ReadOrWrite = W
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Delete()
}
