package main

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

const (
	AccessKey = "y_XTiaH5dywx_R-J-twejWCQRXvBd5jI54YT9ihT"
	SerectKey = "2g0S7zGWZ_zca0BVwYTeugUoZJepYLsYjd5bKGir"
	Bucket    = "taoniuma"
	ImgUrl    = "http://s9isqyrv9.hn-bkt.clouddn.com/"
)

// 封装上传图片到七牛云然后返回状态和图片的
func UploadToQiNiu(file multipart.File, fileSize int64) (string, error) {
	putPlicy := storage.PutPolicy{Scope: Bucket}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + ret.Key
	return url, nil
}
