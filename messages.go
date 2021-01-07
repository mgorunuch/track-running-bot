package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
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

func getLeftDays() (int, int) {
	res := goalEnd.Sub(time.Now())

	hrs := res.Hours()
	if hrs < 0 {
		return 0, 0
	}

	days := math.Floor(hrs / 24)
	hrs = hrs - (days * 24)

	return int(days), int(hrs)
}

func removeMessageDistanceMsg(removedDistance, currentDistance, goal float64, leftDays int) string {
	perDay := getPerDay(currentDistance, goal, leftDays)

	return fmt.Sprintf("Removed distance %.2f.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", removedDistance, currentDistance, getLeftToGoal(goal, currentDistance), perDay)
}

func registerMessageDistanceMsg(registeredDistance, currentDistance, goal float64, leftDays int) string {
	if getLeftToGoal(goal, currentDistance) <= 0 {
		return fmt.Sprintf("Registered distance %.2fkm.\n\n\U0001F973 <b>CONGRATULATIONS</b> 🎉\nYou achieved your goal", registeredDistance)
	}

	return fmt.Sprintf("Registered distance %.2fkm.\n", registeredDistance) + myMessageDistanceMsg(currentDistance, goal, leftDays)
}

func myMessageDistanceMsg(currentDistance, goal float64, leftDays int) string {
	perDay := getPerDay(currentDistance, goal, leftDays)
	base := ""

	days, hrs := getLeftDays()

	if getLeftToGoal(goal, currentDistance) <= 0 {
		base = fmt.Sprintf("\U0001F973 <b>CONGRATULATIONS</b> 🎉\nYou achieved your goal\n\n")
	}

	leftToGoal := getLeftToGoal(goal, currentDistance)
	avgPerDay := currentDistance / (totalDays - float64(days))

	currentPercent := currentDistance / (goal / 100)

	return base + fmt.Sprintf("Your current distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\nLeft time: %d day %d hours\n<i>Per day</i>: <b>%.5fkm</b>\n<i>Avg per day:</i>: %.3fkm\n<i>Ready perc</i>: %.3f%", currentDistance, leftToGoal, days, hrs, perDay, avgPerDay, currentPercent)
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
