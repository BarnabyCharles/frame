package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/BarnabyCharles/frame/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ghodss/yaml"
)

var connectionEs *elasticsearch.Client

func ConnectionEs(serverName, group string) error {
	nacos, err := config.GetNacosConfig(serverName, group)
	if err != nil {
		return err
	}
	var esConfig config.AppConfig
	err = yaml.Unmarshal([]byte(nacos), &esConfig)

	if err != nil {
		return errors.New("将yaml文件转换为结构体格式失败！" + err.Error())
	}

	connectionEs, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esConfig.Es.Url},
	})
	if err != nil {
		return errors.New("链接es失败" + err.Error())
	}
	return nil
}

func InsertGoods(index, docID string, doc map[string]interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), connectionEs)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	return nil
}

// 获取商品数据根据id
func GetGoodsById(docID, index string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: docID,
	}

	res, err := req.Do(context.Background(), connectionEs)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, err
	}

	var doc map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		return nil, err
	}

	return doc, err
}

func GetGoodsByName(index, goodsName, from, size string) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": goodsName,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// 创建 Elasticsearch 查询请求
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  &buf,
	}

	// 发送查询请求
	res, err := req.Do(context.Background(), connectionEs)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 处理响应
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
