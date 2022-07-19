package lucky

import (
	"errors"
	"math/rand"
	"os"
	"time"
)

const (
	Prize10   = "FTX三周年礼盒"
	Prize20   = "FTX x AMG联合棒球帽"
	Prize30   = "FTX x MLB 棒球外套"
	Prize40   = "交易手续费抵扣券10USD"
	Prize50   = "FTX祝福红包"
	Prize60   = "FTX清凉防晒衣"
	Prize70   = "FTX绒绒袜"
	Prize80   = "FTX小龙人暖手充电宝"
	Prize90   = "FTX雪花真空杯+小金勺子"
	Prize100  = "FTX超萌小耳朵发箍"
	Prize110  = "FTX定制纸牌"
	Prize1000 = "谢谢参与"
)

var codePrize = map[int]string{
	10:  Prize10,
	20:  Prize20,
	30:  Prize30,
	40:  Prize40,
	50:  Prize50,
	60:  Prize60,
	70:  Prize70,
	80:  Prize80,
	90:  Prize90,
	100: Prize100,
	110: Prize110,

	1000: Prize1000,
}

var prizeCode = map[string]int{
	Prize10:   10,
	Prize20:   20,
	Prize30:   30,
	Prize40:   40,
	Prize50:   50,
	Prize60:   60,
	Prize70:   70,
	Prize80:   80,
	Prize90:   90,
	Prize100:  100,
	Prize110:  110,
	Prize1000: 1000,
}

var prize = map[string]int{
	Prize10:   2,
	Prize20:   3,
	Prize30:   4,
	Prize40:   5,
	Prize50:   10,
	Prize60:   30,
	Prize70:   50,
	Prize80:   50,
	Prize90:   50,
	Prize100:  50,
	Prize110:  100,
	Prize1000: 646,
}

var prizePool []int

func init() {
	for k, v := range prize {
		for i := 0; i < v; i++ {
			prizePool = append(prizePool, prizeCode[k])
		}
	}

	rand.Shuffle(len(prizePool), func(i, j int) {
		prizePool[i], prizePool[j] = prizePool[j], prizePool[i]
	})
}

func Draw() string {
	if os.Getenv("UNIT_TEST") == "true" {
		return Prize1000
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Int63n(int64(len(prizePool)))

	return codePrize[prizePool[index]]
}

func CheckPrize(name string) error {
	if _, ok := prize[name]; ok {
		return nil
	}

	return errors.New("invalid prize name")
}
