package crawl

import (
	"context"
	"fmt"

	"github.com/WinnieVenice/packer/idl"
	"github.com/WinnieVenice/packer/model"
)

func GetRecentContest(ctx context.Context, platform string) (*idl.RecentContest, error) {
	resp, err := Client().GetRecentContest(ctx, &idl.GetRecentContestRequest{
		Platform: platform,
	})
	if err != nil {
		fmt.Printf("get recent contest failed, platform = (%s), err = (%s)\n", platform, err.Error())
		return nil, err
	}
	return resp, nil
}

func MGetRecentContest(ctx context.Context, platformList []string) (*idl.MGetRecentContestResponse, error) {
	resp, err := Client().MGetRecentContest(ctx, &idl.MGetRecentContestRequest{
		Platform: platformList,
	})
	if err != nil {
		fmt.Printf("get recent contest multi failed, platformList = (%+v), err = (%s)\n", platformList, err.Error())
		return nil, err
	}
	return resp, nil
}

func GetDailyQuestion(ctx context.Context, platform string) (*idl.GetDailyQuestionResponse, error) {
	resp, err := Client().GetDailyQuestion(ctx, &idl.GetDailyQuestionRequest{
		Platform: platform,
	})
	if err != nil {
		fmt.Printf("get daily question failed, platform = (%s), err = (%s)\n", platform, err)
		return nil, err
	}

	return resp, nil
}

func GetUserContestRecord(ctx context.Context, userContestRecord model.UserContestRecord) (*idl.UserContestRecord, error) {
	resp, err := Client().GetUserContestRecord(ctx, &idl.GetUserContestRecordRequest{
		Platform: userContestRecord.Platform,
		Handle:   userContestRecord.Username,
	})
	if err != nil {
		fmt.Printf("get user contest record failed, userContest = (%+v), err = (%s)\n",
			userContestRecord, err.Error())
		return nil, err
	}

	return resp, nil
}

func MGetUserContestRecord(ctx context.Context, userContestRecordList []model.UserContestRecord) (*idl.MGetUserContestRecordResponse, error) {
	var reqList []*idl.GetUserContestRecordRequest
	for _, v := range userContestRecordList {
		reqList = append(reqList, &idl.GetUserContestRecordRequest{
			Platform: v.Platform,
			Handle:   v.Username,
		})
	}
	resp, err := Client().MGetUserContestRecord(ctx, &idl.MGetUserContestRecordRequest{
		GetUserContestRecordRequest: reqList,
	})
	if err != nil {
		fmt.Printf("get user contest record multi failed, userContestList = (%+v), err = (%s)\n",
			userContestRecordList, err.Error())
		return nil, err
	}
	return resp, nil
}
