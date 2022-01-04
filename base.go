package cctv

import "github.com/gocolly/colly/v2"

func newSpider(async bool) *colly.Collector {
	c := colly.NewCollector()
	c.Async = async
	c.IgnoreRobotsTxt = false
	return c
}

type errors []error

func (es *errors) add(err error) {
	if err == nil {
		return
	}
	*es = append(*es, err)
}

func (es errors) first() error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

func callWithOutErr(err error, cbs ...func() error) error {
	if err != nil {
		return err
	}
	for _, cb := range cbs {
		if e := cb(); e != nil {
			return e
		}
	}
	return nil
}
