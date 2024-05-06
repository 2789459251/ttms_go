package docs

import (
	model2 "TTMS_go/ttms/models/model"
	utils "TTMS_go/ttms/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type model interface {
	Index() string
}

func CreateDoc(model model) (string, error) {
	IndexResponse, err := utils.ES.Index().Index(model.Index()).BodyJson(model).Do(context.Background())
	if err != nil {
		fmt.Println("创建文档错误：", err)
		return "", err
	}
	return IndexResponse.Id, nil
}
func DeleteDoc(model model, id string) {
	deleteResponse, err := utils.ES.Delete().Index(model.Index()).Id(id).Refresh("true").Do(context.Background())
	if err != nil {
		fmt.Println("创建文档错误：", err)
		return
	}
	fmt.Printf("创建文档：%#v", deleteResponse)
}
func DeleteDocs(model model, ids []string) {
	bulk := utils.ES.Bulk().Index(model.Index()).Refresh("true")

	for _, id := range ids {
		req := elastic.NewBulkDeleteRequest().Id(id)
		bulk.Add(req)
	}

	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println("批量删除文档错误：", err)
		return
	}
	fmt.Println("批量删除文档:", res.Succeeded())

}

func CreateDocs(models []model) {
	bulk := utils.ES.Bulk().Index(models[0].Index()).Refresh("true")
	for _, m := range models {
		req := elastic.NewBulkCreateRequest().Index(m.Index())
		bulk.Add(req)
	}
	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println("创建文档错误：", err)
		return
	}
	fmt.Println("批量创建文档：", res.Succeeded())
}
func FindDoc(model model) ([]*elastic.SearchHit, error) {
	//todo 注意这里可以修改
	limit := 10
	page := 1
	from := (page - 1) * limit
	query := elastic.NewBoolQuery()
	reslist, err := utils.ES.Search(model.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println("查询文档列表错误：", err)
		return nil, err
	}
	count := reslist.Hits.TotalHits.Value
	fmt.Println("查到的数量：", count)
	for _, hit := range reslist.Hits.Hits {
		fmt.Println(string(hit.Source))
	}
	return reslist.Hits.Hits, nil
}

// 精确匹配是指keyword来匹配
//
//	func FindDocExact(model model, name, text string) ([]*elastic.SearchHit, error) {
//		limit := 2
//		page := 1
//		from := (page - 1) * limit
//		query := elastic.NewMatchQuery(name, text)
//		reslist, err := utils.ES.Search(model.Index()).Query(query).From(from).Size(limit).Do(context.Background())
//		if err != nil {
//			fmt.Println("查询文档列表错误：", err)
//			return nil, err
//		}
//		count := reslist.Hits.TotalHits.Value
//		fmt.Println("查到的数量：", count)
//		return reslist.Hits.Hits, nil
//	}
func FindDocById(client *elastic.Client, index, id string) (json.RawMessage, error) {
	getResult, err := client.Get().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取文档出错：%s\n", err)
	}
	fmt.Println(string(getResult.Source))
	return getResult.Source, err
}

// todo 自己写！
func UpdateMovieDoc(model model, info model2.MovieInfo, id string) (string, error) {
	updateRes, err := utils.ES.Update().Index(model.Index()).Id(id).Doc(map[string]interface{}{
		"info": info,
	}).Do(context.Background())
	if err != nil {
		fmt.Println("更新文档错误：", err)
		return "", err
	}
	fmt.Println(updateRes.Id)
	return updateRes.Id, nil
}
