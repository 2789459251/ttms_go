package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type MovieInfo struct {
	Info string `json:"info"`
}

func (movie MovieInfo) Index() string {
	return "movie_index"
}
func (movie MovieInfo) Mapping() string {
	return `{
		"mapping":{
			"properties":{
				"info":{
					"type":"text"
				}
			}
		}
	}`
}
func GetDocumentById(client *elastic.Client, index, id string) (*MovieInfo, error) {
	// 使用GetService根据ID获取文档
	getResult, err := client.Get().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		fmt.Printf("获取文档错误：%s\n", err)
		return nil, err
	}
	if getResult == nil {
		fmt.Println("未找到文档")
		return nil, fmt.Errorf("未找到文档")
	}

	// 将_source字段解析到Movie结构体
	var movie MovieInfo
	if err := json.Unmarshal(getResult.Source, &movie); err != nil {
		fmt.Printf("解析文档错误：%s\n", err)
		return nil, err
	}

	return &movie, nil
}
