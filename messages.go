package main

import (
	"fmt"
	"math"
	"strings"
)

func removeMessageDistanceMsg(removedDistance, currentDistance, goal float64, leftDays int) string {
	perDay := math.Round((goal - currentDistance) / float64(leftDays))

	return fmt.Sprintf("Removed distance %.2f.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%v</b>", removedDistance, currentDistance, goal-currentDistance, perDay)
}

func registerMessageDistanceMsg(registeredDistance, currentDistance, goal float64, leftDays int) string {
	perDay := math.Round((goal - currentDistance) / float64(leftDays))

	return fmt.Sprintf("Registered distance %.2fkm.\nCurrent distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%v</b>", registeredDistance, currentDistance, goal-currentDistance, perDay)
}

func myMessageDistanceMsg(currentDistance, goal float64, leftDays int) string {
	perDay := math.Round((goal - currentDistance) / float64(leftDays))

	return fmt.Sprintf("Your current distance is: %.2fkm\nLeft to goal: <b>%.2fkm</b>\n<i>Per day</i>: <b>%v</b>", currentDistance, goal-currentDistance, perDay)
}

func statsMessageDistanceMsg(dist map[int]float64, nms map[int]string, goal float64, leftDays int) string {
	msgs := make([]string, 0, len(dist))
	for id, nm := range dist {
		perDay := math.Round((goal - nm) / float64(leftDays))

		msgs = append(
			msgs,
			fmt.Sprintf("Distance stats for user <a href=\"tg://user?id=%v\">%s</a> is:\n<i>Current distance</i>: <b>%.2fkm</b>\n<i>Left to goal</i>: <b>%.2fkm</b>\n<i>Per day</i>: <b>%v</b>", id, nms[id], nm, goal-nm, perDay),
		)
	}

	return strings.Join(msgs, "\n\n")
}
