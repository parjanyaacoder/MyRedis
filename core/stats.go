package core

var KeyspaceStat [4]map[string]int

func UpdateDbStat(num int, metric string, value int) {
	KeyspaceStat[num][metric] = value
}
