package datasets

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/gotdx/quotes"
)

//type DataKind = uint64

const (
	baseKind cache.Kind = 0
)

const (
	BaseXdxr             = cache.PluginMaskBaseData | (baseKind + 1) // 基础数据-除权除息
	BaseKLine            = cache.PluginMaskBaseData | (baseKind + 2) // 基础数据-基础K线
	BaseTransaction      = cache.PluginMaskBaseData | (baseKind + 3) // 基础数据-历史成交
	BaseMinutes          = cache.PluginMaskBaseData | (baseKind + 4) // 基础数据-分时数据
	BaseQuarterlyReports = cache.PluginMaskBaseData | (baseKind + 5) // 基础数据-季报
	BaseSafetyScore      = cache.PluginMaskBaseData | (baseKind + 6) // 基础数据-安全分
	//BaseAggregationData  cache.Kind = 1 << 63
	//BaseTest             DataKind   = 0x8000000000000000
)

// DataSet 数据层, 数据集接口 smart
//
//	数据集是基础数据, 应当遵循结构简单, 尽量减小缓存的文件数量, 加载迅速
//	检索的规则是按日期和代码进行查询
type DataSet interface {
	cache.Trait
	Update(cacheDate, featureDate string)   // 更新数据
	Repair(cacheDate, featureDate string)   // 回补数据
	Increase(snapshot quotes.Snapshot)      // 增量计算, 用快照增量计算特征
	Clone(date string, code string) DataSet // 克隆一个DataSet
}

var (
	mapDataSets = map[cache.Kind]cache.DataSummary{
		BaseXdxr:             cache.Summary(BaseXdxr, "xdxr", "除权除息", cache.DefaultDataProvider),
		BaseKLine:            cache.Summary(BaseKLine, "day", "日K线", cache.DefaultDataProvider),
		BaseTransaction:      cache.Summary(BaseTransaction, "trans", "成交数据", cache.DefaultDataProvider),
		BaseMinutes:          cache.Summary(BaseMinutes, "minutes", "分时数据", cache.DefaultDataProvider),
		BaseQuarterlyReports: cache.Summary(BaseQuarterlyReports, "reports", "季报", cache.DefaultDataProvider),
		BaseSafetyScore:      cache.Summary(BaseSafetyScore, "safetyscore", "安全分", cache.DefaultDataProvider),
	}
)

func GetDataDescript(kind cache.Kind) cache.DataSummary {
	v, ok := mapDataSets[kind]
	if !ok {
		panic("类型不存在")
	}
	return v
}
