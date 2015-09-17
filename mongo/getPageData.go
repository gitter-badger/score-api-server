package mongo

import (
	"sort"
	"strconv"
)

var (
	players []Player
	fields  []Field
	teams   []Team
)

func init() {
	players = GetAllPlayerCol()
	fields = GetAllFieldCol()
	teams = GetAllTeamCol()
}

func GetIndexPageData() (*Index, error) {

	teamArray := make([]string, len(teams))
	for i, team := range teams {
		teamArray[i] = team.Team
	}

	idx := &Index{
		Team:   teamArray,
		Length: len(teams),
	}
	return idx, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {

	var leadersBoard LeadersBoard
	var userScore UserScore

	for _, player := range players {
		var totalPar, totalScore, passedHoleCnt int
		for holeIndex, playerScore := range player.Score {
			if playerScore["total"] != 0 {
				totalPar += fields[holeIndex].Par
				totalScore += playerScore["total"].(int)
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore - totalPar,
			Total: totalScore,
			Name:  player.Name,
			Hole:  passedHoleCnt,
		}
		leadersBoard.Ranking = append(leadersBoard.Ranking, userScore)
	}
	sort.Sort(sortByScore(leadersBoard.Ranking))

	return &leadersBoard, nil
}

func GetScoreEntrySheetPageData(teamName string, holeString string) (*ScoreEntrySheet, error) {
	if len(holeString) == 0 {
		return nil, nil
	}
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	field := fields[holeIndex]
	playersInTheTeam := GetPlayersDataInTheTeam(teamName)

	member := make([]string, 4)
	stroke := make([]int, 4)
	putt := make([]int, 4)
	for playerIndex, player := range playersInTheTeam {
		member[playerIndex] = player.Name
		stroke[playerIndex] = player.Score[holeIndex]["stroke"].(int)
		putt[playerIndex] = player.Score[holeIndex]["putt"].(int)
	}

	scoreEntrySheet := ScoreEntrySheet{
		Team:   teamName,
		Hole:   holeNum,
		Member: member,
		Par:    field.Par,
		Yard:   field.Yard,
		Stroke: stroke,
		Putt:   putt,
		Excnt:  0,
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {

	playersInTheTeam := GetPlayersDataInTheTeam(teamName)

	member := make([]string, len(playersInTheTeam))
	apply := make([]int, len(playersInTheTeam))
	holes := make([]Hole, len(fields))
	totalScore := make([]int, len(playersInTheTeam))
	var totalPar int
	for holeIndex, field := range fields {

		totalPar += field.Par
		score := make([]int, len(playersInTheTeam))
		for playerIndex, player := range playersInTheTeam {
			score[playerIndex] = player.Score[holeIndex]["total"].(int)
			totalScore[playerIndex] += score[playerIndex]
			if holeIndex == 0 {
				member[playerIndex] = player.Name
				apply[playerIndex] = player.Apply
			}
		}
		holes[holeIndex] = Hole{
			Hole:  holeIndex + 1,
			Par:   field.Par,
			Yard:  field.Yard,
			Score: score,
		}
	}

	sum := Sum{
		Par:   totalPar,
		Score: totalScore,
	}
	scoreViewSheet := ScoreViewSheet{
		Team:   teamName,
		Member: member,
		Apply:  apply,
		Hole:   holes,
		Sum:    sum,
	}

	return &scoreViewSheet, nil
}

func GetEntireScorePageData() (*EntireScore, error) {

	rowSize := 25
	columnSize := len(players) + 2
	rankMap := make(map[string]bool)

	rows := make([][]string, rowSize)
	for i := 0; i < rowSize; i++ {
		if i == 0 {
			rows[i] = make([]string, len(teams)*2)
		} else {
			rows[i] = make([]string, columnSize)
		}
	}
	mainColumnNum := 2

	rows[1][0] = "ホール"
	rows[1][1] = "パー"
	rows[20][0] = "Gross"
	rows[20][1] = "-"
	rows[21][0] = "Net"
	rows[21][1] = "-"
	rows[22][0] = "申請"
	rows[22][1] = "-"
	rows[23][0] = "スコア差"
	rows[23][1] = "-"
	rows[24][0] = "順位"
	rows[24][1] = "-"

	var passedPlayerNum int
	for teamIndex, team := range teams {
		rows[0][teamIndex*2] = strconv.Itoa(len(team.Member))
		rows[0][teamIndex*2+1] = "TEAM " + team.Team
		playerInTheTeam := GetPlayersDataInTheTeam(team.Team)
		for playerIndex, player := range playerInTheTeam {
			userDataIndex := playerIndex + passedPlayerNum + mainColumnNum
			rows[1][userDataIndex] = player.Name

			var gross int
			var net int
			for holeIndex, field := range fields {
				holeRowNum := holeIndex + 2
				if playerIndex == 0 {
					if field.Ignore {
						rows[holeRowNum][0] = "-i" + strconv.Itoa(field.Hole)
					} else {
						rows[holeRowNum][0] = strconv.Itoa(field.Hole)
					}
					rows[holeRowNum][1] = strconv.Itoa(field.Par)
				}
				playerTotal := player.Score[holeIndex]["total"].(int)
				gross += playerTotal
				if field.Ignore == false {
					net += playerTotal
				}
				rows[holeRowNum][userDataIndex] = strconv.Itoa(playerTotal)
			}
			rows[20][userDataIndex] = strconv.Itoa(gross)
			rows[21][userDataIndex] = strconv.Itoa(net)
			rows[22][userDataIndex] = strconv.Itoa(player.Apply)
			diff := player.Apply - net
			if diff < 0 {
				diff = diff * -1
			}
			rows[23][userDataIndex] = strconv.Itoa(diff)
			rankMap[strconv.Itoa(diff)] = true
		}
		passedPlayerNum += len(playerInTheTeam)
	}

	var rank []int
	for k := range rankMap {
		intK, _ := strconv.Atoi(k)
		rank = append(rank, intK)
	}

	sort.Ints(rank)
	for i := 0; i < len(rows[24])-2; i++ {
		userDataIndex := i + 2
		for j := 0; j < len(rank); j++ {
			if rows[23][userDataIndex] != strconv.Itoa(rank[j]) {
				continue
			} else {
				rows[24][userDataIndex] = strconv.Itoa(j + 1)
			}
		}
	}

	EntireScore := EntireScore{
		Rows: rows,
	}

	return &EntireScore, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
