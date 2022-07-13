package lucky

import (
	"errors"
	"math/rand"
	"os"
	"time"
)

const (
	Prize10   = "FTX x MLB棒球外套"
	Prize20   = "FTX束口袋背包"
	Prize30   = "FTX棒球帽"
	Prize40   = "FTX灰色T恤"
	Prize50   = "交易手续费抵扣券10USD"
	Prize60   = "FTX小龙人暖手充电宝"
	Prize70   = "FTX雪花真空杯"
	Prize80   = "FTX超萌小耳朵发箍 + 小金勺子"
	Prize1000 = "谢谢参与"
)

var codePrize = map[int8]string{
	1: Prize10,
	2: Prize20,
	3: Prize30,
	4: Prize40,
	5: Prize50,
	6: Prize60,
	7: Prize70,
	8: Prize80,
	9: Prize1000,
}

var prizeCode = map[string]int8{
	Prize10:   1,
	Prize20:   2,
	Prize30:   3,
	Prize40:   4,
	Prize50:   5,
	Prize60:   6,
	Prize70:   7,
	Prize80:   8,
	Prize1000: 9,
}

var prize = map[string]int64{
	Prize10:   5,
	Prize20:   30,
	Prize30:   30,
	Prize40:   50,
	Prize50:   50,
	Prize60:   70,
	Prize70:   70,
	Prize80:   100,
	Prize1000: 595,
}

var prizePool []int8

func init() {
	for k, v := range prize {
		for i := 0; i < int(v); i++ {
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
