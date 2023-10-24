package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pandas"
	"gitee.com/quant1x/pandas/stat"
)

// BasicKLine 基础日K线
func BasicKLine(securityCode string) pandas.DataFrame {
	securityCode = proto.CorrectSecurityCode(securityCode)
	filename := cache.KLineFilename(securityCode)
	df := pandas.ReadCSV(filename)
	return df
}

// KLine 加载日K线宽表
func KLine(securityCode string) pandas.DataFrame {
	securityCode = proto.CorrectSecurityCode(securityCode)
	df := GetCacheKLine(securityCode, false)
	return df
}

// KLineToWeekly 日线转周线
//
//	deprecated: 不推荐使用
func KLineToWeekly(kline pandas.DataFrame) pandas.DataFrame {
	// 周线
	var df pandas.DataFrame
	//date,open,close,high,low,volume,amount,up,down
	var wdate string
	var o, c, h, l, v, a stat.DType
	var bv, sv, ba, sa stat.DType
	var prevClose stat.DType
	for i := 0; i < kline.Nrow(); i++ {
		m := kline.IndexOf(i)
		//date,open,close,high,low,volume,amount,up,down
		// 周线日期以最后一天的日期为准
		_date, ok := m["date"].(string)
		if ok {
			wdate = _date
		}
		// 周线开盘价以第一天OPEN为准
		_open, ok := m["open"].(stat.DType)
		if ok && o == stat.DType(0) {
			o = _open
		}
		// 周线的收盘价以本周最后一个交易日的CLOSE为准
		_close, ok := m["close"].(stat.DType)
		if ok {
			c = _close
		}
		// 涨幅
		zf := (c/prevClose - 1.00) * 100.00
		_high, ok := m["high"].(stat.DType)
		if ok && h == stat.DType(0) {
			h = _high
		}
		if h < _high {
			h = _high
		}
		_low, ok := m["low"].(stat.DType)
		if ok && l == stat.DType(0) {
			l = _low
		}
		if l > _low {
			l = _low
		}
		_vol, ok := m["volume"]
		if ok {
			v += stat.Any2DType(_vol)
		}
		_amount, ok := m["amount"].(stat.DType)
		if ok {
			a += _amount
		}
		_bv, ok := m["bv"]
		if ok {
			bv += stat.Any2DType(_bv)
		}
		_sv, ok := m["sv"]
		if ok {
			sv += stat.Any2DType(_sv)
		}
		_ba, ok := m["ba"]
		if ok {
			ba += stat.Any2DType(_ba)
		}
		_sa, ok := m["sa"]
		if ok {
			sa += stat.Any2DType(_sa)
		}
		dt, _ := api.ParseTime(wdate)
		w := int(dt.Weekday())
		last := false
		today := trading.IndexToday()
		if wdate == today {
			last = true
		}
		// 如果是周五
		if !last && w == 5 {
			last = true
		}
		if !last {
			nextDate := trading.NextTradeDate(wdate)
			ndt, _ := api.ParseTime(nextDate)
			nw := int(ndt.Weekday())
			if nw < w || api.DifferDays(ndt, dt) >= 7 {
				last = true
			}
		}
		if last {
			df0 := pandas.NewDataFrame(
				pandas.NewSeries(stat.SERIES_TYPE_STRING, "date", wdate),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "open", o),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "close", c),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "high", h),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "low", l),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "volume", v),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "amount", a),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "bv", bv),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sv", sv),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "ba", ba),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "sa", sa),
				pandas.NewSeries(stat.SERIES_TYPE_DTYPE, "zf", zf),
			)
			df = df.Concat(df0)
			wdate = ""
			prevClose = c
			o = stat.DType(0)
			c = stat.DType(0)
			h = stat.DType(0)
			l = stat.DType(0)
			v = stat.DType(0)
			a = stat.DType(0)
		}
	}
	return df
}

// MovingAverage 移动平均线(MA)
type MovingAverage struct {
	MA5  float64 // 5日均线
	MA10 float64 // 10日均线
	MA20 float64 // 20日均线
}

// 由日线计算日线以上级别的K线
func periodKLine(checkPeriod func(date ...string) (s, e string), securityCode string, cacheKLine ...[]base.KLine) (df pandas.DataFrame) {
	baseKLines := []base.KLine{}
	if len(cacheKLine) > 0 {
		baseKLines = cacheKLine[0]
	} else {
		baseKLines = base.LoadBasicKline(securityCode)
	}
	if len(baseKLines) == 0 {
		return
	}

	var klines []base.KLine
	var kline base.KLine
	length := len(baseKLines)
	for i, v := range baseKLines {
		// 确定时间, 周线的日期是本周内最后一个交易日
		if len(kline.Date) == 0 {
			// 重新计算周线, 先确认周线范围
			ws, we := checkPeriod(v.Date)
			_ = ws
			//dates := trading.TradeRange(ws, we)
			//days := len(dates)
			//days = 7
			//if days > 0 {
			periodLastDate := trading.FixTradeDate(we)
			//if periodLastDate == "2023-07-30" {
			//	fmt.Println(1)
			//}
			offset := i
			for {
				destDate := baseKLines[offset].Date
				if destDate < periodLastDate {
					offset++
				}
				if offset >= length {
					periodLastDate = destDate
					break
				} else if baseKLines[offset].Date == periodLastDate {
					periodLastDate = baseKLines[offset].Date
					break
				} else if baseKLines[offset].Date > periodLastDate {
					periodLastDate = destDate
					break
				}
			}
			kline.Date = periodLastDate
			//} else {
			//	return
			//}
		}
		// 周线开盘价以第一天OPEN为准
		if kline.Open == stat.DType(0) {
			kline.Open = v.Open
		}
		// 周线的收盘价以本周最后一个交易日的CLOSE为准
		kline.Close = v.Close
		if kline.High == stat.DType(0) {
			kline.High = v.High
		} else if kline.High < v.High {
			kline.High = v.High
		}
		if kline.Low == stat.DType(0) {
			kline.Low = v.Low
		} else if kline.Low > v.Low {
			kline.Low = v.Low
		}
		kline.Volume += v.Volume
		kline.Amount += v.Amount

		// 切换下一周
		if kline.Date == v.Date || i+1 >= len(baseKLines) {
			kline.Date = v.Date
			klines = append(klines, kline)
			kline = base.KLine{}
		}
	}
	df = pandas.LoadStructs(klines)
	//fmt.Println(df)
	return
}

// WeeklyKLine 周线
func WeeklyKLine(securityCode string, cacheKLine ...[]base.KLine) (df pandas.DataFrame) {
	return periodKLine(api.GetWeekDay, securityCode, cacheKLine...)
}

// MonthlyKLine 月K线
func MonthlyKLine(securityCode string, cacheKLine ...[]base.KLine) (df pandas.DataFrame) {
	return periodKLine(api.GetMonthDay, securityCode, cacheKLine...)
}
