package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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

// UnitTypeKey 단위 타입 키
type UnitTypeKey string

// UnitType 단위 타입 정의 구조체
type UnitType struct {
	Units   []string
	Divisor float64
}

// UnitTypes 단위 타입 키에 따른 단위 타입 정의 상수
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

// HumanizeOptions Humanize 함수의 세 번째 인자의 타입으로 Humanize 함수에서 사용하는 옵션을 정의
type HumanizeOptions struct {
	Precision     uint
	InitialUnit   string
	PreferredUnit string
}

type humanizeValue struct {
	Unit  string
	Value float64
}

type convertedValue struct {
	value interface{}
	unit  string
}

// shift https://go.dev/play/p/WAEIzRgSNB 인자로 받은 string 배열을 shift 하고 값을 반환하는 함수
func shift(pToSlice *[]string) string {
	sValue := (*pToSlice)[0]
	*pToSlice = (*pToSlice)[1:len(*pToSlice)]
	return sValue
}

// indexOf string 배열에 대상값 인덱스 반환
func indexOf(strings []string, val string) int {
	for idx, v := range strings {
		if v == val {
			return idx
		}
	}
	return -1
}

// RoundFloat https://gosamples.dev/round-float/ 인자로 받은 value 를 소수점 precision 자리에서 반올림하는 함수
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// convertBaseValueToUnits value의 단위를 조정하는 함수
func convertBaseValueToUnits(value float64, unitArray []string, divisor float64, initialUnit string, preferredUnit string) *convertedValue {
	var sliceIndex = 0
	var unit = ""
	if initialUnit != "" {
		sliceIndex = indexOf(unitArray, initialUnit)
	}
	if sliceIndex != -1 {
		units := unitArray[sliceIndex:]

		unitIndex := indexOf(units, preferredUnit)
		if unitIndex != -1 {
			return &convertedValue{value / math.Pow(divisor, float64(unitIndex)), preferredUnit}
		}

		unit = shift(&units)
		for value >= divisor && len(units) > 0 {
			value = value / divisor
			unit = shift(&units)
		}
	}
	return &convertedValue{value, unit}
}

// Humanize 인자로 들어온 값을 단위 타입 키에 따라 변환된 humanizeValue 를 반환하는 함수
func Humanize(value float64, unitTypeKey UnitTypeKey, options *HumanizeOptions) *humanizeValue {
	types := UnitTypes[unitTypeKey]

	convertedValue := convertBaseValueToUnits(value, types.Units, types.Divisor, options.InitialUnit, options.PreferredUnit)

	result := RoundFloat(convertedValue.value.(float64), options.Precision)

	if value != 0 && result == 0 {
		var offset int
		for i, v := range strings.Split(strconv.FormatFloat(value, 'f', -1, 64), ".")[1] {
			if string(v) != "0" {
				offset = i + 2
				break
			}
		}

		result = RoundFloat(value, uint(offset))
	}

	return &humanizeValue{
		convertedValue.unit,
		result,
	}
}

// FindMaxUnitByValues 인자로 들어온 값(배열)의 최대 단위를 반환하는 함수
func FindMaxUnitByValues(unitTypeKey UnitTypeKey, values interface{}) string {
	var maxFloat64 float64 = 0
	if _, ok := values.(float64); ok {
		maxFloat64 = values.(float64)
	}
	if _, ok := values.([]float64); ok {
		for _, float := range values.([]float64) {
			if maxFloat64 < float {
				maxFloat64 = float
			}
		}
	}

	return Humanize(maxFloat64, unitTypeKey, &HumanizeOptions{Precision: 2}).Unit
}

func main() {
	//var arr = []string{"", "k", "m", "b"}

	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "", ""))
	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "k", ""))
	//fmt.Println(convertBaseValueToUnits(1000, arr, 1000, "k", "b"))

	//number := 12.00
	//fmt.Println(RoundFloat(number, 2))
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

	//fmt.Println(RoundFloat(0.0000000234, 10))
	//fmt.Println(RoundFloat(0.234, 1))
	//fmt.Println(RoundFloat(0.234, 2))
	//fmt.Println(RoundFloat(0.035, 1))
	fmt.Println(strconv.FormatFloat(RoundFloat(value, uint(tt+2)), 'f', -1, 64))

	//fmt.Println(FindMaxUnitByValues(BinaryBytes, []float64{311113.33, 3111.33, 300000000}))
	//fmt.Println(FindMaxUnitByValues(BinaryBytes, []float64{1000, 10000, 100000099000}))

}
