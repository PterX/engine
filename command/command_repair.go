package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	cmder "github.com/spf13/cobra"
	"strings"
	"time"
)

// CmdRepair 补登历史数据
var CmdRepair = &cmder.Command{
	Use:     "repair",
	Example: Application + " repair --all",
	//Args:    args.MinimumNArgs(0),
	Args: func(cmd *cmder.Command, args []string) error {
		return nil
	},
	Short: "回补股市数据",
	Long:  `回补股市数据`,
	Run: func(cmd *cmder.Command, args []string) {
		beginDate := trading.FixTradeDate(flagStartDate.Value)
		endDate := cache.DefaultCanReadDate()
		if len(flagEndDate.Value) > 0 {
			endDate = trading.FixTradeDate(flagEndDate.Value)
		}
		dates := trading.TradeRange(beginDate, endDate)
		count := len(dates)
		fmt.Printf("修复数据: %s => %s"+strings.Repeat("\r\n", 2), dates[0], dates[count-1])
		if flagAll.Value {
			handleRepairAll(dates)
		} else if flagDataSet.Value {
			handleRepairDataSet(dates)
		} else if flagHistory.Value {
			handleRepairFeatures(dates)
		}
	},
}

func init() {
	commandInit(CmdRepair, &flagAll)
	commandInit(CmdRepair, &flagDataSet)
	commandInit(CmdRepair, &flagHistory)
	commandInit(CmdRepair, &flagStartDate)
	commandInit(CmdRepair, &flagEndDate)
}

func handleRepairAll(dates []string) {
	moduleName := "补登全部历史数据"
	count := len(dates)
	fmt.Println()
	fmt.Println()
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.RepairAllFeature(&barIndex, cacheDate, featureDate)
		_ = cacheDate
		_ = featureDate
		bar.Add(1)
	}
	logger.Info("任务执行完毕.", time.Now())
	fmt.Println()
}

func handleRepairDataSet(dates []string) {
	fmt.Println()
	moduleName := "补登数据集合"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		//storages.RepairBaseData(&barIndex, cacheDate, featureDate)
		storages.PluginsRepairBase(&barIndex, cacheDate, featureDate)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

func handleRepairFeatures(dates []string) {
	moduleName := "补登特征数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		bar.Add(1)
		cacheDate, featureDate := cache.CorrectDate(date)
		//storages.Repair(cacheDate, featureDate)
		storages.PluginsRepairFeatures(&barIndex, cacheDate, featureDate)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}
