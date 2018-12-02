package mackerel

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	config = flag.String("config", "./config.json", "mackerel apikey and servicename config file")
)

//Metric has name,time,value
type Metric []struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Value int    `json:"value"`
}

//Client has APIKey and ServiceName
type Client struct {
	APIKey      string `json:"apikey"`
	ServiceName string `json:"servicename"`
}

//NewClient create metric post client
func NewClient() Client {
	c := Client{}
	f, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(f, &c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(c)
	return c
}

//Post metric
func (c Client) Post(m Metric) {
	body, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "https://api.mackerelio.com/api/v0/services/"+c.ServiceName+"/tsdb", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("response status code:", resp.StatusCode)
	}
}
