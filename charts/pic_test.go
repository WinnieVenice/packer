package charts

import (
	"testing"
	"time"

	"github.com/WinnieVenice/packer/idl"
)

func TestDrawRecord(t *testing.T) {
	req := []*idl.UserContestRecord_Record{
		{
			Name:      "Test",
			Timestamp: time.Now().Unix(),
			Rating:    1145,
		},
	}
	t.Log(DrawRecord(req))
}

func TestDrawBindUserDailyDiff(t *testing.T) {
	req := map[string]int64{
		"xiaozhan": 11,
		"ikun":     22,
	}

	t.Log(DrawBindUserDailyDiff(req))
}
