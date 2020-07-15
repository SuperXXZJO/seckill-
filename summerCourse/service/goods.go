package service

import (
	"summerCourse/model"
	"github.com/jinzhu/gorm"
	"log"
)

type Goods struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Num   int    `json:"num"`
}

// 添加商品
func AddGoods(goodName string,goodPrice int,goodNum int) {
	// TODO
	mod :=&model.Goods{
		Model: gorm.Model{},
		Name:  goodName,
		Price: goodPrice,
		Num:   goodNum,
	}
	if err :=model.Goods.AddGoods(mod);err !=nil{
		log.Println(err)
		return
	}

}

//查看所有商品
func SelectGoods() (goods []Goods) {
	_goods, err := model.SelectGoods()
	if err != nil {
		log.Printf("Error get goods info. Error: %s", err)
	}

	//将数据库传回来的good写入切片，并返回
	for _, v := range _goods {
		good := Goods{
			ID:    v.ID,
			Name:  v.Name,
			Price: v.Price,
			Num:   v.Num,
		}
		goods = append(goods, good)
	}
	return goods
}
