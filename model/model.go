package model

import (
	"fmt"
	"sync"

	"github.com/xormplus/xorm"
)

type BaseModel struct {
	*sync.Mutex
	engines chan *xorm.Engine
}

func(this *BaseModel)GetEngine() *xorm.Engine{
	engine := <- this.engines
	if err := engine.Ping(); err != nil{
		fmt.Println(err)
	}
	return engine
}

func (this *BaseModel)InitXorm(group string){
	ms := make(map[string]interface{})

	ms["source"] = "mysql 链接地址 "
	//用户名：密码@tcp(地址：端口)/库名
	if len(ms["source"].(string)) == 0{
		fmt.Println("日志打印---地址错误")
	}
	ms["driver"] = "数据库引擎"
	//mysql
	ms["maxidle"] = "最大空闲数"
	ms["naxopen"] = "最大连接数"
	this.engines <- this.initDb(ms)
}

func (this *BaseModel)initDb(ms map[string]interface{}) *xorm.Engine{
	orm,err := xorm.NewEngine(ms["driver"].(string),ms["source"].(string))
	if err != nil{
		fmt.Println("日志打印错误")
		panic(err)
	}
	return orm
}