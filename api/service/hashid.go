package service

import (
	"github.com/beego/beego"
	"github.com/speps/go-hashids/v2"
)

var hash *hashids.HashID

func init() {
	hd := hashids.NewData()
	hd.Salt = beego.AppConfig.String("APP_KEY")
	hd.MinLength = 5
	hash, _ = hashids.NewWithData(hd)
}

func HashIdEncode(id int) (string, error) {
	return hash.Encode([]int{id})
}

func HashIdDecode(value string) ([]int, error) {
	return hash.DecodeWithError(value)
}
