package cctv

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
	"time"
)

var xwlbUrlFormat = "https://tv.cctv.com/lm/xwlb/day/%s.shtml"

type XwlbData struct {
	Title   string
	Content string
	SrcUrl  string
	PubTime time.Time
}

type XwlbDatas []*XwlbData

func (ds XwlbDatas) fetchContent() error {
	var es = make(errors, 0)
	spider := newSpider(true)
	spider.OnHTML("#about_txt .cnt_bd", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ds[i].Content, _ = e.DOM.Html()
	})
	spider.OnHTML("#chbox01 .text_box_02", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ss := strings.Split(e.ChildTexts("p")[1], "：")
		ds[i].PubTime, _ = time.Parse("2006年01月02日 15:04", ss[1])
	})
	spider.OnError(func(resp *colly.Response, err error) {
		es.add(fmt.Errorf("status code: %d, err = %#v", resp.StatusCode, err))
	})
	for i, d := range ds {
		frame := fmt.Sprintf("#%d", i)
		es.add(spider.Visit(d.SrcUrl + frame))
	}
	spider.Wait() // wait done
	return es.first()
}

// XWLB 每日新闻联播
func XWLB(date string) (XwlbDatas, error) {
	var es = make(errors, 0)
	var ds = make(XwlbDatas, 0)
	spider := newSpider(false)
	spider.OnHTML("ul li a", func(e *colly.HTMLElement) {
		title := e.DOM.Find(".title").Text()
		ds = append(ds, &XwlbData{
			Title:  strings.TrimLeft(title, "[视频]"),
			SrcUrl: e.Attr("href"),
		})
	})
	spider.OnError(func(resp *colly.Response, err error) {
		es.add(fmt.Errorf("status code: %d, err = %#v", resp.StatusCode, err))
	})
	es.add(spider.Visit(fmt.Sprintf(xwlbUrlFormat, date)))
	return ds[1:], callWithOutErr(es.first(), ds[1:].fetchContent)
}
