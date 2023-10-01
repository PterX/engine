package features

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/features/base"
)

type DataKLine struct {
	DataCache
}

func (k *DataKLine) Kind() FeatureKind {
	return FeatureBaseKLine
}

func (k *DataKLine) Name() string {
	return mapFeatures[k.Kind()].Name
}

func (k *DataKLine) Key() string {
	return mapFeatures[k.Kind()].Key
}

func (k *DataKLine) Filename(date, code string) string {
	k.filename = cache.KLineFilename(code)
	return k.filename
}

func (k *DataKLine) Update(cacheDate, featureDate string) {
	base.UpdateAllBasicKLine(k.Code)

}

func (k *DataKLine) Repair(cacheDate, featureDate string) {
	//TODO implement me
	panic("implement me")
}

func (k *DataKLine) Increase(snapshot Snapshot) {
	//TODO K线增量更新数据的条件是缓存的数据最晚的日期是上一个交易日, 否则会缺失缓存数据中最后1条数据和当日之间的数据, 所以只能按照K线的更新方法, 不适合用快照更新
	// 第一步: 取出最后一条数据的记录
	// 第二步: 检查时间的有效性
	// 第三步: 用快照组织K线结构
	// 第四步: 如果不符合快照更新, 则忽略
	panic("implement me")
}

func (k *DataKLine) Clone(date, code string) DataSet {
	var dest = DataKLine{DataCache{Date: date, Code: code}}
	return &dest
}
