package model

type Statement struct {
	ForwardNum int `json:"forward_num"`//转发数 
	BuyPeopleNum int64 `json:"buy_people_num"`//购买人数
	BuyNum int `json:"buy_num"`//已购买份数
	SendCoin int `json:"send_coin"`//已送出的币值
	ClickNum int `json:"click_num"` //点击数
}
