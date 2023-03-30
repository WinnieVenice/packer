package timer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WinnieVenice/jdsd"
	"github.com/WinnieVenice/packer/backend"
	"github.com/WinnieVenice/packer/charts"
	"github.com/WinnieVenice/packer/client/cq"
	"github.com/WinnieVenice/packer/client/crawl"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/util"
	"gorm.io/gorm"

	cqc "github.com/WinnieVenice/packer/client/cq"
	cqs "github.com/WinnieVenice/packer/service/cq"
)

var (
	CacheInterval = 300
	Cache         = util.NewLocalCache(CacheInterval)
	GroupIdList   = []int64{376893667}
	DailyCache    = util.NewLocalCache(CacheInterval)
)

func AddContestInTimer(contest *model.Contest) {
	endTime := contest.StartTime.Add(contest.Duration)
	now := time.Now()
	if contest.StartTime.Before(now) {
		return
	}

	if _, ok := Cache.Get(contest.Url); !ok {
		fmt.Printf("RecentContest, add contest = (%s) into timer\n", contest)
		time.AfterFunc(time.Until(contest.StartTime.Add(-time.Hour)), func() {
			if value, ok := Cache.Get(contest.Url); ok {
				c := value.(*model.Contest)
				if c.StartTime.Add(-2 * time.Hour).Before(time.Now()) {
					msg := c.String()
					cq.MSendGroupMsg(GroupIdList, msg, true)
				}
			}
		})
	}

	Cache.Set(contest.Url, contest, time.Until(endTime))
}

func RecentContest() {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()

	resp, err := crawl.MGetRecentContest(ctx, []string{})
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
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()

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

func TimerPlatformUserRecord() {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()

	users, err := backend.GetAllPlatformUser()
	if err != nil {
		fmt.Printf("TimerUserRecord, GetAllPlatform failed, err = (%+v)\n", err)
		return
	}

	mp := make(map[string]interface{})
	var reqList []model.UserContestRecord
	for _, user := range users {
		strs := strings.Split(user.Id, "|")
		platform, id := strs[0], strs[1]
		reqList = append(reqList, model.UserContestRecord{
			Platform: platform,
			Username: id,
		})

		j, _ := json.Marshal(user)
		mp[user.Id] = j
	}

	resp, err := crawl.MGetUserSubmitRecord(ctx, reqList)
	if err != nil {
		fmt.Printf("TimerUserRecord, rpc MGetUserSubmitRecord failed, err = (%+v)\n", err)
		return
	}

	backend.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, v := range resp.UserSubmitRecord {
			if v == nil {
				continue
			}

			platformId := strings.Join([]string{v.Platform, v.Handle}, "|")
			accept, submit := strconv.FormatInt(int64(v.AcceptCount), 10), strconv.FormatInt(int64(v.SubmitCount), 10)
			j := mp[platformId]
			var oldUser model.PlatformUserInfo
			_ = json.Unmarshal(j.([]byte), &oldUser)

			user := &model.PlatformUserInfo{
				Id:         platformId,
				Accept:     accept,
				Submit:     submit,
				PrevAccept: oldUser.Accept,
				PrevSubmit: oldUser.Submit,
			}
			if err := backend.UpdPlatformUser(user); err != nil {
				fmt.Printf("TimerUserRecord, UpdPlatformUser failed, RollBack, user = (%+v), err = (%+v)\n", user, err)
				return err
			}

		}

		return nil
	})
}

func TimerBindUserDailyRecord() {
	users, err := backend.GetAllUser()
	if err != nil {
		fmt.Printf("TimerBindUserDailyRecord, GetAllUser failed, err = (%+v)\n", err)
		return
	}

	dailyRecord := make(map[string]int64)
	for _, user := range users {
		mp := user.InvConvertMap()
		sumDiffAc, sumDiffSub := int64(0), int64(0)

		for k, v := range mp {
			plaform, id := strings.ReplaceAll(k, "_id", ""), v
			plaformId := fmt.Sprintf("%s|%s", plaform, id)

			var plaformUser *model.PlatformUserInfo
			if val, ok := DailyCache.Get(plaformId); ok {
				_ = json.Unmarshal(val.([]byte), plaformUser)
			} else {
				plaformUser, err = backend.QuePlatformUserById(plaformId)
				if err != nil {
					fmt.Println(err)
					continue
				}
				util.Println("query sql, platform_user = (%+v)", plaformUser)
			}

			ac, _ := strconv.ParseInt(plaformUser.Accept, 10, 64)
			prevAc, _ := strconv.ParseInt(plaformUser.PrevAccept, 10, 64)
			sub, _ := strconv.ParseInt(plaformUser.Submit, 10, 64)
			prevSub, _ := strconv.ParseInt(plaformUser.PrevSubmit, 10, 64)

			sumDiffAc += ac - prevAc
			sumDiffSub += sub - prevSub
		}

		if sumDiffSub > 0 {
			dailyRecord[user.Id] = sumDiffSub
		}
	}
	fmt.Println("daily_record = ", dailyRecord)
	s := ""
	if dailyRecord == nil || len(dailyRecord) <= 0 {
		now := time.Now()
		pre := now
		if value, ok := DailyCache.Get("no_pratice"); ok && value != nil {
			if t, ok := value.(time.Time); ok {
				pre = t
			}
		}
		text := fmt.Sprintf("今日%s无人训练, 距离距离上一次无人训练%s, 已有%d天, 警钟长鸣, 望周知",
			now.Format("2006-01-02 15:04:05"), pre.Format("2006-01-02 15:04:05"), int64(now.Sub(pre).Hours()/24))
		s = fmt.Sprintf("%s\n%s\n", s, text)
	} else {
		filePath, fileName := charts.DrawBindUserDailyDiff(dailyRecord)
		defer time.AfterFunc(time.Minute, func() {
			os.Remove(filePath)
		})

		cqCode := fmt.Sprintf("[CQ:image,file=%s/pic/%s]", cqs.GetHostPort(), fileName)
		s = fmt.Sprintf("%s\n%s\n", s, cqCode)
	}

	fmt.Println("s=", s)
	cqc.MSendGroupMsg(GroupIdList, s, false)
}

func TimerJDSD() {
	s, err := jdsd.Run()
	if err != nil {
		s = fmt.Sprintf("error = (%+v)", err)
	}
	cqc.MSendGroupMsg(GroupIdList, s, false)
}
