package tests

import (
	"testing"
	"time"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	nowMs := model.GetMillis()
	assert.Greater(t, nowMs, int64(0))

	t1 := time.Now()
	millis := model.GetMillisForTime(t1)
	assert.Equal(t, t1.UnixNano()/1000000, millis)

	t2 := model.GetTimeForMillis(millis)
	// they should be within 1 millisecond
	diff := t2.Sub(t1)
	if diff < 0 {
		diff = -diff
	}
	assert.LessOrEqual(t, diff, time.Millisecond)
}
