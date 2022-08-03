package handlers

import (
	"context"
	"fmt"
	"os"
	"packer/dal"
	"packer/http"
	"packer/model"
	"packer/rpc"
	"packer/util"
	"strconv"
	"strings"
	"time"
)

func HandlerWrapper(name string, f model.MsgHandler) model.DefaultHandler {
	return func(commCtx map[string]string, param []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := f(ctx, commCtx, param)
		if err != nil {
			fmt.Printf("handle err = (%+v)\n", err)
			return err
		}
		return nil
	}
}
func GetRecentContest(ctx context.Context, commCtx map[string]string, param []string) error {
	if len(param) <= 0 {
		return fmt.Errorf("param invailed")
	}

	platform := util.MatchPlatform(param[0])

	resp, err := rpc.GetRecentContest(ctx, platform)
	if err != nil {
		fmt.Printf("GetRecentContest, rpc GetRecentContest failed, err = (%s)\n", err.Error())
		return err
	}
	fmt.Println(resp)

	s := ""
	for _, v := range resp.RecentContest {
		s += fmt.Sprintf("%s\n", model.ConvertContest(v).String())
	}
	s += time.Now().String()

	groupIds := getGroupIdsWithCommCtx(commCtx)

	http.MSendGroupMsg(groupIds, s, true)

	return nil
}

func GetUserContestRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	if len(param) <= 1 {
		return fmt.Errorf("param invailed")
	}

	userContest := model.UserContestRecord{
		Platform: util.MatchPlatform(param[0]),
		Username: param[1],
	}

	resp, err := rpc.GetUserContestRecord(ctx, userContest)
	if err != nil {
		fmt.Printf("GetUserContestRecord, rpc GetUserContestRecord failed, err = (%s)\n", err.Error())
		return err
	}
	fmt.Println(resp)
	userRecord := model.ConvertUserRecord(resp)
	userRecord.Record = nil
	filePath, fileName := dal.DrawRecord(resp.GetRecord())
	defer time.AfterFunc(time.Second, func() {
		os.Remove(filePath)
	})

	s := fmt.Sprintf("%s\n", userRecord)
	cqCode := fmt.Sprintf("[CQ:image,file=%s/pic/%s]", http.GetServer().Url, fileName)
	s = fmt.Sprintf("%s\n%s\n%s\n", s, cqCode, time.Now().String())

	groupIds := getGroupIdsWithCommCtx(commCtx)

	http.MSendGroupMsg(groupIds, s, false)

	return nil
}

func GetTimerAllContest(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	s := ""
	timerCacheKV := TimerCache.GetAllKV()
	cnt := 0
	for _, v := range timerCacheKV {
		cnt++
		contest := v.(*util.Item).Value().(*model.Contest)
		s = fmt.Sprintf("%s\n%s\n", s, contest.String())
		if cnt > 10 {
			s = fmt.Sprintf("%s\n%s", s, time.Now().String())
			http.MSendGroupMsg(groupIds, s, true)
			cnt = 0
			s = ""
		}
	}

	if cnt > 0 {
		s = fmt.Sprintf("%s\n%s", s, time.Now().String())
		http.MSendGroupMsg(groupIds, s, true)
	}

	return nil
}

func FetchTimerRecentContest(ctx context.Context, commCtx map[string]string, param []string) error {
	TimerRecentContest()

	s := "执行完成\n"
	groupIds := getGroupIdsWithCommCtx(commCtx)

	http.MSendGroupMsg(groupIds, s, true)
	return nil
}

func AddTimerGroupId(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)

	succGroupIds := []string{}
	for _, groupId := range groupIds {
		ok := true
		for _, v := range TimerGroupIdList {
			if groupId == v {
				ok = false
				break
			}
		}
		if ok {
			succGroupIds = append(succGroupIds, strconv.FormatInt(groupId, 10))
		}
	}
	nowGroupIds := []string{}
	for _, v := range TimerGroupIdList {
		nowGroupIds = append(nowGroupIds, strconv.FormatInt(v, 10))
	}
	s := fmt.Sprintf("将群号(%s)加入定时推送列表成功\n现有推送列表:(%s)",
		strings.Join(succGroupIds, " "), strings.Join(nowGroupIds, " "))

	http.MSendGroupMsg(groupIds, s, true)

	return nil
}
func getGroupIdsWithCommCtx(commCtx map[string]string) []int64 {
	groupIdList := strings.Split(commCtx["group_id"], " ")
	groupIds := []int64{}
	for _, v := range groupIdList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("strconv group_id failed, group_id = (%+v), err = (%+v)\n", commCtx["group_id"], err)
			continue
		}
		groupIds = append(groupIds, id)
	}
	return groupIds
}

func GetAllCommand(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	s := ""
	for k, _ := range model.MsgHandlerMap {
		s = fmt.Sprintf("%s\n%s\n	", s, k)
	}
	s = fmt.Sprintf("%s\n%s\n", s, time.Now().String())

	http.MSendGroupMsg(groupIds, s, true)
	return nil
}
