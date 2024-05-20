package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var gsf singleflight.Group

var errorNotExist = errors.New("not exist")

func main() {
	var wg sync.WaitGroup
	wg.Add(20)

	//模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}

	time.Sleep(4 * time.Second)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

// 获取数据
func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//模拟从db中获取数据
		v, err, _ := gsf.Do(key, func() (interface{}, error) {

			data, err := getDataFromDB(key)
			if err != nil {
				log.Fatalln("error")
			}
			isCache = true
			return data, nil
			//set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}

		//TOOD: set cache
		data = v.(string)
	} else if err != nil {
		return "", err
	}
	return data, nil
}

var isCache = false

// 模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	if !isCache {
		return "", errorNotExist
	}
	log.Printf("get %s from cache\n", key)

	return "cache:" + key, nil
}

// 模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database\n", key)
	time.Sleep(3 * time.Second)
	return "data:" + key, nil
}
