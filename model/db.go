package model

import (
	"encoding/json"
	"fmt"
)

const (
	TableUserInfo         = "user_infos"
	TablePlatformUserInfo = "platform_user_infos"
)

type UserInfo struct {
	Id           string `gorm:"uniqueIndex" json:"id"`
	NickName     string `gorm:"size:9;not null" json:"nick_name"`
	CodeforcesId string `json:"codeforces_id"`
	AtcoderId    string `json:"atcoder_id"`
	CodechefId   string `json:"codechef_id"`
	NowcoderId   string `json:"nowcoder_id"`
	VjudgeId     string `json:"vjudge_id"`
	LeetcodeId   string `json:"leetcode_id"`
	LuoguId      string `json:"luogu_id"`
}

func (u *UserInfo) Convert(v *UserInfo) {
	j, _ := json.Marshal(*v)
	mp := map[string]string{}
	_ = json.Unmarshal(j, &mp)
	for k, v := range mp {
		if len(v) <= 0 {
			delete(mp, k)
		}
	}
	j, _ = json.Marshal(mp)
	_ = json.Unmarshal(j, &u)
}

func (u *UserInfo) ConvertMap(mp map[string]interface{}) {
	j, _ := json.Marshal(mp)
	_ = json.Unmarshal(j, &u)
}
func (u *UserInfo) InvConvertMap() map[string]string {
	mp := make(map[string]interface{})
	j, _ := json.Marshal(u)
	_ = json.Unmarshal(j, &mp)
	ret := make(map[string]string)
	for _, v := range PlaformList {
		if id, ok := mp[fmt.Sprintf("%s_id", v)]; ok && id.(string) != "" {
			ret[v] = id.(string)
		}
	}
	return ret
}
func (u *UserInfo) String() string {
	s := ""
	tmpIds := []string{u.CodeforcesId, u.AtcoderId, u.NowcoderId, u.LuoguId, u.VjudgeId, u.LeetcodeId, u.CodechefId}
	for i := range PlaformList {
		if len(tmpIds[i]) > 0 && tmpIds[i] != "" {
			s = fmt.Sprintf("%s\n%s:%s", s, PlaformList[i], tmpIds[i])
		}
	}
	return s
}

type PlatformUserInfo struct {
	// Id plantform|id
	Id         string `gorm:"uniqueIndex" json:"id"`
	Accept     string
	Submit     string
	PrevAccept string
	PrevSubmit string
}

func (u *PlatformUserInfo) Convert(v *PlatformUserInfo) {
	j, _ := json.Marshal(*v)
	mp := map[string]string{}
	_ = json.Unmarshal(j, &mp)
	for k, v := range mp {
		if len(v) <= 0 {
			delete(mp, k)
		}
	}
	j, _ = json.Marshal(mp)
	_ = json.Unmarshal(j, &u)
}

func (u *PlatformUserInfo) String() string {
	return fmt.Sprintf("当前: 提交数量: %s, ac数量: %s\n昨日: 提交数量: %s, ac数量%s", u.Submit, u.Accept, u.PrevSubmit, u.PrevAccept)
}
