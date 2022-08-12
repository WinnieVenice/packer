package timer

import (
	"context"
	"fmt"
	"time"

	"github.com/WinnieVenice/packer/client/cq"
	"github.com/WinnieVenice/packer/client/crawl"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/util"
)

var (
	CacheInterval = 300
	Cache         = util.NewLocalCache(CacheInterval)
	GroupIdList   = []int64{376893667}
)

func RecentContest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := crawl.MGetRecentContest(ctx, model.PlaformList)
	if err != nil {
		fmt.Printf("RecentContest, rpc MGetRecentContest failed, err = (%s)\n", err.Error())
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

func DailyQuestion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := crawl.GetDailyQuestion(ctx, model.PlatformLc)
	if err != nil {
		fmt.Printf("DailyQuestion, rpc GetDailyQuestion failed, err = (%s)\n", err.Error())
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
	cq.MSendGroupMsg(GroupIdList, msg, true)
}

func AddContestInTimer(contest *model.Contest) {
	msg := fmt.Sprintf("%s\n%s\n", contest.String(), time.Now().String())
	endTime := contest.StartTime.Add(contest.Duration)

	if contest.StartTime.Before(time.Now()) {
		return
	}

	if _, ok := Cache.Get(contest.Url); ok {
		return
	}

	Cache.Set(contest.Url, contest, time.Until(endTime))

	fmt.Printf("RecentContest, add contest = (%s) into timer\n", contest)

	time.AfterFunc(time.Until(contest.StartTime.Add(-time.Hour)), func() {
		cq.MSendGroupMsg(GroupIdList, msg, true)
	})
}
