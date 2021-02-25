package parse

import (
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var globalProtoMap map[string]*desc.FileDescriptor
var IsCached = true
var lk sync.RWMutex

func init() {
	globalProtoMap = make(map[string]*desc.FileDescriptor)
}

func getProto(path string) *desc.FileDescriptor {
	lk.Lock()
	defer lk.Unlock()

	if IsCached {
		fd, ok := globalProtoMap[path]
		if ok {
			log.Println("getProto path:%v cached", path)
			return fd
		}
	}
	p := protoparse.Parser{}
	fds, err := p.ParseFiles(path)
	if err != nil {
		log.Println("getProto ParseFiles error:%v", err)
		return nil
	}
	log.Println("JsonToPb fd %v, err %v", fds[0], err)
	fd := fds[0]

	if IsCached {
		globalProtoMap[path] = fd
	}

	return fd
}

func JsonToPb(protoPath, messageName string, jsonStr []byte) ([]byte, error) {
	fd := getProto(protoPath)
	msg := fd.FindMessage(messageName)
	dymsg := dynamic.NewMessage(msg)
	err := dymsg.UnmarshalJSON(jsonStr)
	if err != nil {
		log.Println("JsonToPb UnmarshalJSON error:%v", err)
		return nil, nil
	}

	any, err := ptypes.MarshalAny(dymsg)
	if err != nil {
		log.Println("JsonToPb MarshalAny error:%v", err)
		return nil, nil
	}
	return any.Value, nil
}

func PbToJson(protoPath, messageName string, protoData []byte) ([]byte, error) {
	fd := getProto(protoPath)
	msg := fd.FindMessage(messageName)
	dymsg := dynamic.NewMessage(msg)

	err := proto.Unmarshal(protoData, dymsg)
	jsonByte, err := dymsg.MarshalJSON()
	return jsonByte, err
}
