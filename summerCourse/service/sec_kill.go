package service

import (
	"summerCourse/model"
	"log"
	"sync"
	"time"
)

type User struct {
	UserId string
	GoodsId  uint
}

//订单管道
var OrderChan = make(chan User, 1024)

var ItemMap = make(map[uint]*Item)

type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int
	sellCh    chan int
	done      chan struct{}
	Lock      sync.Mutex
}

// TODO 写一个定时任务，每天定时从数据库加载数据到Map
func initMap(){

	//获取当前时间
	nowTime :=time.Now()
	//获取第二天的时间
	nextTime := nowTime.Add(time.Hour*24)
	//定时 每日0点 从数据库加载数据
	renewTime :=time.Date(nextTime.Year(),nextTime.Month(),nextTime.Day(),0,0,0,0,nextTime.Location())
	d := renewTime.Sub(nowTime)
	timer :=time.NewTimer(d)
	<-timer.C
	res,err := model.SelectGoods()
	if err != nil {
		log.Println(err)
		return
	}
	for _,v :=range res {
		mod := &Item{
			ID:        v.ID,
			Name:      v.Name,
			Total:     v.Num,
			Left:      v.Num,
			IsSoldOut: false,
			leftCh:    nil,
			sellCh:    nil,
			done:      nil,
			Lock:      sync.Mutex{},
		}
		ItemMap[v.ID] = mod
	}



}


//func initMap() {
//	item := &Item{
//		ID:        1,
//		Name:      "测试",
//		Total:     100,
//		Left:      100,
//		IsSoldOut: false,
//		leftCh:    make(chan int),
//		sellCh:    make(chan int),
//	}
//	ItemMap[item.ID] = item
//}


//获取秒杀商品信息
func getItem(itemId uint) *Item{
	return ItemMap[itemId]
}


//监听订单管道，处理秒杀订单
func order() {
	for {
		user := <- OrderChan
		item := getItem(user.GoodsId)
		item.SecKilling(user.UserId)
	}
}


//秒杀
func (item *Item) SecKilling(userId string) {

	//加锁
	item.Lock.Lock()
	defer item.Lock.Unlock()
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }
	if item.IsSoldOut {
		return
	}
	item.BuyGoods(1)

	MakeOrder(userId, item.ID,1)


}

// 定时下架
func (item *Item) OffShelve() {
	beginTime := time.Now()
	// 获取第二天时间
	//nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	//offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	offShelveTime := beginTime.Add(time.Second*5)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	<-timer.C
	delete(ItemMap, item.ID)
	close(item.done)

}
// 出售商品
func (item *Item) SalesGoods() {
	for {
		select {
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {  //监听商品是否售完
				item.IsSoldOut = true
			}

		case item.leftCh <- item.Left:  //监听商品库存
		case <-item.Done():
			log.Println("我自闭了")
			return
		}
	}
}

func (item *Item) Done() <-chan struct{} {
	if item.done == nil {
		item.done = make(chan struct{})
	}
	d := item.done
	return d
}

//监听秒杀时商品的数据
func (item *Item) Monitor() {
	go item.SalesGoods()
}

// 获取剩余库存
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num
}


//开启秒杀服务
func InitService() {
	//从数据库中加载秒杀商品数据
	initMap()

	for _,item := range ItemMap{
		item.Monitor()
		go item.OffShelve()
	}
	for i := 0; i < 10; i++ {
		go order()
	}
}
