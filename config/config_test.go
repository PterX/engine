package config

import (
	"fmt"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
)

func TestConfig(t *testing.T) {
	config, found := LoadConfig()
	fmt.Println(found)
	fmt.Println(config)
	strategyCode := 82
	v := GetStrategyParameterByCode(strategyCode)
	fmt.Println(v)
}

func TestBlocks(t *testing.T) {
	sectorCode := "sh880884"
	blk := securities.GetBlockInfo(sectorCode)
	fmt.Println(len(blk.ConstituentStocks))
}
