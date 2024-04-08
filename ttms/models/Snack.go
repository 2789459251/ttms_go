package models

import (
	utils "TTMS_go/ttms/util"
	"fmt"
	"gorm.io/gorm"
	"log"
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
			str += " AND name LIKE '%" + c + "%'"
		}
	}
	utils.DB.Where(str).Find(&snacks)
	return
}
func Insertsnack(snack Snack) {
	utils.DB.Create(&snack)
}
func QuerysnackByid(id int) (s Snack) {
	utils.DB.Where("id = ?", id).First(&s)
	return
}
func (s Snack) Refleshsnack() (err error) {
	err = utils.DB.Updates(s).Error
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
	s := QuerysnackByid(id)
	return utils.DB.Where("id = ?", id).Delete(s).Error
}
func DeleteSnackByNamekey(nameKey string) error {
	snacks := SearchSnack(nameKey)
	return utils.DB.Delete(snacks).Error
}
func (s *Snack) RefreshSnack() {

}
