package models

import (
	"encoding/json"
)

func ConvertToString(seat [][]int) (string, error) {
	// 序列化为 JSON 字符串
	seatJSON, err := json.Marshal(seat)
	if err != nil {
		return "", err
	}

	// 将 JSON 字符串转换为字符串
	seatString := string(seatJSON)
	return seatString, nil
}

func ConvertTo2DIntSlice(str string) ([][]int, error) {
	// 将字符串转换为字节数组
	bytes := []byte(str)

	// 反序列化为 [][]int 类型
	var result [][]int
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
