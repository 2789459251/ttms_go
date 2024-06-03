package models

import (
	utils "TTMS_go/ttms/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
)

type Snack struct {
	gorm.Model
	mu      sync.RWMutex
	Name    string
	Picture string
	Info    string
	Stock   int     //库存量
	Price   float64 //价格
}

func (snack Snack) TableName() string {
	return "snack_basic"
}
func Showsnacks() (snacks []Snack) {
	utils.DB.Find(&snacks)
	return
}
func SearchSnack(name string) (snacks []Snack) {
	str := ""
	for i, i2 := range name {
		c := string(i2)
		if i == 0 {
			str += "name LIKE '%" + c + "%'"
		} else {
			str += " OR name LIKE '%" + c + "%'"
		}
	}
	utils.DB.Where(str).Find(&snacks)
	return
}
func Insertsnack(snack Snack) {
	utils.DB.Create(&snack)
}
func QuerysnackByid(id string) (s Snack) {
	utils.DB.Where("id = ?", id).First(&s)
	return
}
func (s Snack) Refleshsnack() (err error) {
	err = utils.DB.Updates(&s).Error
	return
}
func (s *Snack) GetStock() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Stock
}
func (s *Snack) UpdateStock(Func func() (err error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if Func() != nil {
		log.Println(fmt.Sprintln("更新操作有错误，事务回滚"))
	}
}

func DeleteSnackByid(id int) error {
	id_str := strconv.Itoa(id)
	s := QuerysnackByid(id_str)
	ss := []Snack{}
	if s.Name == "" {
		return errors.New("没有id为" + id_str + "的零食!")
	}
	ss = append(ss, s)
	return utils.DB.Delete(ss).Error
}
func DeleteSnackByNamekey(nameKey string) error {
	snacks := SearchSnack(nameKey)
	if len(snacks) == 0 {
		return errors.New("没有name包含" + nameKey + "关键字的零食!")
	}
	return utils.DB.Delete(snacks).Error
}
func (s *Snack) RefreshSnack() {
	utils.DB.Model(s).Save(&s)
}

func FindSnackByIds(ids []string) []Snack {
	snacks := []Snack{}
	for _, id := range ids {
		snack := Snack{}
		utils.DB.Where("id = ?", id).First(&snack)
		snacks = append(snacks, snack)
	}
	return snacks
}

func GetsnackByid(id string) (s Snack) {
	utils.DB.Where("id = ?", id).First(&s)
	return
}
