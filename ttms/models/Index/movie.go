package Index

import (
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
)

func CreateIndex(index string, mapping string) {
	if Isexist(index) {
		fmt.Println("已存在", index, "索引，执行删除操作！")
		DeleteIndex(index)
	}
	createIndex, err := utils.ES.CreateIndex(index).BodyString(mapping).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(createIndex, "创建索引成功")
}

func Isexist(index string) bool {
	exist, err := utils.ES.IndexExists(index).Do(context.Background())
	if err != nil {
		fmt.Println("判断索引是否存在的方法出错：", err)
	}
	return exist
}

func DeleteIndex(index string) {
	deleteindex, err := utils.ES.DeleteIndex(index).Do(context.Background())
	if err != nil {
		fmt.Println("删除索引：", err.Error())
		return
	}
	fmt.Println(deleteindex, "删除索引成功")
}
