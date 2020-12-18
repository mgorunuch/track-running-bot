package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	medal1, medal2, medal3, medalCommon = "🥇", "🥈", "🥉", "🎗"
)

func getPerDay(currentDistance, goal float64, leftDays int) float64 {
	res := (goal - currentDistance) / float64(leftDays)

	if res < 0 {
		return 0
	}

	return res
}

func getLeftToGoal(goal, currentDistance float64) float64 {
	res := goal - currentDistance

	if res < 0 {
		return 0
	}

	return res
}

func removeMessageDistanceMsg(removedDistance, currentDistance, goal float64, leftDays int) string {
	perDay := getPerDay(currentDistance, goal, leftDays)

	return fmt.Sprintf("Removed distance %.2f.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", removedDistance, currentDistance, getLeftToGoal(goal, currentDistance), perDay)
}

func registerMessageDistanceMsg(registeredDistance, currentDistance, goal float64, leftDays int) string {
	perDay := getPerDay(currentDistance, goal, leftDays)

	if getLeftToGoal(goal, currentDistance) <= 0 {
		return fmt.Sprintf("Registered distance %.2fkm.\n\n\U0001F973 <b>CONGRATULATIONS</b> 🎉\nYou achieved your goal", registeredDistance)
	}

	return fmt.Sprintf("Registered distance %.2fkm.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", registeredDistance, currentDistance, getLeftToGoal(goal, currentDistance), perDay)
}

func myMessageDistanceMsg(currentDistance, goal float64, leftDays int) string {
	perDay := getPerDay(currentDistance, goal, leftDays)
	base := ""

	if getLeftToGoal(goal, currentDistance) <= 0 {
		base = fmt.Sprintf("\U0001F973 <b>CONGRATULATIONS</b> 🎉\nYou achieved your goal")
	}

	return base + fmt.Sprintf("Your current distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", currentDistance, getLeftToGoal(goal, currentDistance), perDay)
}

type sortingDat struct {
	id   int
	dist float64
	name string
}

func statsMessageDistanceMsg(dist map[int]float64, nms map[int]string, goal float64, leftDays int) string {
	msgs := make([]string, 0, len(dist))
	sortingData := make([]sortingDat, 0, len(dist))
	for id, nm := range dist {
		sortingData = append(sortingData, sortingDat{
			dist: nm,
			name: nms[id],
		})
	}

	sort.Slice(sortingData, func(i, j int) bool {
		return sortingData[i].dist > sortingData[j].dist
	})

	for i, v := range sortingData {
		medal := ""
		switch i {
		case 0:
			medal = medal1
		case 1:
			medal = medal2
		case 2:
			medal = medal3
		default:
			medal = medalCommon
		}

		msgs = append(
			msgs,
			fmt.Sprintf("%v %v place: <a href=\"tg://user?id=%v\">%s</a> - <b>%.2fkm</b>", medal, i+1, v.id, v.name, v.dist),
		)
	}

	return strings.Join(msgs, "\n")
}
