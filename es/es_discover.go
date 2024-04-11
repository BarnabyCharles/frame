package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ghodss/yaml"

	"github.com/BarnabyCharles/frame/config"
)

var connectionEsDiscover *elasticsearch.Client

func ConnectionEsDiscover(serverName, group string) error {
	nacos, err := config.GetNacosConfig(serverName, group)
	if err != nil {
		return err
	}
	var esConfig config.AppConfig
	err = yaml.Unmarshal([]byte(nacos), &esConfig)

	if err != nil {
		return errors.New("将yaml文件转换为结构体格式失败！" + err.Error())
	}
	if err != nil {
		return errors.New("链接es失败" + err.Error())
	}
	connectionEsDiscover, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esConfig.EsDiscover.Url},
	})
	if err != nil {
		return errors.New("链接es_discover失败" + err.Error())
	}
	return nil
}

func GetDataByNamesFromEs(index, content string, from, size int) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": content,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  &buf,
	}

	res, err := req.Do(context.Background(), connectionEsDiscover)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
