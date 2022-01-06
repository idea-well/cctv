package cctv

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type jingJiRes struct {
	Data struct {
		List JingJiDatas `json:"list"`
	} `json:"data"`
}

type JingJiData struct {
	Id      string `json:"id"`         // 新闻ID
	Desc    string `json:"brief"`      // 新闻简介
	Title   string `json:"title"`      // 新闻标题
	Content string `json:"content"`    // 新闻内容
	SrcUrl  string `json:"url"`        // 新闻链接
	PubTime string `json:"focus_date"` // 发布时间
}

func (j *JingJiData) PubTimeFormat(layout string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", j.PubTime)
	return t.Format(layout)
}

type JingJiDatas []*JingJiData

func (ds JingJiDatas) sliceFrom(id string) JingJiDatas {
	for i := range ds {
		if ds[i].Id == id {
			return ds[0:i]
		}
	}
	return ds
}

func (ds JingJiDatas) fetchContent() error {
	var es = make(errors, 0)
	var lock = make(chan struct{}, 5)
	spider := newSpider(true)
	spider.OnHTML("#content_area", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ds[i].Content, _ = e.DOM.Html()
	})
	spider.OnError(func(resp *colly.Response, err error) {
		es.add(fmt.Errorf("status code: %d, err = %#v", resp.StatusCode, err))
	})
	spider.OnResponse(func(_ *colly.Response) { <-lock })
	for i := range ds {
		lock <- struct{}{}
		frame := fmt.Sprintf("#%d", i)
		es.add(spider.Visit(ds[i].SrcUrl + frame))
	}
	spider.Wait() // wait done
	return es.first()
}

var jingJiApi = "https://news.cctv.com/2019/07/gaiban/cmsdatainterface/page/economy_zixun_1.jsonp"

// JingJi 经济频道
func JingJi(lastId string) (JingJiDatas, error) {
	datas, err := doFetchJingJi()
	if err == nil && len(datas) > 0 {
		datas = datas.sliceFrom(lastId)
	}
	return datas, callWithOutErr(err, datas.fetchContent)
}

func doFetchJingJi() (JingJiDatas, error) {
	var res = new(jingJiRes)
	resp, err := http.Get(jingJiApi)
	return res.Data.List, callWithOutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithOutErr(err, func() error {
			return json.Unmarshal(bts[14:len(bts)-1], res)
		})
	})
}
