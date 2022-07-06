package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// UnitTypeKey 단위 정의 키
type UnitTypeKey string

// UnitType ...
type UnitType struct {
	Units   []string
	Divisor float64
}

const (
	Count                = UnitTypeKey("Count")
	Percentage           = UnitTypeKey("Percentage")
	Core                 = UnitTypeKey("Core")
	Numeric              = UnitTypeKey("Numeric")
	DecimalBytes         = UnitTypeKey("DecimalBytes")
	DecimalBytesWithoutB = UnitTypeKey("DecimalBytesWithoutB")
	BinaryBytes          = UnitTypeKey("BinaryBytes")
	BinaryBytesWithoutB  = UnitTypeKey("BinaryBytesWithoutB")
	SI                   = UnitTypeKey("SI")
	DecimalBytesPerSec   = UnitTypeKey("DecimalBytesPerSec")
	PacketsPerSec        = UnitTypeKey("PacketsPerSec")
	Seconds              = UnitTypeKey("Seconds")
)

var (
	UnitTypes = map[UnitTypeKey]UnitType{
		Count: {
			Units:   []string{""},
			Divisor: 1,
		},
		Percentage: {
			Units:   []string{"%"},
			Divisor: 1,
		},
		Core: {
			Units:   []string{"Core"},
			Divisor: 1,
		},
		Numeric: {
			Units:   []string{"", "k", "m", "b"},
			Divisor: 1000,
		},
		DecimalBytes: {
			Units:   []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"},
			Divisor: 1000,
		},
		DecimalBytesWithoutB: {
			Units:   []string{"", "k", "M", "G", "T", "P", "E"},
			Divisor: 1000,
		},
		BinaryBytes: {
			Units:   []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"},
			Divisor: 1024,
		},
		BinaryBytesWithoutB: {
			Units:   []string{"i", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"},
			Divisor: 1024,
		},
		SI: {
			Units:   []string{"", "k", "M", "G", "T", "P", "E"},
			Divisor: 1000,
		},
		DecimalBytesPerSec: {
			Units:   []string{"Bps", "KBps", "MBps", "GBps", "TBps", "PBps", "EBps"},
			Divisor: 1000,
		},
		PacketsPerSec: {
			Units:   []string{"pps", "kpps"},
			Divisor: 1000,
		},
		Seconds: {
			Units:   []string{"ns", "μs", "ms", "s"},
			Divisor: 1000,
		},
	}
)

type HumanizeOptions struct {
	Precision     uint
	InitialUnit   string
	PreferredUnit string
}

type HumanizeValue struct {
	Unit  string
	Value float64
}

type convertedValue struct {
	value interface{}
	unit  string
}

// shift https://go.dev/play/p/WAEIzRgSNB
func shift(pToSlice *[]string) string {
	sValue := (*pToSlice)[0]
	*pToSlice = (*pToSlice)[1:len(*pToSlice)]
	return sValue
}

func indexOf(arr []string, val string) int {
	for pos, v := range arr {
		if v == val {
			return pos
		}
	}
	return -1
}

// roundFloat https://gosamples.dev/round-float/
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func convertBaseValueToUnits(value float64, unitArray []string, divisor float64, initialUnit string, preferredUnit string) *convertedValue {
	var sliceIndex = 0
	if initialUnit != "" {
		sliceIndex = indexOf(unitArray, initialUnit)
	}
	if sliceIndex == -1 {
		panic("입력 단위를 확인 해주세요.")
	}
	units := unitArray[sliceIndex:]

	unitIndex := indexOf(units, preferredUnit)
	if unitIndex != -1 {
		return &convertedValue{value / math.Pow(divisor, float64(unitIndex)), preferredUnit}
	}

	unit := shift(&units)
	for value >= divisor && len(units) > 0 {
		value = value / divisor
		unit = shift(&units)
	}

	return &convertedValue{value, unit}
}

func Humanize(value float64, key UnitTypeKey, options *HumanizeOptions) *HumanizeValue {
	types := UnitTypes[key]

	convertedValue := convertBaseValueToUnits(value, types.Units, types.Divisor, options.InitialUnit, options.PreferredUnit)

	result := roundFloat(convertedValue.value.(float64), options.Precision)

	if result == 0 {
		var offset int
		if value != 0 {
			for i, v := range strings.Split(strconv.FormatFloat(value, 'f', -1, 64), ".")[1] {
				if string(v) != "0" {
					offset = i + 2
					break
				}
			}
		}

		result = roundFloat(value, uint(offset))
	}

	return &HumanizeValue{
		convertedValue.unit,
		result,
	}
}

func FindMaxUnitByValues(unitTypeKey UnitTypeKey, values interface{}) string {
	var maxFloat64 float64
	if str, ok := values.(float64); ok {
		maxFloat64 = values.(float64)
	} else {
		fmt.Println("float 64 not ok", str)
	}

	if str, ok := values.([]float64); ok {
		fmt.Println("[]float 64 ok", str)
		for _, float := range values.([]float64) {
			if maxFloat64 < float {
				maxFloat64 = float
			}
		}
	} else {
		fmt.Println("[]float 64 not ok", str)
	}

	return Humanize(maxFloat64, unitTypeKey, &HumanizeOptions{Precision: 2}).Unit
}

func main() {
	//var arr = []string{"", "k", "m", "b"}

	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "", ""))
	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "k", ""))
	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "k", "b"))

	//number := 12.00
	//fmt.Println(roundFloat(number, 2))
	//
	fmt.Println(fmt.Sprintf("%v", Humanize(0.00002, DecimalBytes, &HumanizeOptions{2, "", ""})))

	value := 0.0

	fmt.Println(strconv.FormatFloat(value, 'f', -1, 64))

	var ttt = strconv.FormatFloat(value, 'f', -1, 64)

	var tt int
	if ttt != "0" {
		for i, v := range strings.Split(ttt, ".")[1] {
			if string(v) != "0" {
				tt = i
				break
			}
		}
	}

	//fmt.Println(roundFloat(0.0000000234, 10))
	//fmt.Println(roundFloat(0.234, 1))
	//fmt.Println(roundFloat(0.234, 2))
	//fmt.Println(roundFloat(0.035, 1))
	fmt.Println(strconv.FormatFloat(roundFloat(value, uint(tt+2)), 'f', -1, 64))

	//fmt.Println(FindMaxUnitByValues(BinaryBytes, []float64{311113.33, 3111.33, 300000000}))
	//fmt.Println(FindMaxUnitByValues(BinaryBytes, []float64{1000, 10000, 100000099000}))

}
