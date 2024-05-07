package service

import (
	"TTMS_go/ttms/models"
	utils "TTMS_go/ttms/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddTheatre(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	column, _ := strconv.Atoi(c.Request.FormValue("columns"))
	row, _ := strconv.Atoi(c.Request.FormValue("rows"))
	t := models.Theatre{
		Name:  c.Request.FormValue("name"),
		N:     column,
		M:     row,
		Inuse: 0,
	}
	t.Num = t.N * t.M
	if t.Num <= 0 || t.Num > 50 {
		utils.RespFail(c.Writer, "无效添加操作，请检查输入是否合法！")
		return
	}
	t.Info = c.Request.FormValue("info")
	if t.Num <= 25 {
		t.Info += "全景声MAX激光厅"
	} else {
		t.Info += "RGB基色激光厅"
	}
	Seat := make([][]int, 0)
	for i := 0; i < t.M; i++ {
		n := make([]int, t.N)
		Seat = append(Seat, n)
	}
	//seat, _ := json.Marshal(Seat)
	t.Seat, _ = models.ConvertToString(Seat)
	models.CreateTheatre(t)
	utils.RespOk(c.Writer, t, "添加放映厅成功")
	return
}
func RemoveTheatre(c *gin.Context) {
	if !isLimited(c) {
		return
	}
	id := c.Request.FormValue("id")
	err := models.DeleteTheatre(id)
	if err != nil {
		utils.RespFail(c.Writer, "删除失败:"+err.Error())
		return
	}
	utils.RespOk(c.Writer, "", "删除成功")
}
func GetAllTheatre(c *gin.Context) {
	t := models.FindAllTheatre()
	if len(t) == 0 {
		utils.RespFail(c.Writer, "尚未添加放映厅")
		return
	}
	utils.RespOk(c.Writer, t, "返回所有放映厅")
}
