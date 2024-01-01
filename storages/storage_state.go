package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/trader"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/gox/logger"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	// 状态文件扩展名
	orderStateFileExtension = ".done"
)

// Touch 创建一个空文件
func Touch(filename string) error {
	_ = api.CheckFilepath(filename, true)
	return os.WriteFile(filename, nil, 0644)
}

// 获取状态机路径
func state_filepath(state_date string) string {
	flagPath := filepath.Join(cache.GetQmtCachePath(), "var", state_date)
	return flagPath
}

// 获取状态文件前缀
func state_prefix(state_date, qmtStrategyName string, direction trader.Direction) string {
	qmtStrategyName = strings.ToLower(qmtStrategyName)
	prefix := fmt.Sprintf("%s-%s-%s-%s", state_date, traderConfig.AccountId, qmtStrategyName, direction.Flag())
	return prefix
}

// 订单状态文件前缀
func order_state_prefix(state_date string, model models.Strategy, direction trader.Direction) string {
	qmtStrategyName := models.QmtStrategyName(model)
	prefix := state_prefix(state_date, qmtStrategyName, direction)
	return prefix
}

// 获得订单标识文件名
func order_state_filename(date string, model models.Strategy, code string, direction trader.Direction) string {
	state_date := trading.FixTradeDate(date, cache.CACHE_DATE)
	orderFlagPath := state_filepath(state_date)
	prefix := order_state_prefix(date, model, direction)
	securityCode := proto.CorrectSecurityCode(code)
	filename := fmt.Sprintf("%s-%s.done", prefix, securityCode)
	state_filename := path.Join(orderFlagPath, filename)
	return state_filename
}

// CheckOrderState 检查订单执行状态
func CheckOrderState(date string, model models.Strategy, code string, direction trader.Direction) bool {
	filename := order_state_filename(date, model, code, direction)
	return api.FileExist(filename)
}

// PushOrderState 推送订单完成状态
func PushOrderState(date string, model models.Strategy, code string, direction trader.Direction) error {
	filename := order_state_filename(date, model, code, direction)
	return Touch(filename)
}

// CountStrategyOrders 统计策略订单数
func CountStrategyOrders(date string, model models.Strategy, direction trader.Direction) int {
	stateDate := trading.FixTradeDate(date, cache.CACHE_DATE)
	orderFlagPath := state_filepath(stateDate)
	filenamePrefix := order_state_prefix(stateDate, model, direction)
	pattern := filepath.Join(orderFlagPath, filenamePrefix+"-*"+orderStateFileExtension)
	files, err := filepath.Glob(pattern)
	if err != nil {
		logger.Error(err)
		return 0
	}
	return len(files)
}

// FetchListForFirstPurchase 获取指定日期交易的个股列表
func FetchListForFirstPurchase(date, qmtStrategyName string, direction trader.Direction) []string {
	stateDate := trading.FixTradeDate(date, cache.CACHE_DATE)
	orderFlagPath := state_filepath(stateDate)
	filenamePrefix := state_prefix(stateDate, qmtStrategyName, direction)
	var list []string
	prefix := filepath.Join(orderFlagPath, filenamePrefix+"-")
	pattern := prefix + "*" + orderStateFileExtension
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		logger.Error(err)
		return list
	}
	for _, filename := range files {
		after, found := strings.CutPrefix(filename, prefix)
		if !found {
			continue
		}
		before, found := strings.CutSuffix(after, orderStateFileExtension)
		if !found {
			continue
		}
		list = append(list, before)
	}
	return list
}
