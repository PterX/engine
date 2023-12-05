package models

import (
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/num"
)

// QuoteSnapshot 即时行情快照
type QuoteSnapshot struct {
	Date                  string               // 交易日期
	ServerTime            string               // 时间
	SecurityCode          string               // 证券代码
	ExchangeState         quotes.ExchangeState // 交易状态
	State                 quotes.TradeState    // 上市公司状态
	Market                uint8                // 市场ID
	Code                  string               `name:"证券代码"`  // 代码
	Name                  string               `name:"证券名称"`  // 证券名称
	Active                uint16               `name:"活跃度"`   // 活跃度
	LastClose             float64              `name:"昨收"`    // 昨收
	Open                  float64              `name:"开盘价"`   // 开盘
	OpeningChangeRate     float64              `name:"开盘涨幅%"` // 开盘
	Price                 float64              `name:"现价"`    // 现价
	ChangeRate            float64              `name:"涨跌幅%"`  // 涨跌幅
	PremiumRate           float64              `name:"溢价率%"`  // 集合竞价买入溢价, hedgeRatio
	High                  float64              // 最高
	Low                   float64              // 最低
	Vol                   int                  // 总量
	CurVol                int                  // 现量
	Amount                float64              // 总金额
	SVol                  int                  // 内盘
	BVol                  int                  // 外盘
	IndexOpenAmount       int                  // 指数-集合竞价成交金额=开盘成交金额
	StockOpenAmount       int                  // 个股-集合竞价成交金额=开盘成交金额
	OpenVolume            int                  `name:"开盘量"` // 集合竞价-开盘量, 单位是股
	CloseVolume           int                  `name:"收盘量"` /// 集合竞价-收盘量, 单位是股
	IndexUp               int                  // 指数有效-上涨数
	IndexUpLimit          int                  // 指数有效-涨停数
	IndexDown             int                  // 指数有效-下跌数
	IndexDownLimit        int                  // 指数有效-跌停数
	OpenBiddingDirection  int                  `name:"开盘竞价" dataframe:"开盘竞价"` // 竞价方向, 交易当日集合竞价开盘时更新
	OpenVolumeDirection   int                  `name:"开盘竞量" dataframe:"开盘竞量"` // 委托量差, 交易当日集合竞价开盘时更新
	CloseBiddingDirection int                  `name:"收盘竞价" dataframe:"收盘竞价"` // 竞价方向, 交易当日集合竞价收盘时更新
	CloseVolumeDirection  int                  `name:"收盘竞量" dataframe:"收盘竞量"` // 委托量差, 交易当日集合竞价收盘时更新
	Rate                  float64              // 涨速
	TopNo                 int                  // 板块排名
	TopCode               string               // 领涨个股
	TopName               string               // 领涨个股名称
	TopRate               float64              // 领涨个股涨幅
	ZhanTing              int                  // 涨停数
	Ling                  int                  // 平盘数
	Count                 int                  // 总数
	Capital               float64              `name:"流通盘"`    // 流通盘
	FreeCapital           float64              `name:"自由流通股本"` // 自由流通股本
	OpenTurnZ             float64              `name:"开盘换手Z%"` // 开盘换手
	OpenQuantityRatio     float64              `name:"开盘量比"`
	QuantityRatio         float64              `name:"量比"`
	ChangePower           float64              `name:"涨跌力度"` // 开盘金额除以开盘涨幅
	AverageBiddingVolume  int                  `name:"委托均量"` // 委托均量
}

//func (q QuoteSnapshot) QuantityRelativeRatio() {
//	return q.OpenQuantityRatio = float64(q.OpenVolume) / history.GetMV5()
//}

func QuoteSnapshotFromProtocol(v quotes.Snapshot) QuoteSnapshot {
	snapshot := QuoteSnapshot{}
	_ = api.Copy(&snapshot, &v)
	securityCode := proto.GetSecurityCode(v.Market, v.Code)
	snapshot.Code = securityCode
	snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
	snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
	snapshot.PremiumRate = num.NetChangeRate(snapshot.Open, snapshot.Price)
	snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
	// 涨跌力度
	snapshot.ChangePower = float64(snapshot.OpenVolume) / snapshot.OpeningChangeRate
	snapshot.AverageBiddingVolume = v.AverageBiddingVolume()

	// 补全F10相关
	f10 := smart.GetL5F10(securityCode)
	if f10 != nil {
		snapshot.Name = f10.SecurityName
		snapshot.Capital = f10.Capital
		snapshot.FreeCapital = f10.FreeCapital
		snapshot.OpenTurnZ = f10.TurnZ(snapshot.OpenVolume)
	}
	// 补全扩展相关
	history := smart.GetL5History(securityCode)
	if history != nil && history.MV5 > 0 {
		lastMinuteVolume := history.GetMV5()
		snapshot.OpenQuantityRatio = float64(snapshot.OpenVolume) / lastMinuteVolume
		minuteVolume := float64(snapshot.Vol) / float64(trading.Minutes(snapshot.Date))
		snapshot.QuantityRatio = minuteVolume / lastMinuteVolume
	}
	return snapshot
}

// BatchSnapShot 批量获取即时行情数据快照
func BatchSnapShot(codes []string) []QuoteSnapshot {
	tdxApi := gotdx.GetTdxApi()
	list := []QuoteSnapshot{}
	var err error
	var hq []quotes.Snapshot
	retryTimes := 0
	for retryTimes < quotes.DefaultRetryTimes {
		hq, err = tdxApi.GetSnapshot(codes)
		if err == nil && hq != nil {
			break
		}
		retryTimes++
	}

	if err != nil {
		logger.Errorf("获取即时行情数据失败", err)
		return list
	}

	for _, v := range hq {
		if v.State != quotes.TDX_SECURITY_TRADE_STATE_NORMAL {
			// 非正常交易的记录忽略掉
			continue
		}
		//snapshot := QuoteSnapshot{}
		//_ = api.Copy(&snapshot, &v)
		//securityCode := proto.GetSecurityCode(v.Market, v.Code)
		//snapshot.Code = securityCode
		//snapshot.OpeningChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Open)
		//snapshot.ChangeRate = num.NetChangeRate(snapshot.LastClose, snapshot.Price)
		//snapshot.PremiumRate = num.NetChangeRate(snapshot.Open, snapshot.Price)
		//snapshot.OpenBiddingDirection, snapshot.OpenVolumeDirection = v.CheckDirection()
		//// 涨跌力度
		//snapshot.ChangePower = float64(snapshot.OpenVolume) / snapshot.OpeningChangeRate
		//snapshot.AverageBiddingVolume = v.AverageBiddingVolume()
		//
		//// 补全F10相关
		//f10 := smart.GetL5F10(securityCode)
		//if f10 != nil {
		//	snapshot.Name = f10.SecurityName
		//	snapshot.Capital = f10.Capital
		//	snapshot.FreeCapital = f10.FreeCapital
		//	snapshot.OpenTurnZ = f10.TurnZ(snapshot.OpenVolume)
		//}
		//// 补全扩展相关
		//history := smart.GetL5History(securityCode)
		//if history != nil && history.MV5 > 0 {
		//	snapshot.OpenQuantityRatio = float64(snapshot.OpenVolume) / history.GetMV5()
		//}
		snapshot := QuoteSnapshotFromProtocol(v)
		list = append(list, snapshot)
	}
	return list
}

// GetTick 获取一只股票的tick数据
//
//	该函数用于测试
func GetTick(securityCode string) *QuoteSnapshot {
	list := BatchSnapShot([]string{securityCode})
	if len(list) > 0 {
		v := list[0]
		return &v
	}
	return nil
}
