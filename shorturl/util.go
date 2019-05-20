package shorturl

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)
var ValidChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLSMap struct {
	Data map[string]string
	Domain string
	Count int
	ConfigPath string
	Mu sync.RWMutex
}

type SerialStruct struct {
	data map[string]string
	count int
}

func New(domain , configPath string)*URLSMap{
	return &URLSMap{
		make(map[string]string),
		domain	,
		0,
		path.Join(configPath, "urlsmap"),
		sync.RWMutex{},
	}
}


func (u* URLSMap) LoadFromFile(){
	file, err := os.Open(u.ConfigPath)
	defer file.Close()
	if err != nil{
		log.Println(err)
		u.Data = make(map[string]string)
		return
	}
	dec := json.NewDecoder(file)
	dec.Decode(u)
	fmt.Println(u)
	if err != nil{
		log.Println("=======")
		log.Println(err)
		u.Data = make(map[string]string)
		return
	}
}

func (u* URLSMap) Persist(){
	file, err := os.OpenFile(u.ConfigPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil{
		log.Fatalln(err)
	}
	defer file.Close()
	data ,err := json.Marshal(u)
	file.Write(data)
}

func (u* URLSMap)Set(longurl string) string{
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Count++
	short := parseIntToChars(u.Count)
	shorturl := u.Domain+"/O"+short
	u.Data["/O"+short] = longurl
	return shorturl
}


func (u *URLSMap)Get(shorturl string)(string, bool){
	u.Mu.RLock()
	defer u.Mu.RUnlock()
	longurl , ok:= u.Data[shorturl]
	return longurl, ok
}

func parseIntToChars(count int)string{
	base := len(ValidChar)
	res := []byte{}
	for count>0 {
		res = append(res, ValidChar[count%base])
		count = count/base
	}
	return string(res)
}

