package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/parse"
	"main/user_info"
	"os"

	"google.golang.org/protobuf/proto"
)

func getFileContent(path string) (body []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	body, err = ioutil.ReadAll(file)
	return
}

func test() {

}

func main() {
	gene_json := false

	if gene_json {
		ui := &user_info.UserInfo{
			Username: "name",
			Age:      uint32(188),
			Graduate: "college",
		}
		ui.BookList = append(ui.BookList, &user_info.Book{
			Name:    "China Book",
			PageNum: 1900,
		})
		ui.BookList = append(ui.BookList, &user_info.Book{
			Name:    "Japan Book",
			PageNum: 1901,
		})
		jsret, _ := json.Marshal(ui)
		fmt.Println("json:", string(jsret))
		return
	}

	//convert json with proto desc to obj
	pbfile := "/Users/yuandan_15/Documents/test/protojson/src/user_info/user.proto"

	jsonContent := `{"username":"name","age":188,"graduate":"college","bookList":[{"name":"China Book","pageNum":1900},{"name":"Japan Book","pageNum":1901}]}`
	pbobjBin, err := parse.JsonToPb(pbfile, "user_info.UserInfo", []byte(jsonContent))
	if err != nil {
		panic(err)
	}

	pbObj := &user_info.UserInfo{}
	err = proto.Unmarshal(pbobjBin, pbObj)
	if err != nil {
		panic(err)
	}
	fmt.Println("last:")
	fmt.Println(pbObj)
}
