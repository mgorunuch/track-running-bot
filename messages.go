package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	medal1, medal2, medal3, medalCommon = "🥇", "🥈", "🥉", "🎗"
)

func removeMessageDistanceMsg(removedDistance, currentDistance, goal float64, leftDays int) string {
	perDay := (goal - currentDistance) / float64(leftDays)

	return fmt.Sprintf("Removed distance %.2f.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", removedDistance, currentDistance, goal-currentDistance, perDay)
}

func registerMessageDistanceMsg(registeredDistance, currentDistance, goal float64, leftDays int) string {
	perDay := (goal - currentDistance) / float64(leftDays)

	return fmt.Sprintf("Registered distance %.2fkm.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", registeredDistance, currentDistance, goal-currentDistance, perDay)
}

func myMessageDistanceMsg(currentDistance, goal float64, leftDays int) string {
	perDay := (goal - currentDistance) / float64(leftDays)

	return fmt.Sprintf("Your current distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", currentDistance, goal-currentDistance, perDay)
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
