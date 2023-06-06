package wiki

import (
	"github.com/KomeiDiSanXian/BlueArchive/bluearchive/web"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type Event struct {
	EventName   string
	Description string
	BeginAt     int64
	EndAt       int64
	PictureURL  string
}

type Events []*Event

var URL = "https://ba.gamekee.com/v1/wiki/index"

var Headers = map[string]string{
	"game-alias": "ba",
	"Connection": "close",
}

// func (es *Events) DeleteEvent(event *Event) {
// 	index := -1
// 	for i, e := range *es {
// 		if e == event {
// 			index = i
// 			break
// 		}
// 	}
// 	if index == -1 {
// 		return
// 	}
// 	(*es)[index] = (*es)[len(*es)-1]
// 	*es = (*es)[:len(*es)-1]
// }

func (es *Events) PrintEvent(layout string) []string {
	strs := make([]string, 0, len(*es))
	for _, event := range *es {
		h, m, s, isStarted := event.RemainingTime()
		if h < 0 {
			// es.DeleteEvent(event)
			continue
		}
		event.FixDescription()
		startTime := time.Unix(event.BeginAt, 0).Format(layout)
		endTime := time.Unix(event.EndAt, 0).Format(layout)
		nonstartfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n距离开始剩余时间: %d 小时 %d 分钟 %d 秒"
		startedfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n活动剩余时间: %d 小时 %d 分钟 %d 秒"
		if !isStarted {
			strs = append(strs, fmt.Sprintf(nonstartfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		} else {
			strs = append(strs, fmt.Sprintf(startedfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		}

	}
	return strs
}

func (e *Event) FixDescription() {
	// 删除 <br>
	e.Description = strings.ReplaceAll(e.Description, "<br>", "")
}

func (e *Event) RemainingTime() (hours, minutes, seconds int64, isStarted bool) {
	now := time.Now().Unix()
	beforeStart := e.BeginAt - now
	// 活动未开始
	if beforeStart > 0 {
		duration := time.Duration(beforeStart) * time.Second
		hours = int64(duration.Hours())
		minutes = int64(duration.Minutes()) % 60
		seconds = beforeStart % 60
		return
	}
	isStarted = true
	remain := e.EndAt - now
	// 活动结束
	if remain < 0 {
		return -1, -1, -1, isStarted
	}
	duration := time.Duration(remain) * time.Second
	hours = int64(duration.Hours())
	minutes = int64(duration.Minutes()) % 60
	seconds = remain % 60
	return
}

func (e *Event) DownloadPicture(downloadTo string) error {
	if e.PictureURL == "" {
		return errors.New("picture not found")
	}
	link := "http:" + e.PictureURL
	resp, err := web.MakeRequest(link, Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s := strings.Split(link, "/")
	name := s[len(s)-1]
	f, err := os.Create(downloadTo + "/" + name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func Request() (Events, error) {
	resp, err := web.MakeRequest(URL, Headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if gjson.GetBytes(jsonBytes, "code").Int() != 0 {
		return nil, errors.New("wiki response code not zero")
	}

	events := gjson.GetBytes(jsonBytes, "data.4.list")
	result := make(Events, 0, len(events.Array()))
	for _, value := range events.Array() {
		wiki := &Event{
			EventName:   value.Get("title").Str,
			Description: value.Get("description").Str,
			BeginAt:     value.Get("begin_at").Int(),
			EndAt:       value.Get("end_at").Int(),
			PictureURL:  value.Get("picture").Str,
		}
		result = append(result, wiki)
	}

	return result, nil
}
