package handlers

import (
	"testing"
)

func TestGetUserContestRecordPicture(t *testing.T) {
	f := HandlerWrapper("GetPic", GetUserContestRecordPicture)
	if f == nil {
		t.Errorf("f is nil")
		return
	}
	_ = f(map[string]string{"group_id": "1111"}, []string{"cf", "tourist"})
}
