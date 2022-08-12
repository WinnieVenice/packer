package qqbot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/WinnieVenice/packer/charts"
	cqc "github.com/WinnieVenice/packer/client/cq"
	"github.com/WinnieVenice/packer/client/crawl"
	"github.com/WinnieVenice/packer/core/timer"
	"github.com/WinnieVenice/packer/model"
	cqs "github.com/WinnieVenice/packer/service/cq"
	"github.com/WinnieVenice/packer/util"
)

func HandlerWrapper(funcMap map[string]string, f model.MsgHandler) model.DefaultHandler {
	handler := model.DefaultHandler{}
	handler.Handler = func(commCtx map[string]string, param []string) error {
		ctx := context.Background()
		err := f(ctx, commCtx, param)
		if err != nil {
			errMsg := fmt.Sprintf("handle err = (%+v)\n", err)
			fmt.Println(errMsg)
			groupIds := getGroupIdsWithCommCtx(commCtx)
			cqc.MSendGroupMsg(groupIds, errMsg, true)
			return err
		}
		return nil
	}
	jsonFuncMap, err := json.Marshal(funcMap)
	if err != nil {
		fmt.Printf("HandlerWrapper marshal funcMap failed, funcMap = (%+v), err = (%+v)\n", funcMap, err)
		return handler
	}
	err = json.Unmarshal(jsonFuncMap, &handler)

	if err != nil {
		fmt.Printf("HandlerWrapper unmarshal funcMap failed, funcMap = (%+v), err = (%+v)\n", funcMap, err)
		return handler
	}

	return handler
}

func SendRecentContest(ctx context.Context, commCtx map[string]string, param []string) error {
	if len(param) <= 0 {
		return fmt.Errorf("param invailed")
	}

	platform := util.MatchPlatform(param[0])

	resp, err := crawl.GetRecentContest(ctx, platform)
	if err != nil {
		fmt.Printf("SendRecentContest, rpc SendRecentContest failed, err = (%s)\n", err.Error())
		return err
	}
	fmt.Println(resp)

	s := ""
	for _, v := range resp.RecentContest {
		s += fmt.Sprintf("%s\n", model.ConvertContest(v).String())
	}
	s += time.Now().String()

	groupIds := getGroupIdsWithCommCtx(commCtx)

	cqc.MSendGroupMsg(groupIds, s, true)

	return nil
}

func SendUserContestRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	if len(param) <= 1 {
		return fmt.Errorf("param invailed")
	}

	userContest := model.UserContestRecord{
		Platform: util.MatchPlatform(param[0]),
		Username: param[1],
	}

	resp, err := crawl.GetUserContestRecord(ctx, userContest)
	if err != nil {
		fmt.Printf("SendUserContestRecord, rpc SendUserContestRecord failed, err = (%s)\n", err.Error())
		return err
	}
	fmt.Println(resp)
	userRecord := model.ConvertUserRecord(resp)
	userRecord.Record = nil
	filePath, fileName := charts.DrawRecord(resp.GetRecord())
	defer time.AfterFunc(time.Second, func() {
		os.Remove(filePath)
	})

	s := fmt.Sprintf("%s\n", userRecord)
	cqCode := fmt.Sprintf("[CQ:image,file=%s/pic/%s]", cqs.GetHostPort(), fileName)
	s = fmt.Sprintf("%s\n%s\n%s\n", s, cqCode, time.Now().String())

	groupIds := getGroupIdsWithCommCtx(commCtx)

	cqc.MSendGroupMsg(groupIds, s, false)

	return nil
}

func GetTimerAllContest(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	s := ""
	timerCacheKV := timer.Cache.GetAllKV()
	cnt := 0
	for _, v := range timerCacheKV {
		cnt++
		contest := v.(*util.Item).Value().(*model.Contest)
		s = fmt.Sprintf("%s\n%s\n", s, contest.String())
		if cnt > 10 {
			s = fmt.Sprintf("%s\n%s", s, time.Now().String())
			cqc.MSendGroupMsg(groupIds, s, true)
			cnt = 0
			s = ""
		}
	}

	if cnt > 0 {
		s = fmt.Sprintf("%s\n%s", s, time.Now().String())
		cqc.MSendGroupMsg(groupIds, s, true)
	}

	return nil
}

func FetchTimerRecentContest(ctx context.Context, commCtx map[string]string, param []string) error {
	timer.RecentContest()

	s := "执行完成\n"
	groupIds := getGroupIdsWithCommCtx(commCtx)

	cqc.MSendGroupMsg(groupIds, s, true)
	return nil
}

func AddTimerGroupId(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)

	succGroupIds := []string{}
	for _, groupId := range groupIds {
		ok := true
		for _, v := range timer.GroupIdList {
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
	for _, v := range timer.GroupIdList {
		nowGroupIds = append(nowGroupIds, strconv.FormatInt(v, 10))
	}
	s := fmt.Sprintf("将群号(%s)加入定时推送列表成功\n现有推送列表:(%s)",
		strings.Join(succGroupIds, " "), strings.Join(nowGroupIds, " "))

	cqc.MSendGroupMsg(groupIds, s, true)

	return nil
}
func getGroupIdsWithCommCtx(commCtx map[string]string) []int64 {
	groupIdList := strings.Split(commCtx["group_id"], " ")
	var groupIds []int64
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

func SendHelp(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	s := ""
	for _, f := range model.MsgHandlerMap {
		s = fmt.Sprintf("%s\n%s\n", s, f.String())
	}
	s = fmt.Sprintf("%s\n%s\n", s, time.Now().String())

	cqc.MSendGroupMsg(groupIds, s, true)
	return nil
}
