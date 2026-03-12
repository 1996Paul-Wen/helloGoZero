package infra

import (
	"fmt"
	"io"
	"log"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/config"
	"github.com/elastic/go-elasticsearch/v7"
)

var esClient *elasticsearch.Client

func InitEsClient(conf config.ESConfigModel) {
	cfg := elasticsearch.Config{
		Addresses: conf.Hosts,
	}
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	res, err := esClient.Info()
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	var respBody []byte
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	fmt.Printf("es info is %+v\n", string(respBody))

}

func LoadESClient() *elasticsearch.Client {
	if esClient == nil {
		log.Fatal("es client not init")
	}
	return esClient
}
