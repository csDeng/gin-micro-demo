package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

const port = 53759

func TestThrottle(t *testing.T) {

	for i := 0; i < 10; i++ {
		go func(i int) {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/hello", port))
			if err != nil {
				log.Printf("%d get fail, err=%v \r\n", i, err)
				return
			}
			defer func() {
				if resp != nil {
					resp.Body.Close()
				}
			}()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(i, "读取body失败")
				return
			}

			log.Println(i, string(body))

		}(i)
	}

	time.Sleep(time.Duration(10) * time.Second)
}

func TestRpcThrottle(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/b", port))
			if err != nil {
				log.Printf("%d get fail, err=%v \r\n", i, err)
				return
			}
			defer func() {
				if resp != nil {
					resp.Body.Close()
				}
			}()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(i, "读取body失败")
				return
			}

			log.Println(i, string(body))

		}(i)
	}

	time.Sleep(time.Duration(10) * time.Second)
}
