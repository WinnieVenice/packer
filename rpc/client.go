package rpc

import (
	"context"
	"fmt"
	"packer/model"
	"packer/pb"
)

var (
	Client *model.RpcClient
)

func NewClient() *model.RpcClient {
	client := &model.RpcClient{}
	client.GetCrawlServiceClient()
	return client
}

func GetClient() *model.RpcClient {
	if Client == nil {
		Client = NewClient()
	}
	return Client
}

func GetRecentContest(ctx context.Context, platform string) (*pb.RecentContest, error) {
	cli := GetClient()
	client, err := (*cli).GetCrawlServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := (*client).GetRecentContest(ctx, &pb.GetRecentContestRequest{
		Platform: platform,
	})
	if err != nil {
		fmt.Printf("get recent contest failed, platform = (%s), err = (%s)\n", platform, err.Error())
		return nil, err
	}

	return resp, nil
}

func MGetRecentContest(ctx context.Context, platformList []string) (*pb.MGetRecentContestResponse, error) {
	cli := GetClient()
	client, err := (*cli).GetCrawlServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := (*client).MGetRecentContest(ctx, &pb.MGetRecentContestRequest{
		Platform: platformList,
	})
	if err != nil {
		fmt.Printf("get recent contest multi failed, platformList = (%+v), err = (%s)\n", platformList, err.Error())
		return nil, err
	}

	return resp, nil
}

func GetDailyQuestion(ctx context.Context, platform string) (*pb.GetDailyQuestionResponse, error) {
	cli := GetClient()
	client, err := (*cli).GetCrawlServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := (*client).GetDailyQuestion(ctx, &pb.GetDailyQuestionRequest{
		Platform: platform,
	})
	if err != nil {
		fmt.Printf("get daily question failed, platform = (%s), err = (%s)\n", platform, err)
		return nil, err
	}

	return resp, nil
}

func GetUserContestRecord(ctx context.Context, userContestRecord model.UserContestRecord) (*pb.UserContestRecord, error) {
	cli := GetClient()
	client, err := (*cli).GetCrawlServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := (*client).GetUserContestRecord(ctx, &pb.GetUserContestRecordRequest{
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

func MGetUserContestRecord(ctx context.Context, userContestRecordList []model.UserContestRecord) (*pb.MGetUserContestRecordResponse, error) {
	cli := GetClient()
	client, err := (*cli).GetCrawlServiceClient()
	if err != nil {
		return nil, err
	}

	reqList := []*pb.GetUserContestRecordRequest{}
	for _, v := range userContestRecordList {
		reqList = append(reqList, &pb.GetUserContestRecordRequest{
			Platform: v.Platform,
			Handle:   v.Username,
		})
	}
	resp, err := (*client).MGetUserContestRecord(ctx, &pb.MGetUserContestRecordRequest{
		GetUserContestRecordRequest: reqList,
	})
	if err != nil {
		fmt.Printf("get user contest record multi failed, userContestList = (%+v), err = (%s)\n",
			userContestRecordList, err.Error())
		return nil, err
	}
	return resp, nil
}
