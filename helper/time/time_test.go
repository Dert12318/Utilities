package helperTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	t.Run("a", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2022, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2022-10-10", Now().Format("2006-01-02"))
	})
	t.Run("b", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2021, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2021-10-10", Now().Format("2006-01-02"))
	})
	t.Run("c", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2020, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2020-10-10", Now().Format("2006-01-02"))
	})
}

func TestNow2(t *testing.T) {
	t.Run("a2", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2022, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2022-10-10", Now().Format("2006-01-02"))
	})
	t.Run("b2", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2021, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2021-10-10", Now().Format("2006-01-02"))
	})
	t.Run("c2", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		Mock(time.Date(2020, 10, 10, 0, 0, 0, 0, loc))
		defer ResetMock()

		assert.Equal(t, "2020-10-10", Now().Format("2006-01-02"))
	})
}
