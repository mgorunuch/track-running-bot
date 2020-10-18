package main

import (
	"fmt"
	"strings"
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

func statsMessageDistanceMsg(dist map[int]float64, nms map[int]string, goal float64, leftDays int) string {
	msgs := make([]string, 0, len(dist))
	for id, nm := range dist {
		perDay := (goal - nm) / float64(leftDays)

		msgs = append(
			msgs,
			fmt.Sprintf("Distance stats for user <a href=\"tg://user?id=%v\">%s</a> is:\n<i>Current distance</i>: <b>%.2fkm</b>\n<i>Left to goal</i>: <b>%.2fkm</b>\n<i>Per day</i>: <b>%.2fkm</b>", id, nms[id], nm, goal-nm, perDay),
		)
	}

	return strings.Join(msgs, "\n\n")
}
