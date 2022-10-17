package environ

import (
	"log"
	"strings"
	"time"
	"unicode"

	"code.olapie.com/conv"
)

type Manager interface {
	Has(key string) bool
	Get(key string) any
	Set(key string, value any)
	LoadConfigFile(filename string) error
}

var defaultManager Manager = NewDefaultManager()

func Has(key string) bool {
	return defaultManager.Has(key)
}

func Get(key string) any {
	return defaultManager.Get(key)
}

func Set(key string, value any) {
	defaultManager.Set(key, value)
}

func LoadConfigFile(filename string) error {
	return defaultManager.LoadConfigFile(filename)
}

func MustLoadConfigFile(filename string) {
	if err := defaultManager.LoadConfigFile(filename); err != nil {
		log.Panicf("Cannot load config file %s: %v", filename, err)
	}
}

func String(key string, defaultValue string) string {
	s, err := conv.ToString(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return s
}

func MustString(key string) string {
	v, _ := conv.ToString(defaultManager.Get(key))
	if v == "" {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func Int(key string, defaultValue int) int {
	v, err := conv.ToInt(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func MustInt(key string) int {
	v, err := conv.ToInt(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func Int64(key string, defaultValue int64) int64 {
	v, err := conv.ToInt64(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func MustInt64(key string) int64 {
	v, err := conv.ToInt64(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func Float64(key string, defaultValue float64) float64 {
	v, err := conv.ToFloat64(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func MustFloat64(key string) float64 {
	v, err := conv.ToFloat64(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func Duration(key string, defaultValue time.Duration) time.Duration {
	v, err := conv.ToDuration(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func MustDuration(key string) time.Duration {
	v, err := conv.ToDuration(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func Bool(key string, defaultValue bool) bool {
	if !defaultManager.Has(key) {
		return defaultValue
	}
	v, err := conv.ToBool(defaultManager.Get(key))
	if err != nil {

	}
	return v
}

func MustBool(key string) bool {
	v, err := conv.ToBool(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	return v
}

func IntSlice(key string, defaultValue []int) []int {
	if !defaultManager.Has(key) {
		return defaultValue
	}
	v, _ := conv.ToIntSlice(defaultManager.Get(key))
	return v
}

func MustIntSlice(key string) []int {
	v, err := conv.ToIntSlice(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	if len(v) == 0 {
		log.Panicf("%s is empty", key)
	}
	return v
}

func StringSlice(key string, defaultValue []string) []string {
	if !defaultManager.Has(key) {
		return defaultValue
	}
	v, _ := conv.ToStringSlice(defaultManager.Get(key))
	return v
}

func MustStringSlice(key string) []string {
	v, err := conv.ToStringSlice(defaultManager.Get(key))
	if err != nil {
		log.Panicf("%s is not defined", key)
	}
	if len(v) == 0 {
		log.Panicf("%s is empty", key)
	}
	return v
}

// func Map(key string, defaultValue map[string]any) map[string]any {
// 	if !defaultManager.Has(key) {
// 		return defaultValue
// 	}
// 	return cast.ToStringMap(defaultManager.Get(key))
// }

// func MustMap(key string) map[string]any {
// 	if !defaultManager.Has(key) {
// 		log.Panicf("%s is not defined", key)
// 	}
// 	v := cast.ToStringMap(defaultManager.Get(key))
// 	if len(v) == 0 {
// 		log.Panicf("%s is empty", key)
// 	}
// 	return v
// }

func SizeInBytes(key string, defaultValue int) int {
	if !defaultManager.Has(key) {
		return defaultValue
	}
	s, err := conv.ToString(defaultManager.Get(key))
	if err != nil {
		return defaultValue
	}
	return int(parseSizeInBytes(s))
}

func MustSizeInBytes(key string) int {
	if !defaultManager.Has(key) {
		log.Panicf("%s is not defined", key)
	}
	v := defaultManager.Get(key)
	s, err := conv.ToString(v)
	if err != nil {
		log.Panicf("Cast to string %v: %v", v, err)
	}
	if len(s) == 0 {
		log.Panicf("%s is empty", key)
	}
	return int(parseSizeInBytes(s))
}

// parseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes
func parseSizeInBytes(sizeStr string) uint {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size, _ := conv.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return safeMul(uint(size), multiplier)
}

func safeMul(a, b uint) uint {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		return 0
	}
	return c
}
