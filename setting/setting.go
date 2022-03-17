package setting

import (
	"awesomeProject2/utils"
	"encoding/json"
	"fmt"
	"os"
)

type Setting struct {
	Address string
	Port    string
	DbHost  string
	DbPort  string
	DbUser  string
	DbPass  string
	DbName  string
}

var options Setting

func Load(filename string) *Setting {
	bytes, e := utils.LoadFile(filename)
	e = json.Unmarshal(bytes, &options)
	if e != nil {
		fmt.Println(e)
		return nil
	}
	return &options
}

func Save(filename string, s *Setting) {
	bytes, e := json.Marshal(s)
	if e != nil {
		fmt.Println(e)
		return
	}
	e = os.WriteFile(filename, bytes, 0777)
	if e != nil {
		fmt.Println(e)
	}
}
