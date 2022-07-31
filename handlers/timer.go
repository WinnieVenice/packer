package handlers

import (
	"context"
	"fmt"
	"packer/http"
	"packer/model"
	"packer/rpc"
	"packer/util"
	"time"
)

var (
	TimerCacheInterval = 300
	TimerCache         = util.NewLocalCache(TimerCacheInterval)
	TimerGroupIdList   = []int64{376893667}
)

func TimerRecentContest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := rpc.MGetRecentContest(ctx, model.PlaformList)
	if err != nil {
		fmt.Printf("TimerRecentContest, rpc MGetRecentContest failed, err = (%s)\n", err.Error())
		return
	}

	if resp.GetRecentContest() == nil {
		return
	}

	for _, v := range resp.RecentContest {
		if v.GetRecentContest() == nil {
			continue
		}
		for _, c := range v.RecentContest {
			AddContestInTimer(model.ConvertContest(c))
		}
	}

}

func TimerDailyQuestion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := rpc.GetDailyQuestion(ctx, model.PlatformLc)
	if err != nil {
		fmt.Printf("TimerDailyQuestion, rpc GetDailyQuestion failed, err = (%s)\n", err.Error())
		return
	}

	if resp.GetProblem() == nil || len(resp.GetProblem()) <= 0 {
		return
	}

	msg := ""
	for _, p := range resp.Problem {
		problem := model.ConvertProblem(p)
		msg = fmt.Sprintf("%s\n\n%s", msg, problem.String())
	}
	http.MSendGroupMsg(TimerGroupIdList, msg, true)
}

func AddContestInTimer(contest *model.Contest) {
	msg := fmt.Sprintf("%s\n%s\n", contest.String(), time.Now().String())
	endTime := contest.StartTime.Add(contest.Duration)

	if contest.StartTime.Before(time.Now()) {
		return
	}

	if _, ok := TimerCache.Get(contest.Url); ok {
		return
	}
	TimerCache.Set(contest.Url, contest, int(time.Until(endTime)))

	fmt.Printf("TimerRecentContest, add contest = (%s) into timer\n", contest)

	time.AfterFunc(time.Until(contest.StartTime.Add(-time.Hour)), func() {
		http.MSendGroupMsg(TimerGroupIdList, msg, true)
	})
}
