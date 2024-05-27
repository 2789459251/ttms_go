package main

import (
	utils "TTMS_go/ttms/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/http"
)

const url_ = "http://video.cdn.zy520.online/"

func Upload(c *gin.Context) {

	str, err := upload(c.Request, c.Writer, c)
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
		return
	}
	utils.RespOk(c.Writer, str, "成功")
}

func upload(r *http.Request, w http.ResponseWriter, c *gin.Context) (string, error) {
	putPolicy := storage.PutPolicy{Scope: "zsy-ttms"}
	mac := qbox.NewMac("C5IdEtxZFLx0oS9iY8C1yrS6rLYnfzB4XV8rt4HU", "MpTfRFklA1fWfmym2YhSxebykuxmi6SoBCFOd1Y4")
	upTocken := putPolicy.UploadToken(mac)

	cfg := storage.Config{Zone: &storage.ZoneHuabei, UseHTTPS: false, UseCdnDomains: false}
	file, head, err := r.FormFile("picture")
	if err != nil {
		utils.RespFail(c.Writer, "文件读取失败："+err.Error())
	}
	fmt.Println(head.Header)
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	fmt.Println(head.Filename)
	err = formUploader.Put(context.Background(), &ret, upTocken, head.Filename, file, head.Size, &putExtra)
	fmt.Println(ret.Key)
	return url_ + ret.Key, err
}
