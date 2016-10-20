package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/confl"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	//"github.com/webx-top/webx/lib/com"
)

var cfg = &config{}
var memberTemplate = "\t%v\t%v\t`db:\"%v\" bson:\"%v\" comment:\"%v\" json:\"%v\" xml:\"%v\"`"
var replaces = &map[string]string{
	"packageName":  "",
	"imports":      "",
	"structName":   "",
	"attributes":   "",
	"tableName":    "",
	"beforeInsert": "",
	"beforeUpdate": "",
	"beforeDelete": "",
}
var structTemplate = `//Do not edit this file, which is automatically generated by the generator.
package {{packageName}}

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	{{imports}}
)

type {{structName}} struct {
	trans	*factory.Transaction
	objects []*{{structName}}
	
{{attributes}}
}

func (this *{{structName}}) Trans() *factory.Transaction {
	return this.trans
}

func (this *{{structName}}) Use(trans *factory.Transaction) factory.Model {
	this.trans = trans
	return this
}

func (this *{{structName}}) Objects() []*{{structName}} {
	if this.objects == nil {
		return nil
	}
	return this.objects[:]
}

func (this *{{structName}}) NewObjects() *[]*{{structName}} {
	this.objects = []*{{structName}}{}
	return &this.objects
}

func (this *{{structName}}) Param() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetTrans(this.trans).SetCollection("{{tableName}}").SetModel(this)
}

func (this *{{structName}}) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetRecv(this).SetMiddleware(mw).One()
}

func (this *{{structName}}) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetPage(page).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *{{structName}}) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = this.NewObjects()
	}
	return this.Param().SetArgs(args...).SetOffset(offset).SetSize(size).SetRecv(recv).SetMiddleware(mw).List()
}

func (this *{{structName}}) Add() (interface{}, error) {
	{{beforeInsert}}
	return this.Param().SetSend(this).Insert()
}

func (this *{{structName}}) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	{{beforeUpdate}}
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Update()
}

func (this *{{structName}}) Upsert(mw func(db.Result) db.Result, args ...interface{}) error {
	return this.Param().SetArgs(args...).SetSend(this).SetMiddleware(mw).Upsert(func(){
		{{beforeUpdate}}
	},func(){
		{{beforeInsert}}
	})
}

func (this *{{structName}}) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	{{beforeDelete}}
	return this.Param().SetMiddleware(mw).Delete()
}

`

type config struct {
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Host           string          `json:"host"`
	Engine         string          `json:"engine"`
	Database       string          `json:"database"`
	SaveDir        string          `json:"saveDir"`
	Prefix         string          `json:"prefix"`
	PackageName    string          `json:"packageName"`
	Schema         string          `json:"schema"`
	AutoTimeFields *AutoTimeFields `json:"autoTime"`
}

type AutoTimeFields struct {
	//update操作时，某个字段自动设置为当前时间（map的键和值分别为表名称和字段名称。当表名称设置为“*”时，代表所有表中的这个字段）
	Update map[string][]string `json:"update"`

	//insert操作时，某个字段自动设置为当前时间（map的键和值分别为表名称和字段名称。当表名称设置为“*”时，代表所有表中的这个字段）
	Insert map[string][]string `json:"insert"`
}

func main() {
	confFile := flag.String(`c`, ``, `-c conf.yaml`)
	username := flag.String(`u`, `root`, `-u user`)
	password := flag.String(`p`, ``, `-p password`)
	host := flag.String(`h`, `localhost`, `-h host`)
	engine := flag.String(`e`, `mysql`, `-e engine`)
	database := flag.String(`d`, `blog`, `-d database`)
	saveDir := flag.String(`o`, `dbschema`, `-o targetDir`)
	prefix := flag.String(`pre`, ``, `-pre prefix`)
	packageName := flag.String(`pkg`, `dbschema`, `-pkg packageName`)
	schema := flag.String(`schema`, `public`, `-schema schemaName`)
	autoTime := flag.String(`autoTime`, `{"update":{"*":["updated"]},"insert":{"*":["created"]}}`, `-autoTime <json-data>`)
	flag.Parse()

	cfg.Username = *username
	cfg.Password = *password
	cfg.Host = *host
	cfg.Engine = *engine
	cfg.Database = *database
	cfg.SaveDir = *saveDir
	cfg.Prefix = *prefix
	cfg.PackageName = *packageName
	cfg.Schema = *schema

	var err error
	if len(*confFile) > 0 {
		_, err = confl.DecodeFile(*confFile, cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	//com.Dump(cfg)
	var sess sqlbuilder.Database
	switch cfg.Engine {
	case `mymysql`, `mysql`:
		fallthrough
	default:
		settings := mysql.ConnectionURL{
			Host:     cfg.Host,
			Database: cfg.Database,
			User:     cfg.Username,
			Password: cfg.Password,
		}
		sess, err = mysql.Open(settings)
	}
	if err != nil {
		log.Fatal(err)
	}
	if cfg.AutoTimeFields == nil && len(*autoTime) > 0 {
		cfg.AutoTimeFields = &AutoTimeFields{}

		// JSON
		if (*autoTime)[0] == '{' {
			err = json.Unmarshal([]byte(*autoTime), cfg.AutoTimeFields)
			if err != nil {
				log.Fatal(err)
			}
		} else { // update(*:updated)/insert(*:created) 括号内的格式：<表1>:<字段1>,<字段2>,<...字段N>;<表2>:<字段1>,<字段2>,<...字段N>
			cfg.AutoTimeFields.Update = make(map[string][]string)
			cfg.AutoTimeFields.Insert = make(map[string][]string)
			for _, par := range strings.Split(*autoTime, `/`) {
				par = strings.TrimSpace(par)
				switch {
				case strings.HasPrefix(par, `update(`):
					par = strings.TrimPrefix(par, `update(`)
					par = strings.TrimSuffix(par, `)`)
					for _, item := range strings.Split(par, `;`) {
						t := strings.SplitN(item, `:`, 2)
						if len(t) > 1 {
							t[0] = strings.TrimSpace(t[0])
							t[1] = strings.TrimSpace(t[1])
							if len(t[0]) == 0 || len(t[1]) == 0 {
								continue
							}
							if _, ok := cfg.AutoTimeFields.Update[t[0]]; !ok {
								cfg.AutoTimeFields.Update[t[0]] = []string{}
							}
							for _, field := range strings.Split(t[1], `,`) {
								field = strings.TrimSpace(field)
								if len(field) == 0 {
									continue
								}
								cfg.AutoTimeFields.Update[t[0]] = append(cfg.AutoTimeFields.Update[t[0]], field)
							}
						}
					}

				case strings.HasPrefix(par, `insert(`):
					par = strings.TrimPrefix(par, `insert(`)
					par = strings.TrimSuffix(par, `)`)
					for _, item := range strings.Split(par, `;`) {
						t := strings.SplitN(item, `:`, 2)
						if len(t) > 1 {
							t[0] = strings.TrimSpace(t[0])
							t[1] = strings.TrimSpace(t[1])
							if len(t[0]) == 0 || len(t[1]) == 0 {
								continue
							}
							if _, ok := cfg.AutoTimeFields.Insert[t[0]]; !ok {
								cfg.AutoTimeFields.Insert[t[0]] = []string{}
							}
							for _, field := range strings.Split(t[1], `,`) {
								field = strings.TrimSpace(field)
								if len(field) == 0 {
									continue
								}
								cfg.AutoTimeFields.Insert[t[0]] = append(cfg.AutoTimeFields.Insert[t[0]], field)
							}
						}
					}
				}
			}
		}
	}
	defer sess.Close()
	tables, err := sess.Collections()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Found tables: %v`, tables)
	err = os.MkdirAll(cfg.SaveDir, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	allFields := map[string]map[string]factory.FieldInfo{}
	hasPrefix := len(cfg.Prefix) > 0
	for _, tableName := range tables {
		structName := TableToStructName(tableName, cfg.Prefix)
		imports := ``
		goFields, fields := GetTableFields(cfg.Engine, sess, tableName)
		fieldBlock := strings.Join(goFields, "\n")
		noPrefixTableName := tableName
		if hasPrefix {
			noPrefixTableName = strings.TrimPrefix(tableName, cfg.Prefix)
		}
		replaceMap := *replaces
		replaceMap["packageName"] = cfg.PackageName
		replaceMap["structName"] = structName
		replaceMap["attributes"] = fieldBlock
		replaceMap["tableName"] = noPrefixTableName
		replaceMap["beforeInsert"] = ""
		replaceMap["beforeUpdate"] = ""
		replaceMap["beforeDelete"] = ""

		importTime := false
		if cfg.AutoTimeFields != nil {
			_fieldNames, ok := cfg.AutoTimeFields.Insert[`*`]
			if !ok {
				_fieldNames, ok = cfg.AutoTimeFields.Insert[tableName]
			}
			if ok && len(_fieldNames) > 0 {
				beforeInsert := ``
				newLine := ``
				for _, _fieldName := range _fieldNames {
					fieldInf, ok := fields[_fieldName]
					if !ok {
						continue
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `int64`, `uint64`:
						beforeInsert += newLine + `this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				replaceMap["beforeInsert"] = beforeInsert
			}
			_fieldNames, ok = cfg.AutoTimeFields.Update[`*`]
			if !ok {
				_fieldNames, ok = cfg.AutoTimeFields.Update[tableName]
			}
			if ok && len(_fieldNames) > 0 {
				beforeUpdate := ``
				newLine := ``
				for _, _fieldName := range _fieldNames {
					fieldInf, ok := fields[_fieldName]
					if !ok {
						continue
					}
					switch fieldInf.GoType {
					case `uint`, `int`, `int64`, `uint64`:
						beforeUpdate += newLine + `this.` + fieldInf.GoName + ` = ` + fieldInf.GoType + `(time.Now().Unix())`
						newLine = "\n\t"
						importTime = true
					case `string`:
						//TODO
					}
				}
				replaceMap["beforeUpdate"] = beforeUpdate
			}
		}
		if importTime {
			imports += "\n\t" + `"time"`
		}

		replaceMap["imports"] = imports

		content := structTemplate
		for tag, val := range replaceMap {
			content = strings.Replace(content, `{{`+tag+`}}`, val, -1)
		}

		saveAs := filepath.Join(cfg.SaveDir, structName) + `.go`
		file, err := os.Create(saveAs)
		if err == nil {
			_, err = file.WriteString(content)
		}
		if err != nil {
			log.Println(err)
		} else {
			log.Println(`Generated struct:`, structName)
		}

		allFields[noPrefixTableName] = fields
	}

	content := `//Do not edit this file, which is automatically generated by the generator.
package {{packageName}}

import (
	"github.com/webx-top/db/lib/factory"
)

func init(){
	{{initCode}}
}

`
	dataContent := strings.Replace(fmt.Sprintf(`factory.Fields=%#v`+"\n", allFields), `map[string]factory.FieldInfo`, `map[string]*factory.FieldInfo`, -1)
	dataContent = strings.Replace(dataContent, `:factory.FieldInfo`, `:&factory.FieldInfo`, -1)
	content = strings.Replace(content, `{{packageName}}`, cfg.PackageName, -1)
	content = strings.Replace(content, `{{initCode}}`, dataContent, -1)
	saveAs := filepath.Join(cfg.SaveDir, `init`) + `.go`
	file, err := os.Create(saveAs)
	if err == nil {
		_, err = file.WriteString(content)
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println(`Generated init.go`)
	}

	log.Println(`End.`)
}

func TableToStructName(tableName string, prefix string) string {
	if len(prefix) > 0 {
		tableName = strings.TrimPrefix(tableName, prefix)
	}
	return factory.ToCamleCase(tableName)
}

func DataType(fieldInfo *factory.FieldInfo) string {
	switch fieldInfo.DataType {
	case `int`, `tinyint`, `smallint`, `mediumint`:
		if fieldInfo.Unsigned {
			return `uint`
		}
		return `int`
	case `bigint`:
		if fieldInfo.Unsigned {
			return `uint64`
		}
		return `int64`
	case `decimal`, `double`:
		return `float64`
	case `float`:
		return `float32`
	case `bit`, `binary`, `varbinary`, `tinyblob`, `blob`, `mediumblob`, `longblob`: //二进制
		return `byte[]`
	case `geometry`, `point`, `linestring`, `polygon`, `multipoint`, `multilinestring`, `multipolygon`, `geometrycollection`: //几何图形
		return `byte[]`

	//postgreSQL
	case `boolean`:
		return `bool`
	case `oid`:
		if fieldInfo.Unsigned {
			return `uint64`
		}
		return `int64`

	default:
		return `string`
	}
}

func GetTableFields(engine string, d sqlbuilder.Database, tableName string) ([]string, map[string]factory.FieldInfo) {
	switch engine {
	case "mymysql", "mysql":
		fallthrough
	default:
		return getMySQLTableFields(d, tableName)
	}
}
