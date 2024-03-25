package main

import (
	"errors"
	"log"
	"time"

	"github.com/KudinovKV/FastEmbededCache/cache"
)

func main() {
	driver := cache.New(time.Minute)

	driver.Set("statue of liberty", "40.68960612218659, -74.0456618251789", time.Minute*2)

	data, err := driver.Get("statue of liberty")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	driver.Delete("statue of liberty")

	_, err = driver.Get("statue of liberty")
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			log.Println("element not found")
		}
	}
}
