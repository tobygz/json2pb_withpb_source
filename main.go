package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/parse"
	"main/user_info"
	"net/http"
	"os"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	testJsonToPb()
	fmt.Fprintf(w, "htllo")
}

func main() {
	parse.StartMonitor()

	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(":8081", nil)
}

func getFileContent(path string) (body []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	body, err = ioutil.ReadAll(file)
	return
}

func genJson() {
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
}

func test() {
	//path := "/Users/yuandan_15/Documents/test/protojson/json2pb_withpb_source/conf"
	//parse.GetAllFile(path, nil)
	mp := make(map[string]*parse.ConfInfo)
	mp["First"] = &parse.ConfInfo{
		Name: "First",
	}
	mp["Second"] = &parse.ConfInfo{
		Name: "Second",
	}

	mp1 := make(map[string]*parse.ConfInfo)
	for k, v := range mp {
		mp1[k] = v
	}
	mp1["First"].Name = "1111First"
	mpJson, _ := json.Marshal(mp)
	log.Println("result:", string(mpJson))
}

func testJsonToPb() {
	f, err := os.Open("./user.json")
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)

	pbobjBin, err := parse.JsonToPb("user.proto", "user_info.UserInfo", content)
	if err != nil {
		log.Println("testJsonToPb err:", err)
		return
	}

	log.Println(pbobjBin)
}
