package mongo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Config struct {
	Mongo struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Database string `toml:"database"`
	} `toml:"mongo"`
}

// Structs for Collections
type Player struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Apply    int           `json:"apply"`
	Editable bool          `json:"editable"`
	Score    []bson.M      `json:"score"`
	Team     mgo.DBRef     `json:"team"`
}

type Team struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Member []struct {
		Player mgo.DBRef `json:"player"`
	} `json:"member"`
	Team    string `json:"team"`
	Defined bool   `json:"defined"`
	Excnt   []int  `json:"excnt"`
}

type Field struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Hole           int           `json:"hole"`
	DrivingContest bool          `json:"drivingContest"`
	Ignore         bool          `json:"ignore"`
	Image          string        `json:"image"`
	NearPin        bool          `json:"nearPin"`
	Par            int           `json:"par"`
	Yard           int           `json:"yard"`
}

// Structs for Page
type UserScore struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
	Hole  int    `json:"hole"`
	Total int    `json:"total"`
}

type Hole struct {
	Hole  int   `json:"hole"`
	Par   int   `json:"par"`
	Yard  int   `json:"yard"`
	Score []int `json:"score"`
}

type Sum struct {
	Par   int   `json:"par"`
	Score []int `json:"score"`
}

type Reaction struct {
	Name        string `json:"mame"`
	ContentType int    `json:"contentType"`
	Content     string `json:"content"`
}

type Thread struct {
	ThreadId  string     `json:"threadId"`
	Msg       string     `json:"msg"`
	ImgUrl    string     `json:"imgUrl"`
	ColorCode string     `json:"colorCode"`
	Positive  bool       `json:"positive"`
	Reactions []Reaction `json:"reactions"`
}

type sortByScore []UserScore

type Index struct {
	Team   []string `json:"team"`
	Length int      `json:"length"`
}

type LeadersBoard struct {
	Ranking []UserScore `json:"ranking"`
}

type ScoreEntrySheet struct {
	Team   string   `json:"team"`
	Hole   int      `json:"hole"`
	Member []string `json:"member"`
	Par    int      `json:"par"`
	Yard   int      `json:"yard"`
	Total  []int    `json:"total"`
	Putt   []int    `json:"putt"`
	Excnt  int      `json:"excnt"`
}

type ScoreViewSheet struct {
	Team    string   `json:"team"`
	Member  []string `json:"member"`
	Apply   []int    `json:"apply"`
	Hole    []Hole   `json:"hole"`
	Sum     Sum      `json:"sum"`
	Defined bool     `json:"defined"`
}

type EntireScore struct {
	Rows [][]string `json:"rows"`
}

type TimeLine struct {
	Threads []Thread `json:"threads"`
}

type PostApplyScore struct {
	Member []string `json:"member"`
	Apply  []int    `json:"apply"`
}

type PostDefinedTeam struct {
	Team string `json:"team"`
}

type PostTeamScore struct {
	Member []string `json:"member"`
	Total  []int    `json:"total"`
	Putt   []int    `json:"putt"`
	Excnt  int      `json:"excnt"`
}

type Status struct {
	Status string `json:"status"`
}
