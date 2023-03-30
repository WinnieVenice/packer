package qqbot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/WinnieVenice/packer/backend"
	"github.com/WinnieVenice/packer/charts"
	cqc "github.com/WinnieVenice/packer/client/cq"
	"github.com/WinnieVenice/packer/client/crawl"
	"github.com/WinnieVenice/packer/core/timer"
	"github.com/WinnieVenice/packer/model"
	cqs "github.com/WinnieVenice/packer/service/cq"
	"github.com/WinnieVenice/packer/util"
)

func getGroupIdsWithCommCtx(commCtx map[string]string) []int64 {
	groupIdList := strings.Split(commCtx["group_id"], " ")
	var groupIds []int64
	for _, v := range groupIdList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("strconv group_id failed, group_id = (%+v), err = (%+v)\n", commCtx["group_id"], err)
			continue
		}
		if id <= 0 {
			continue
		}
		groupIds = append(groupIds, id)
	}
	return groupIds
}

func getMsgWrapper(uid int64, msg string) string {
	s := ""
	if uid > 0 {
		s = model.CqCodeAt(uid)
	}
	s = fmt.Sprintf("%s\n%s\n%s", s, msg, time.Now().Format("2006-01-02 15:04:05"))
	return s
}

func HandlerWrapper(funcMap map[string]string, f model.MsgHandler) model.DefaultHandler {
	handler := model.DefaultHandler{}
	handler.Handler = func(commCtx map[string]string, param []string) error {
		ctx := context.Background()
		err := f(ctx, commCtx, param)
		if err != nil {
			uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
			errMsg := getMsgWrapper(uId, fmt.Sprintf("handle err = (%+v)\n", err))
			fmt.Println(errMsg)
			groupIds := getGroupIdsWithCommCtx(commCtx)
			if len(groupIds) > 0 {
				cqc.MSendGroupMsg(groupIds, errMsg, false)
			} else {
				cqc.SendPrivateMsg(uId, errMsg, false)
			}
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
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	s = getMsgWrapper(uId, s)

	groupIds := getGroupIdsWithCommCtx(commCtx)

	cqc.MSendGroupMsg(groupIds, s, false)

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
	s = fmt.Sprintf("%s\n%s\n", s, cqCode)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	groupIds := getGroupIdsWithCommCtx(commCtx)
	s = getMsgWrapper(uId, s)

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
			s = getMsgWrapper(0, s)
			cqc.MSendGroupMsg(groupIds, s, false)
			cnt = 0
			s = ""
		}
	}

	if cnt > 0 {
		s = getMsgWrapper(0, s)
		cqc.MSendGroupMsg(groupIds, s, false)
	}

	return nil
}

func FetchTimerRecentContest(ctx context.Context, commCtx map[string]string, param []string) error {
	timer.RecentContest()

	s := "执行完成\n"
	groupIds := getGroupIdsWithCommCtx(commCtx)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	s = getMsgWrapper(uId, s)

	cqc.MSendGroupMsg(groupIds, s, false)
	return nil
}

func AddTimerGroupId(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)

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
	s = getMsgWrapper(uId, s)

	cqc.MSendGroupMsg(groupIds, s, false)

	return nil
}

func SendHelp(ctx context.Context, commCtx map[string]string, param []string) error {
	groupIds := getGroupIdsWithCommCtx(commCtx)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	s := ""
	for _, f := range model.MsgHandlerMap {
		s = fmt.Sprintf("%s\n%s\n", s, f.String())
	}
	s = getMsgWrapper(uId, s)

	cqc.MSendGroupMsg(groupIds, s, false)
	return nil
}

func InitUser(ctx context.Context, commCtx map[string]string, param []string) error {
	if len(param) <= 0 {
		return fmt.Errorf("param is invalid")
	}
	uIdStr := commCtx["user_id"]
	uId, _ := strconv.ParseInt(uIdStr, 10, 64)
	user, err := backend.QueUserById(uIdStr)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("user is inited")
	}
	nickName := param[1]
	userInfo := model.UserInfo{
		Id:       uIdStr,
		NickName: nickName,
	}
	if err := backend.AddUser(&userInfo); err != nil {
		util.Println("init user failed, err = (%+v)", err)
		return err
	}
	s := getMsgWrapper(uId, "初始化用户成功")
	groupIds := getGroupIdsWithCommCtx(commCtx)
	if len(groupIds) > 0 {
		cqc.MSendGroupMsg(groupIds, s, false)
	} else {
		cqc.SendPrivateMsg(uId, s, false)
	}
	return nil
}

func BindUser(ctx context.Context, commCtx map[string]string, param []string) error {
	uIdStr := commCtx["user_id"]
	uId, err := strconv.ParseInt(uIdStr, 10, 64)
	if err != nil {
		return err
	}
	groupIds := getGroupIdsWithCommCtx(commCtx)
	userMap := map[string]interface{}{
		"id": uIdStr,
	}
	for _, v := range param {
		x := strings.Split(v, ":")
		if len(x) != 2 || len(x[0]) <= 0 || len(x[1]) <= 0 {
			continue
		}
		platform, id := util.MatchPlatform(x[0]), x[1]

		if resp, err := crawl.GetUserSubmitRecord(ctx, platform, id); err != nil {
			return err
		} else {
			platformUserInfo := model.PlatformUserInfo{
				Id:         strings.Join([]string{platform, id}, "|"),
				Accept:     strconv.FormatInt(int64(resp.AcceptCount), 10),
				Submit:     strconv.FormatInt(int64(resp.SubmitCount), 10),
				PrevAccept: strconv.FormatInt(int64(resp.AcceptCount), 10),
				PrevSubmit: strconv.FormatInt(int64(resp.SubmitCount), 10),
			}
			_ = backend.AddPlatformUser(&platformUserInfo)

			userMap[fmt.Sprintf("%s_id", platform)] = id

		}
	}
	userInfo := model.UserInfo{}
	userInfo.ConvertMap(userMap)
	fmt.Println(userInfo.String())

	if err = backend.UpdUser(&userInfo); err != nil {
		return err
	}

	s := getMsgWrapper(uId, "绑定成功")

	if len(groupIds) > 0 {
		cqc.MSendGroupMsg(groupIds, s, false)
	} else {
		cqc.SendPrivateMsg(uId, s, false)
	}

	return nil
}

func QueryUser(ctx context.Context, commCtx map[string]string, param []string) error {
	uIdStr := commCtx["user_id"]
	uId, err := strconv.ParseInt(uIdStr, 10, 64)
	if err != nil {
		return err
	}

	userInfo, err := backend.QueUserById(uIdStr)
	if err != nil {
		return err
	}

	groupIds := getGroupIdsWithCommCtx(commCtx)
	s := getMsgWrapper(uId, userInfo.String())

	if len(groupIds) > 0 {
		cqc.MSendGroupMsg(groupIds, s, false)
	} else {
		cqc.SendPrivateMsg(uId, s, false)
	}

	return nil
}

func MGetUserSubmitRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	uIdStr := commCtx["user_id"]
	/*
		uId, err := strconv.ParseInt(uIdStr, 10, 64)
		if err != nil {
			return err
		}
	*/

	var userList []model.UserContestRecord
	if len(param) <= 1 {
		userInfo, err := backend.QueUserById(uIdStr)
		if err != nil {
			return err
		}

		IdMap := userInfo.InvConvertMap()
		if len(param) == 1 {
			platform := util.MatchPlatform(param[0])
			if id, ok := IdMap[platform]; ok {
				userList = append(userList, model.UserContestRecord{
					Platform: platform,
					Username: id,
				})
			}
		} else {
			for _, v := range model.PlaformList {
				if id, ok := IdMap[v]; ok {
					userList = append(userList, model.UserContestRecord{
						Platform: v,
						Username: id,
					})
				}
			}
		}
	} else {
		for i := range param {
			if i+1 >= len(param) {
				continue
			}

			platform, id := util.MatchPlatform(param[i]), param[i+1]
			userList = append(userList, model.UserContestRecord{
				Platform: platform,
				Username: id,
			})
		}
	}

	resp, err := crawl.MGetUserSubmitRecord(ctx, userList)
	if err != nil {
		fmt.Printf("MGetUserSubmitRecord, rpc MGetUserSubmitRecord failed, err = (%+v)", err)
		return nil
	}
	fmt.Println(resp)

	return nil
}

func FetchTimerPlatformUserRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	timer.TimerPlatformUserRecord()

	s := "执行完成\n"
	groupIds := getGroupIdsWithCommCtx(commCtx)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	s = getMsgWrapper(uId, s)

	cqc.MSendGroupMsg(groupIds, s, false)
	return nil
}

func GetBindUserRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	uIdStr := commCtx["user_id"]
	uId, err := strconv.ParseInt(uIdStr, 10, 64)
	if err != nil {
		return err
	}

	userInfo, err := backend.QueUserById(uIdStr)
	if err != nil {
		return err
	}

	s := ""
	for k, v := range userInfo.InvConvertMap() {
		plaform, id := strings.ReplaceAll(k, "_id", ""), v
		platformId := fmt.Sprintf("%s|%s", plaform, id)

		user, err := backend.QuePlatformUserById(platformId)
		if err != nil {
			fmt.Printf("GetBindUserRecord, QuePlaformUserById failed, err = (%+v)", err)
			continue
		}

		s = fmt.Sprintf("%s\n%s", s, user.String())
	}

	groupIds := getGroupIdsWithCommCtx(commCtx)
	s = getMsgWrapper(uId, s)

	if len(groupIds) > 0 {
		cqc.MSendGroupMsg(groupIds, s, false)
	} else {
		cqc.SendPrivateMsg(uId, s, false)
	}

	return nil
}

func FetchTimerBindUserDailyRecord(ctx context.Context, commCtx map[string]string, param []string) error {
	timer.TimerBindUserDailyRecord()

	s := "执行完成\n"
	groupIds := getGroupIdsWithCommCtx(commCtx)
	uId, _ := strconv.ParseInt(commCtx["user_id"], 10, 64)
	s = getMsgWrapper(uId, s)

	cqc.MSendGroupMsg(groupIds, s, false)
	return nil
}
