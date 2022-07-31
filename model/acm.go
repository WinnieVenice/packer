package model

import (
	"fmt"
	"packer/pb"
	"time"
)

type Contest struct {
	Name      string
	Url       string
	StartTime time.Time
	Duration  time.Duration
}

type Problem struct {
	Platform   string
	Url        string
	Id         string
	Name       string
	Difficulty string
}

type UserRecord struct {
	ProfileUrl string // 用户页面URL
	Rating     int32  // 当前rating 没有参加过则为0
	Length     int32  // 参加比赛场次
	Record     []*ContestRecord
	Platform   string
	Username   string
}

type ContestRecord struct {
	Name      string    // 比赛名称
	Url       string    // 比赛链接
	StartTime time.Time // 比赛时间 单位：秒
	Rating    int32     // 结果rating
}
type UserContestRecord struct {
	Username string
	Platform string
}

func ConvertContest(c *pb.RecentContest_ContestMessage) *Contest {
	r := Contest{}
	r.Name = c.Name
	r.Url = c.Url
	r.StartTime = time.Unix(c.Timestamp, 0)
	r.Duration = time.Unix(int64(c.Duration), 0).Sub(time.Unix(0, 0))
	return &r
}

func (c *Contest) String() string {
	s := ""
	s += fmt.Sprintf("比赛名字: %s\n", c.Name)
	s += fmt.Sprintf("比赛地址: %s\n", c.Url)
	s += fmt.Sprintf("开始时间: %+v\n", c.StartTime)
	s += fmt.Sprintf("持续时间: %+v", c.Duration)
	return s
}

func ConvertProblem(p *pb.GetDailyQuestionResponse_Problem) *Problem {
	r := Problem{}
	r.Platform = p.Platform
	r.Url = p.Url
	r.Id = p.Id
	r.Name = p.Name
	r.Difficulty = p.Difficulty
	return &r
}

func (p *Problem) String() string {
	s := ""
	s += fmt.Sprintf("题目平台: %s\n", p.Platform)
	s += fmt.Sprintf("题目地址: %s\n", p.Url)
	s += fmt.Sprintf("题目ID: %s\n", p.Id)
	s += fmt.Sprintf("题目名称: %s\n", p.Name)
	s += fmt.Sprintf("题目难度: %s", p.Difficulty)
	return s
}

func ConvertUserRecord(ur *pb.UserContestRecord) *UserRecord {
	r := UserRecord{}
	r.ProfileUrl = ur.ProfileUrl
	r.Rating = ur.Rating
	r.Length = ur.Length
	r.Record = []*ContestRecord{}
	for _, v := range ur.Record {
		if ur.Record == nil {
			continue
		}
		r.Record = append(r.Record, ConvertContestRecord(v))
	}
	r.Platform = ur.Platform
	r.Username = ur.Handle
	return &r
}

func (ur *UserRecord) String() string {
	s := ""
	s += fmt.Sprintf("平台: %s\n", ur.Platform)
	s += fmt.Sprintf("用户名: %s\n", ur.Username)
	s += fmt.Sprintf("个人页面: %s\n", ur.ProfileUrl)
	s += fmt.Sprintf("Rating: %+v\n", ur.Rating)
	s += fmt.Sprintf("参赛场次: %+v\n", ur.Length)
	s += fmt.Sprintf("比赛记录: %s\n", BlankSpace)
	for _, v := range ur.Record {
		s += fmt.Sprintf("%s\n%s\n", (*v).String(), BlankSpace)
	}
	return s
}

func ConvertContestRecord(cr *pb.UserContestRecord_Record) *ContestRecord {
	r := ContestRecord{}
	r.Name = cr.Name
	r.Url = cr.Url
	r.StartTime = time.Unix(cr.Timestamp, 0)
	r.Rating = cr.Rating
	return &r
}

func (cr *ContestRecord) String() string {
	s := ""
	s += fmt.Sprintf("比赛名称: %s\n", cr.Name)
	s += fmt.Sprintf("比赛地址: %s\n", cr.Url)
	s += fmt.Sprintf("比赛时间: %s\n", cr.StartTime)
	s += fmt.Sprintf("结束Rating: %+v\n", cr.Rating)
	return s
}
