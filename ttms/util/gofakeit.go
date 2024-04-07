package utils

import (
	"github.com/brianvoe/gofakeit"
	"time"
)

func GenerateSMSCode() int {
	gofakeit.Seed(time.Now().UnixNano())
	return gofakeit.Number(100000, 999999)
}
