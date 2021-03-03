package parse

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

func updateProto(fname string) {
	fpath := fmt.Sprintf("%s/%s", CONF_PATH, fname)
	p := protoparse.Parser{}
	fds, err := p.ParseFiles(fpath)
	if err != nil {
		log.Println("loadProto ParseFiles error:%v", err)
		return
	}
	fd := fds[0]

	lk.Lock()
	defer lk.Unlock()
	globalProtoMap[fpath] = fd
}

func loadProto(fname string) *desc.FileDescriptor {
	lk.Lock()
	defer lk.Unlock()

	fpath := fmt.Sprintf("%s/%s", CONF_PATH, fname)
	fd, ok := globalProtoMap[fpath]
	if ok {
		log.Println("loadProto path:%s cached", fpath)
		return fd
	}

	p := protoparse.Parser{}
	fds, err := p.ParseFiles(fpath)
	if err != nil {
		log.Println("loadProto ParseFiles error:%v", err)
		return nil
	}
	log.Println("JsonToPb fd %v, err %v", fds[0], err)
	fd = fds[0]

	globalProtoMap[fpath] = fd
	return fd
}

func JsonToPb(protoFname, messageName string, jsonStr []byte) ([]byte, error) {
	fd := loadProto(protoFname)
	msg := fd.FindMessage(messageName)
	dymsg := dynamic.NewMessage(msg)
	err := dymsg.UnmarshalJSON(jsonStr)
	if err != nil {
		log.Println("JsonToPb UnmarshalJSON error:%v", err)
		return nil, err
	}
	jsval, _ := dymsg.MarshalJSON()
	log.Println("JsonToPb jsval:", string(jsval))
	any, err := ptypes.MarshalAny(dymsg)
	if err != nil {
		log.Println("JsonToPb MarshalAny error:%v", err)
		return nil, nil
	}
	return any.Value, nil
}

func PbToJson(protoPath, messageName string, protoData []byte) ([]byte, error) {
	fd := loadProto(protoPath)
	msg := fd.FindMessage(messageName)
	dymsg := dynamic.NewMessage(msg)

	err := proto.Unmarshal(protoData, dymsg)
	jsonByte, err := dymsg.MarshalJSON()
	return jsonByte, err
}
