package gosync

import(
        "github.com/astaxie/beego/logs"
)

func init(){
     log := logs.NewLogger(10000)
     logs.SetLogger("file", `{"filename":"gosync.log"}`)
}