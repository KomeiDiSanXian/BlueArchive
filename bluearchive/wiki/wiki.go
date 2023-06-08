// Package wiki 从 wiki (https://ba.gamekee.com/v1/wiki/index) 获取信息
package wiki

import (
	"errors"
	"fmt"
	"github.com/KomeiDiSanXian/BlueArchive/bluearchive/web"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// WikiData 存储公告和活动
type WikiData struct {
	Events        *Events
	Announcements *Announcements
}

// Event 用于存储活动信息
type Event struct {
	EventName   string
	Description string
	BeginAt     int64
	EndAt       int64
	PictureURL  *Picture
}

// Announcement 公告
type Announcement struct {
	Title       string
	Summary     string
	PictureURLs []*Picture
}

// Picture 存储图片信息
type Picture struct {
	Name string
	URL  string
}

// Events 活动信息切片
type Events []*Event

// Announcements 公告切片
type Announcements []*Announcement

// URL wiki URL
var URL = "https://ba.gamekee.com/v1/wiki/index"

// Headers wiki headers
var Headers = map[string]string{
	"game-alias": "ba",
	"Connection": "close",
}

// NewWikiData 创建空的WikiData，返回其指针
func NewWikiData() *WikiData {
	return &WikiData{}
}

// NewPictureByURL 使用url创建图片信息，并基于url命名
func NewPictureByURL(url string) *Picture {
	s := strings.Split(url, "/")
	return &Picture{
		Name: s[len(s)-1],
		URL:  url,
	}
}

// PrintEvent 输入时间格式，输出字符串切片
//
// 每个字符串都是格式化Event的 开始/结束/剩余时间 等
func (es *Events) PrintEvent(layout string) []string {
	strs := make([]string, 0, len(*es))
	for _, event := range *es {
		h, m, s, isStarted := event.remainingTime()
		if h < 0 {
			continue
		}
		event.fixDescription()
		startTime := time.Unix(event.BeginAt, 0).Format(layout)
		endTime := time.Unix(event.EndAt, 0).Format(layout)
		nonstartfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n距离开始剩余时间: %d 小时 %d 分钟 %d 秒\n"
		startedfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n活动剩余时间: %d 小时 %d 分钟 %d 秒\n"
		if !isStarted {
			strs = append(strs, fmt.Sprintf(nonstartfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		} else {
			strs = append(strs, fmt.Sprintf(startedfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		}

	}
	return strs
}

// PrintAnnouncements 返回每条公告的格式化信息
func (as *Announcements) PrintAnnouncements() []string {
	strs := make([]string, 0, len(*as))
	for _, a := range *as {
		str := fmt.Sprintf("%s\n\n%s...", a.Title, a.Summary)
		strs = append(strs, str)
	}
	return strs
}

// fixDescription 删除可能存在的 <br>
func (e *Event) fixDescription() {
	// 删除 <br>
	e.Description = strings.ReplaceAll(e.Description, "<br>", "")
}

// remainingTime 计算活动的剩余时间
//
// 如果活动未开始，输出的剩余时间是距离活动的开始时间
//
// 如果活动进行中，输出的剩余时间是距离活动的结束时间
//
// 剩余时间 比如 3661 会输出 1h1min1s 每个数字单独输出
func (e *Event) remainingTime() (hours, minutes, seconds int64, isStarted bool) {
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

// Download 将会下载 Picture到downloadTo 路径
func (p *Picture) Download(downloadTo string) error {
	if p.URL == "" {
		return errors.New("picture url not found")
	}
	link := "http:" + p.URL
	resp, err := web.MakeRequest(link, Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(downloadTo + "/" + p.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Request 请求获取wiki中的数据
func (w *WikiData) Request() error {
	resp, err := web.MakeRequest(URL, Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if gjson.GetBytes(jsonBytes, "code").Int() != 0 {
		return errors.New("wiki response code not zero")
	}

	w.Events = w.Events.getEvents(jsonBytes)
	w.Announcements = w.Announcements.getAnnouncements(jsonBytes)

	return nil
}

// getEvents 从jsonBytes 中获取活动信息
func (es *Events) getEvents(jsonBytes []byte) *Events {
	events := gjson.GetBytes(jsonBytes, "data.4.list").Array()
	result := make(Events, 0, len(events))
	for _, value := range events {
		picurl := value.Get("picture").Str
		event := &Event{
			EventName:   value.Get("title").Str,
			Description: value.Get("description").Str,
			BeginAt:     value.Get("begin_at").Int(),
			EndAt:       value.Get("end_at").Int(),
			PictureURL:  NewPictureByURL(picurl),
		}
		result = append(result, event)
	}
	return &result
}

// getAnnouncements 从jsonBytes 中获取公告信息
func (as *Announcements) getAnnouncements(jsonBytes []byte) *Announcements {
	announcements := gjson.GetBytes(jsonBytes, "data.2.list").Array()
	result := make(Announcements, 0, len(announcements))
	for _, value := range announcements {
		announcement := &Announcement{
			Title:   value.Get("title").Str,
			Summary: value.Get("summary").Str,
		}
		picurls := value.Get("thumb").Str
		if picurls != "" {
			urls := strings.Split(picurls, ",")
			pics := make([]*Picture, 0, len(urls))
			for _, url := range urls {
				pics = append(pics, NewPictureByURL(url))
			}
			announcement.PictureURLs = pics
		}
		result = append(result, announcement)
	}
	return &result
}

// 使用URL的最后一个path 给图片命名
func (p *Picture) NamedByURL() {
	s := strings.Split(p.URL, "/")
	p.Name = s[len(s)-1]
}
