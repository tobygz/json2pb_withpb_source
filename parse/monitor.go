package parse

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var CONF_PATH string = "/Users/yuandan_15/Documents/test/protojson/json2pb_withpb_source/conf/"

func initAllConfFile() error {
	mpNow := make(map[string]string)
	err := GetAllFile(CONF_PATH, mpNow)
	if err != nil {
		//todo, add log
		return err
	}
	var newInfo *ConfInfo

	//log.Println("initAllConfFile 1111")
	for name, _ := range mpNow {
		newInfo, err = GetConfInfo(name)
		if err != nil {
			//todo, fatal error, update failed
			log.Println("GetConfInfo name:%s err: %v", name, err)
			continue
		}
		oldInfo, ok := allConfFileMap[name]
		if !ok || !newInfo.IsSame(oldInfo) {
			//load it
			log.Println("load file name:", name, " md5:", newInfo.Md5String)
			updateProto(name)
			//update map
			allConfFileMap[name] = newInfo
		}
	}
	return nil
}

func GetAllFile(pathname string, mp map[string]string) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".proto") {
			mp[fi.Name()] = fi.Name()
		}
	}
	return err
}

func GetConfInfo(name string) (*ConfInfo, error) {
	confInfo := &ConfInfo{}
	pathname := fmt.Sprintf("%s/%s", CONF_PATH, name)
	confInfo.Name = name
	confInfo.FilePath = pathname
	var err error
	confInfo.LastModSec, confInfo.Md5String, err = GetFileModTimeMD5(confInfo.FilePath)
	return confInfo, err
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetFileModTimeMD5(path string) (lastUpdate int64, md5Str string, err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return
	}
	lastUpdate = fi.ModTime().Unix()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	md5Str = md5V(string(content))
	return
}
