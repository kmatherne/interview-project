package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"github.com/go-redis/redis"
)

const (

	REDISADD = "127.0.0.1"
	REDISPORT = ":8888"
	REDISPASS = "siege87751"
	REDISDB = "0"
	REDISSTACK = "endpoint_query"
	REDISENDPOINT = "endpoint:"
	REDISDATA = "data:"
)

func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISADD + REDISPORT,
		Password: REDISPASS,
		DB:       REDISDB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func formatRequest(endpoint, data map[string]string) ([]http.Request, error) {
	requests := []http.Request{}
	for _, value := range data {
		strUrl := endpoint["url"]
		var objMap map[string]string
		err := json.Unmarshal([]byte(value), &objMap)
		if err != nil {
			return requests, err
		}
		for k, v := range objMap {
			strUrl = strings.Replace(strUrl, "{"+k+"}", url.QueryEscape(v), -1)
		}
		req, err := http.NewRequest(endpoint["method"], strUrl, nil)
		if err != nil {
			return requests, err
		}
		requests = append(requests, *req)
	}
	return requests, nil
}

func main() {
	client, err := newRedisClient()
	if err != nil {
		log.Print(err.Error())
		return
	}
	for {
		result, err := client.LPop(REDISSTACK).Result()
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Print(result)
		endpoint, err := client.HGetAll(REDISENDPOINT + result).Result()
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Print(endpoint)
		data, err := client.HGetAll(REDISDATA + result).Result()
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Print(data)
		requests, err := formatRequest(endpoint, data)
		if err != nil {
			log.Print(err.Error())
			return
		}
		c := new(http.Client)
		for _, req := range requests {
			resp, err := c.Do(&req)
			if err != nil {
				log.Print(err.Error())
				return
			}
			log.Print(resp)
		}
	}
}

