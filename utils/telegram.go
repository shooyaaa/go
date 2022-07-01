package utils

import (
	"fmt"
	"github.com/shooyaaa/log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/core/library"
	core "github.com/shooyaaa/core/network"
)

type IM interface {
	Send(message IMMessage)
	DefaultChannel() interface{}
}

type IMMessage interface {
	GetText() *string
	GetImage() *string
	GetReceiver() interface{}
}

type TextMessage struct {
	Text  string
	To    interface{}
	Title interface{}
}

func (tm *TextMessage) GetImage() *string {
	return nil
}

func (tm *TextMessage) GetReceiver() interface{} {
	return &tm.To
}

func (tm *TextMessage) GetText() *string {
	return &tm.Text
}

type ImageMessage struct {
	Image string
	To    interface{}
	Desc  string
}

func (tm *ImageMessage) GetImage() *string {
	return &tm.Image
}

func (tm *ImageMessage) GetReceiver() interface{} {
	return &tm.To
}

func (tm *ImageMessage) GetText() *string {
	return &tm.Desc
}

type Job interface {
	Start() IMMessage
}

const apiUrl = "https://api.telegram.org/%v/%v"
const (
	CmdSendText  = "sendMessage"
	CmdSendImage = "sendPhoto"
)

type Telegram struct {
	jobs []Job
}

func (t *Telegram) AddJobs(job Job) {
	if t.jobs == nil {
		t.jobs = []Job{}
	}
	t.jobs = append(t.jobs, job)
}

func (t *Telegram) SendImage(file string, id interface{}) (interface{}, error) {
	fmt.Println("using telegram to send image ")
	params := make(map[string]interface{})
	params["chat_id"] = id
	params["photo"] = file
	http := core.Http{}
	ret, err := http.HttpPostForm(t.url(CmdSendImage), params)
	if err != nil {
		fmt.Printf("Error response from Telegram: %v, err: %v", ret, err)
	}
	return ret, nil
}

func (t *Telegram) DefaultChanel() interface{} {
	return config.TelegramUserID
}

func (t *Telegram) Send(message IMMessage) {
	if message == nil {
		return
	}
	to := message.GetReceiver()
	_, ret := to.(int)
	if ret == false {
		to = t.DefaultChanel()
	}
	if message.GetImage() != nil {
		t.SendImage(*message.GetImage(), to)
	} else {
		t.SendText(*message.GetText(), to)
	}
}

func (t *Telegram) SendText(text string, id interface{}) (interface{}, error) {
	fmt.Println("using telegram to send text")
	params := make(map[string]interface{})
	params["chat_id"] = id
	params["text"] = text
	params["parse_mode"] = "HTML"
	http := core.Http{}
	ret, err := http.HttpPostJson(t.url(CmdSendText), params)
	if err != nil {
		fmt.Printf("Error response from Telegram: %v, err: %v", ret, err)
	}
	return ret, nil
}

func (t *Telegram) url(cmd string) string {
	return fmt.Sprintf(apiUrl, config.TelegramBotToken, cmd)
}

func (t *Telegram) Do() {
	for _, job := range t.jobs {
		msg := job.Start()
		t.Send(msg)
	}
}

func (t *Telegram) Cron() {
	fmt.Println("Telegram cron")
	du := FundRequest(161028)
	http := core.Http{}

	ret, _ := http.HttpGet(du.Url, nil, du.Headers)
	if buff, ok := ret.(string); ok {
		t.SendImage("#"+buff, config.TelegramUserID)
	} else {
		log.DebugF("get fund data failed with %v", ret)
	}

	now := time.Now()
	date := now.Format("2006-01-02T15:04")
	now = now.AddDate(-1, 0, 0)
	start := now.Format("2006-01-02T15:04")
	cd := CryptoCurrency("ETH", start, date)
	cd1 := CryptoCurrency("BTC", start, date)
	cd2 := CryptoCurrency("DOGE", "2020-05-20T00:00", date)

	charDataList := []*library.ChartData{}
	if cd != nil {
		charDataList = append(charDataList, ToOneHundred(cd))
	}
	if cd1 != nil {
		charDataList = append(charDataList, ToOneHundred(cd1))
	}
	if cd2 != nil {
		charDataList = append(charDataList, ToOneHundred(cd2))
	}
	ret = library.Line(charDataList)
	t.SendImage("@"+ret.(string), config.TelegramUserID)
	os.Remove(ret.(string))
}

func FromTwoArray(entries interface{}) *library.ChartData {
	x := make([]time.Time, 0)
	y := make([]float64, 0)
	for _, item := range entries.([]interface{}) {
		epoch := int64((item.([]interface{})[0]).(float64))
		unix := time.Unix(int64(epoch/1000), 0)
		x = append(x, unix)
		y = append(y, (item.([]interface{})[1]).(float64))
	}
	return &library.ChartData{
		X: x,
		Y: y,
	}
}

func ToOneHundred(cd *library.ChartData) *library.ChartData {
	temp := []float64{}
	temp = append(temp, cd.Y...)
	sort.Float64s(temp)
	min, max := temp[0], temp[len(temp)-1]
	ratio := (max - min) / 100.0
	for v := range cd.Y {
		cd.Y[v] = (cd.Y[v] - min) / ratio
	}
	return cd
}

func CryptoCurrency(name string, start string, end string) *library.ChartData {
	url := fmt.Sprintf(CryptoCurrencyUrl, name, start, end)
	du := &DataUrl{
		Url: url,
	}
	http := core.Http{}
	ret, _ := http.HttpGet(du.Url, nil, du.Headers)
	if info, ok := ret.(map[string]interface{}); ok {
		entries, _ := info["data"].(map[string]interface{})["entries"]
		cd := FromTwoArray(entries)
		cd.Name = name
		return cd
	}
	return nil
}

func DownloadHtmlAsDoc(url string) (*goquery.Document, error) {
	ret, err := http.Get(url)
	if ret == nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(ret.Body)
	return doc, err
}

const url = "https://www.smzdm.com/"

type Smzdm struct {
}

func (sm *Smzdm) Start() IMMessage {
	doc, err := DownloadHtmlAsDoc(url)
	msg := TextMessage{}
	if err != nil {
		msg.Text = err.Error()
	}
	score := 0
	doc.Find("#feed-main-list li").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Find(".feed-block .z-feed-img img").Attr("src")
		subTitle := s.Find(".feed-block .z-feed-content .z-highlight a")
		title := s.Find(".feed-block-title a")
		vote := s.Find(".unvoted-wrap span")
		date := s.Find(".feed-block-extras")
		link := s.Find(".z-btn-red")
		platform := date.Find("a")
		if title.Nodes == nil {
			return
		}
		if vote.Nodes == nil {
			return
		}
		href := ""
		for _, attr := range link.Nodes[0].Attr {
			if attr.Key == "href" {
				href = attr.Val
			}
		}
		good, _ := strconv.Atoi(vote.Nodes[0].FirstChild.Data)
		bad, _ := strconv.Atoi(vote.Nodes[1].FirstChild.Data)
		if good+bad < 50 {
			return
		}
		ratio := math.Pow(float64(good/(good+bad)), float64(good+bad))
		if int(ratio) > score {
			fmt.Printf("Review %d: %s, %v, %v\n", i, img, date.Nodes[0].FirstChild.Data, platform)
			msg.Text = fmt.Sprintf("<a href='%v'>%s:%s, %s</a>, %s, 值:%v, 不值:%v", href, date.Nodes[0].FirstChild.Data, title.Nodes[0].FirstChild.Data,
				subTitle.Nodes[0].FirstChild.Data, img, good, bad)
			score = int(ratio)
		}
	})
	if len(msg.Text) > 0 {
		return &msg
	} else {
		return nil
	}
}
