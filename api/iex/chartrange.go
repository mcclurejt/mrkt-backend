package iex

type ChartRange string

func (c ChartRange) String() string {
	return EnumToString(c)
}

const (
	ChartRangeMax     ChartRange = "max"
	ChartRange5y      ChartRange = "5y"
	ChartRange2y      ChartRange = "2y"
	ChartRange1y      ChartRange = "1y"
	ChartRangeYTD     ChartRange = "ytd"
	ChartRange6m      ChartRange = "6m"
	ChartRange3m      ChartRange = "3m"
	ChartRange1m      ChartRange = "1m"
	ChartRange1mm     ChartRange = "1mm"
	ChartRange5d      ChartRange = "5d"
	ChartRange5dm     ChartRange = "5dm"
	ChartRangeDynamic ChartRange = "dynamic"
)

var chartValidRanges = map[string]bool{
	ChartRangeMax.String():     true,
	ChartRange5y.String():      true,
	ChartRange2y.String():      true,
	ChartRange1y.String():      true,
	ChartRangeYTD.String():     true,
	ChartRange6m.String():      true,
	ChartRange3m.String():      true,
	ChartRange1m.String():      true,
	ChartRange1mm.String():     true,
	ChartRange5d.String():      true,
	ChartRange5dm.String():     true,
	ChartRangeDynamic.String(): true,
}
