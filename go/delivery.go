package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	REDISADD      = "127.0.0.1"
	REDISPORT     = ":8888"
	REDISPASS     = "siege87751"
	REDISDB       = 0
	REDISSTACK    = "endpoint_query"
	REDISENDPOINT = "endpoint:"
	REDISDATA     = "data:"
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

func writeLog(file *os.File, text string) {
	text = text + " - " + time.Now().UTC().String() + "\n"
	_, err := file.Write([]byte(text))
	if err != nil {
		log.Fatal("Could not write log", err)
	}
}

func handleRequest(client *redis.Client, id string, file *os.File) {
	log.Print("Handling requests")
	endpoint, err := client.HGetAll(REDISENDPOINT + id).Result()
	if err != nil {
		writeLog(file, "Error retrieving from "+REDISENDPOINT+id+" : "+err.Error())
		return
	}
	data, err := client.HGetAll(REDISDATA + id).Result()
	if err != nil {
		writeLog(file, "Error retrieving from "+REDISDATA+id+" : "+err.Error())
		return
	}
	requests, err := formatRequest(endpoint, data)
	if err != nil {
		writeLog(file, "Error creating requests : "+err.Error())
		return
	}
	c := new(http.Client)
	for _, req := range requests {
		delivery := time.Now().UTC().String()
		resp, err := c.Do(&req)
		if err != nil {
			writeLog(file, "Error sending request : "+err.Error())
			continue
		}
		responseTime := time.Now().UTC().String()
		body := "Empty Body"
		if resp.StatusCode == 200 {
			byts, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				writeLog(file, "Error writing response body : "+err.Error())
			}
			body = string(byts)
		}
		go writeLog(file, "URL: "+req.URL.String()+" - Method: "+req.Method+" - Delivery Time: "+delivery+" - Response Time: "+responseTime+" - Status: "+resp.Status+" - Body: "+body)
	}
}

func main() {
	file, err := os.Create("logs.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
		return
	}
	defer file.Close()
	client, err := newRedisClient()
	if err != nil {
		log.Fatal("Cannot create Redis Client", err.Error())
		return
	}
	for {
		result, err := client.LPop(REDISSTACK).Result()
		if err != nil && err.Error() != "redis: nil" {
			writeLog(file, "Error retrieving from stack : "+err.Error())
			continue
		}
		if err.Error() == "redis: nil" {
			continue
		}
		go handleRequest(client, result, file)
	}
}

