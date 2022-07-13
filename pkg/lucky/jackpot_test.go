package lucky

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/foxdex/ftx-site/dao/mock"
	"github.com/golang/mock/gomock"
)

func TestJackpot(t *testing.T) {
	t.Parallel()

	t.Run("use initial value", func(t *testing.T) {
		mockDao := mock.NewMockILucky(gomock.NewController(t))
		mockDao.EXPECT().Count(gomock.Any()).Return(int64(0), nil)

		InitJackpot(mockDao)

		assert.Equal(t, uint64(initialValue), GetJackpot())

		i := 0
		for i < 1000 {
			IncJackpot()
			i++
		}
		assert.Equal(t, uint64(initialValue+1000), GetJackpot())

		i = 0

		for i < maxValue+1 {
			IncJackpot()
			i++
		}
		assert.Equal(t, uint64(maxValue), GetJackpot())
	})

	t.Run("use already exist value", func(t *testing.T) {
		once = sync.Once{}
		mockDao := mock.NewMockILucky(gomock.NewController(t))
		mockDao.EXPECT().Count(gomock.Any()).Return(int64(10000), nil)

		InitJackpot(mockDao)

		assert.Equal(t, uint64(10000+initialValue), GetJackpot())

		i := 0
		for i < 1000 {
			IncJackpot()
			i++
		}
		assert.Equal(t, uint64(10000+initialValue+1000), GetJackpot())

		i = 0

		for i < maxValue+1 {
			IncJackpot()
			i++
		}
		assert.Equal(t, uint64(maxValue), GetJackpot())
	})
}
