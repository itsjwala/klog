package report

import (
	"fmt"
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/service"
)

type monthAggregator struct {
	y int
}

func NewMonthAggregator() Aggregator {
	return &monthAggregator{-1}
}

func (a *monthAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *monthAggregator) DateHash(date Date) Hash {
	return Hash(service.NewMonthHash(date))
}

func (a *monthAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    "). // 2020
		CellL("   ")   // Dec
}

func (a *monthAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	if date.Year() != a.y {
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Month
	table.CellR(lib.PrettyMonth(date.Month())[:3])
}
