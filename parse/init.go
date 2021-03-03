package parse

import (
	"log"
	"sync"
	"time"

	"github.com/jhump/protoreflect/desc"
)

var globalProtoMap map[string]*desc.FileDescriptor
var IsCached = true
var lk sync.RWMutex

type ConfInfo struct {
	Name       string
	FilePath   string
	Md5String  string
	LastModSec int64
}

func (c *ConfInfo) IsSame(n *ConfInfo) bool {
	return c.Md5String == n.Md5String && c.LastModSec == n.LastModSec
}

var allConfFileMap map[string]*ConfInfo

func init() {
	//反序列化对象
	globalProtoMap = make(map[string]*desc.FileDescriptor)

	//conf文件监控
	allConfFileMap = make(map[string]*ConfInfo)
}
func StartMonitor() {
	tk := time.NewTicker(time.Second * 1)
	go func() {
		initAllConfFile()
		for {
			<-tk.C
			err := initAllConfFile()
			if err != nil {
				log.Println("err in initAllConfFile:", err)
			}
		}
	}()
}
