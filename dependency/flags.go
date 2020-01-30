package dependency

import (
	"net"
	"time"
)

// FlagSet is an interface that
type FlagSet interface {
	FullFlagSet
	ShorthandFlagSet
}

// FullFlagSet is the interface for defining full flags
type FullFlagSet interface {
	Uint(name string, value uint, usage string) *uint
	Int16(name string, value int16, usage string) *int16
	UintSlice(name string, value []uint, usage string) *[]uint
	IntSlice(name string, value []int, usage string) *[]int
	Uint16(name string, value uint16, usage string) *uint16
	BoolSlice(name string, value []bool, usage string) *[]bool
	Float32(name string, value float32, usage string) *float32
	IP(name string, value net.IP, usage string) *net.IP
	Duration(name string, value time.Duration, usage string) *time.Duration
	Float64(name string, value float64, usage string) *float64
	IPMask(name string, value net.IPMask, usage string) *net.IPMask
	Int64(name string, value int64, usage string) *int64
	IPNet(name string, value net.IPNet, usage string) *net.IPNet
	Uint64(name string, value uint64, usage string) *uint64
	Uint8(name string, value uint8, usage string) *uint8
	StringToString(name string, value map[string]string, usage string) *map[string]string
	Uint32(name string, value uint32, usage string) *uint32
	StringToInt(name string, value map[string]int, usage string) *map[string]int
	Int64Slice(name string, value []int64, usage string) *[]int64
	Count(name string, usage string) *int
	Int32Slice(name string, value []int32, usage string) *[]int32
	BytesHex(name string, value []byte, usage string) *[]byte
	BytesBase64(name string, value []byte, usage string) *[]byte
	Int8(name string, value int8, usage string) *int8
	StringSlice(name string, value []string, usage string) *[]string
	StringArray(name string, value []string, usage string) *[]string
	Float32Slice(name string, value []float32, usage string) *[]float32
	String(name string, value string, usage string) *string
	Bool(name string, value bool, usage string) *bool
	Float64Slice(name string, value []float64, usage string) *[]float64
	DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration
	Int(name string, value int, usage string) *int
	Int32(name string, value int32, usage string) *int32
	StringToInt64(name string, value map[string]int64, usage string) *map[string]int64
	IPSlice(name string, value []net.IP, usage string) *[]net.IP
}

// ShorthandFlagSet is the interface that allows shorthand flags to be defined
type ShorthandFlagSet interface {
	UintP(name, shorthand string, value uint, usage string) *uint
	Int16P(name, shorthand string, value int16, usage string) *int16
	UintSliceP(name, shorthand string, value []uint, usage string) *[]uint
	IntSliceP(name, shorthand string, value []int, usage string) *[]int
	Uint16P(name, shorthand string, value uint16, usage string) *uint16
	BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool
	Float32P(name, shorthand string, value float32, usage string) *float32
	IPP(name, shorthand string, value net.IP, usage string) *net.IP
	DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration
	Float64P(name, shorthand string, value float64, usage string) *float64
	IPMaskP(name, shorthand string, value net.IPMask, usage string) *net.IPMask
	Int64P(name, shorthand string, value int64, usage string) *int64
	IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet
	Uint64P(name, shorthand string, value uint64, usage string) *uint64
	Uint8P(name, shorthand string, value uint8, usage string) *uint8
	StringToStringP(name, shorthand string, value map[string]string, usage string) *map[string]string
	Uint32P(name, shorthand string, value uint32, usage string) *uint32
	StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int
	Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64
	CountP(name, shorthand string, usage string) *int
	Int32SliceP(name, shorthand string, value []int32, usage string) *[]int32
	BytesHexP(name, shorthand string, value []byte, usage string) *[]byte
	BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte
	Int8P(name, shorthand string, value int8, usage string) *int8
	StringSliceP(name, shorthand string, value []string, usage string) *[]string
	StringArrayP(name, shorthand string, value []string, usage string) *[]string
	Float32SliceP(name, shorthand string, value []float32, usage string) *[]float32
	StringP(name, shorthand string, value string, usage string) *string
	BoolP(name, shorthand string, value bool, usage string) *bool
	Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64
	DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration
	IntP(name, shorthand string, value int, usage string) *int
	Int32P(name, shorthand string, value int32, usage string) *int32
	StringToInt64P(name, shorthand string, value map[string]int64, usage string) *map[string]int64
	IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP
}

// ConfigGetter is the interface that allows config to be retrieved
type ConfigGetter interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
}
