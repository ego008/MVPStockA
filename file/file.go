package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func Create(f string) {
	file,err:=os.Create(f)
	if err!=nil{
		fmt.Println(err)
	}
	defer file.Close()
}


func Read(f string) string {
	body,err:=ioutil.ReadFile(f)
	if err!=nil {
		fmt.Println(err)
	}
	return string(body)
}

func Write(f , data string) {
	content:=[]byte(data)
	err:=ioutil.WriteFile(f,content,0777)
	if err!=nil {
		fmt.Println(err)
	}
}