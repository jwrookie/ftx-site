package utils

import (
	"sync"
	"time"

	"github.com/foxdex/ftx-site/pkg/log"

	"go.uber.org/zap"
)

const (
	NanosecondsPerMillisecond        = uint64(time.Millisecond) / uint64(time.Nanosecond)
	TimeSecFormat             string = "2006-01-02 15:04:05"
	TimeDayFormat             string = "2006-01-02"
	TimeZoneChina             string = "Asia/Shanghai"
)

var (
	loadLocationResultMap sync.Map
	loadLocationMutex     sync.Mutex
)

type loadLocationResult struct {
	location *time.Location
	err      error
}

// UnixMilliNow 返回当前的 Unix 时间，使用 Millisecond 作为单位
func UnixMilliNow() uint64 {
	now := time.Now().UnixNano() / int64(NanosecondsPerMillisecond)
	return uint64(now)
}

// ParseMilliToTime 将一个 Unix Millisecond 时间转换为 time.Time 结构体
func ParseMilliToTime(timestamp uint64) time.Time {
	var t = time.Unix(0, 0)
	return t.Add(time.Duration(timestamp * NanosecondsPerMillisecond))
}

// ParseTimeToMilli 将一个 time.Time 结构体 转换为 Unix Millisecond Timestamp
func ParseTimeToMilli(t time.Time) uint64 {
	return uint64(t.UnixNano()) / NanosecondsPerMillisecond
}

// FormatMillisecond 将毫秒类型的时间戳，格式化为指定的format，时区为"Asia/Shanghai"
func FormatMillisecond(timestamp uint64, format string) (string, error) {
	local, err := LoadLocation(TimeZoneChina)
	if err != nil {
		return "", err
	}
	return ParseMilliToTime(timestamp).In(local).Format(format), nil
}

// DayStart 当日开始
func DayStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// DayEnd 当日结束
func DayEnd(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 1e9-1, t.Location())
}

// LoadLocation 实现了带有 cache 的 time.LoadLocation。尽量使用本函数而不是 time.LoadLocation 从而提升系统性能。
//
//
// Go 内置的 time.LoadLocation 会在每次被调用时读取本地的目录，速度比较慢，同时容易造成 too many open files 的错误。
// 可以参考下面的 issue 获取更多信息：
// https://github.com/golang/go/issues/24844
// https://github.com/golang/go/issues/26106
//
func LoadLocation(name string) (*time.Location, error) {
	result, ok := loadLocationResultMap.Load(name)
	if !ok {
		loadLocationMutex.Lock()
		defer loadLocationMutex.Unlock()

		result, ok = loadLocationResultMap.Load(name)
		if ok {
			result := result.(*loadLocationResult)
			return result.location, result.err
		}

		location, err := time.LoadLocation(name)
		result := &loadLocationResult{location: location, err: err}
		loadLocationResultMap.Store(name, result)
		if err != nil {
			log.Log.Error("failed to load location",
				zap.String("name", name),
				zap.Error(err),
			)
		}
		return location, err
	} else {
		result := result.(*loadLocationResult)
		return result.location, result.err
	}
}
