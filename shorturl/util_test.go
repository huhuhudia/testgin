package shorturl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)



func TestURLSMap_LoadFromFile(t *testing.T) {
	u := New("127.0.0.1:8000", "/home/xwt/programming/packages/src/testgin/config")
	u.LoadFromFile()
	for k, v := range u.Data{
		fmt.Println("k ", k, " v ", v)
	}
}

func TestURLSMap_Set(t *testing.T) {
	u := New("127.0.0.1:8000", "/home/xwt/programming/packages/src/testgin/config")
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i:=0; i<1000000; i++{
			u.Set(fmt.Sprintf("127.0.0.1:8000/asdjiajsd%d",i))
		}
		wg.Done()
	}()
	go func() {
		for i:=0; i<1000000; i++{
			u.Set(fmt.Sprintf("127.0.0.1:8000/asdjiajsd%d",i+100))
		}
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, u.Count, len(u.Data), "the two value should be equal")
	log.Println("count : ", u.Count)
	u.Persist()
}

