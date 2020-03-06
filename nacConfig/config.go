package nacConfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
)

const (
	environment = "dev"
	//environment = "test"
	//environment = "pro"
)

func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Println(errors.New("Can not get current file info"))
	}
	return file
}

func Substr(dirctory string, position, length int) string {
	runes := []rune(dirctory)
	l := position + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[position:l])
}

func GetParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func GetConfig() Config {
	JsonParse := NewJsonStruct()
	v := Config{}
	currentPath := CurrentFile()
	parentPath := GetParentDirectory(currentPath)
	JsonParse.Load(parentPath+"/nac-"+environment+".json", &v)
	return v
}
