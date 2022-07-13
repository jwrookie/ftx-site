package lucky

import (
	"sync"

	"github.com/foxdex/ftx-site/pkg/db"

	"github.com/foxdex/ftx-site/dao"
)

const (
	initialValue = 5000
	maxValue     = 40000
)

var (
	once           sync.Once
	defaultJackpot *jackpot
)

type jackpot struct {
	mutex  sync.RWMutex
	number uint64
}

func (j *jackpot) getJackpot() uint64 {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	return j.number
}

func (j *jackpot) incJackpot() uint64 {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.number == maxValue {
		return maxValue
	}

	j.number++

	return j.number
}

func InitJackpot(luckyDao dao.ILucky) {
	once.Do(func() {
		count, err := luckyDao.Count(db.Mysql())
		if err != nil {
			panic(err)
		}

		totalCount := initialValue + count
		if totalCount > maxValue {
			totalCount = maxValue
		}
		defaultJackpot = &jackpot{
			number: uint64(totalCount),
		}
	})
}

func GetJackpot() uint64 {
	return defaultJackpot.getJackpot()
}

func IncJackpot() {
	defaultJackpot.incJackpot()
}
