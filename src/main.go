package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Payload struct {
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji"`
	Username    string       `json:"username"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text string `json:"text"`
}

const icon = ":hourglass_flowing_sand:"
const userName = "TimeLimit"

func myHandler() error {
	slackURL := os.Getenv("URL")
	born, err := time.Parse("2006/01/02", os.Getenv("BORN"))
	if err != nil {
		return err
	}

	lifetime, err := strconv.Atoi(os.Getenv("LIFETIME"))
	if err != nil {
		return err
	}

	now := time.Now().Round(time.Second)
	dyingDay := born.AddDate(lifetime, 0, 0)
	nextYear := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, time.Local)

	yearLeft := nextYear.Sub(now)
	yearCurrent := now.Sub(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local))
	lifeLeft := dyingDay.Sub(now)
	lifeCurrent := now.Sub(born)

	lifeLeftDays := int(lifeLeft.Seconds()/86400)
	lifeLeftHours := int(lifeLeft.Hours())
	lifeCurrentHours := int(lifeCurrent.Hours())
	lifePercent := float64(lifeCurrentHours) / float64(lifeCurrentHours + lifeLeftHours) * 100

	yearLeftDays := int(yearLeft.Seconds()/86400)
	yearLeftHours := int(yearLeft.Hours())
	yearCurrentHours := int(yearCurrent.Hours())
	yearPercent := float64(yearCurrentHours) / float64(yearCurrentHours + yearLeftHours) * 100

	t := fmt.Sprintf("Your life left: *%d* days\n", lifeLeftDays)
	t += fmt.Sprintf("Your life left: *%d* hours\n", lifeLeftHours)
	t += fmt.Sprintf("Your life is current: *%.2f* %%\n\n", lifePercent)
	t += fmt.Sprintf("This year left: *%d* days\n", yearLeftDays)
	t += fmt.Sprintf("This year left: *%d* hours\n", yearLeftHours)
	t += fmt.Sprintf("This year is current: *%.2f* %%\n", yearPercent)

	p := Payload{
		Text:        fmt.Sprintln(now),
		IconEmoji:   icon,
		Username:    userName,
		Attachments: []Attachment{{Text: t}},
	}

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(myHandler)
}
