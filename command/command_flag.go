package command

import (
	flags "github.com/spf13/cobra"
	"unsafe"
)

var (
	flagHistory   = cmdFlag[bool]{Name: "history", Value: false, Usage: "历史特征数据"}
	flagStartDate = cmdFlag[string]{Name: "start", Value: "", Usage: "开始日期"}
	flagEndDate   = cmdFlag[string]{Name: "end", Value: "", Usage: "结束日期"}
)

type cmdFlag[T ~int | ~bool | ~string] struct {
	Name  string
	Usage string
	Value T
}

func (cf *cmdFlag[T]) init(cmd *flags.Command) {
	switch v := any(cf.Value).(type) {
	case bool:
		cmd.Flags().BoolVar((*bool)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case int:
		cmd.Flags().IntVar((*int)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case string:
		cmd.Flags().StringVar((*string)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	}
}
