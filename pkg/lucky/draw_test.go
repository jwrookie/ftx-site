package lucky

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPrize(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name  string
		prize string
		err   error
	}

	cases := []TestCase{
		{
			name:  "success",
			prize: Prize40,
		},
		{
			name:  "failed",
			prize: "xxx",
			err:   errors.New("invalid prize name"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckPrize(tc.prize)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestDraw(t *testing.T) {
	t.Parallel()

	err := os.Setenv("UNIT_TEST", "false")
	assert.NoError(t, err)
	t.Run("check prize pool", func(t *testing.T) {
		assert.Equal(t, 1000, len(prizePool))
		t.Log(prizePool)

		pool := map[string]int{}
		for _, v := range prizePool {
			pool[codePrize[v]] = pool[codePrize[v]] + 1
		}

		t.Log(pool)
		assert.Equal(t, prize, pool)
	})

	t.Run("draw", func(t *testing.T) {
		pool := map[string]int64{}

		i := 0
		for i < 10000 {
			prize := Draw()
			pool[prize] = pool[prize] + 1
			i++
		}

		for k := range prizeCode {
			t.Logf("%s : %d", k, pool[k])
		}

		assert.Greater(t, pool[Prize1000], pool[Prize110])
		assert.Greater(t, pool[Prize110], pool[Prize100])
		assert.Greater(t, pool[Prize100], pool[Prize60])
		assert.Greater(t, pool[Prize90], pool[Prize60])
		assert.Greater(t, pool[Prize80], pool[Prize60])
		assert.Greater(t, pool[Prize70], pool[Prize60])
		assert.Greater(t, pool[Prize60], pool[Prize50])
		assert.Greater(t, pool[Prize50], pool[Prize40])
	})
}
