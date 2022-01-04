package cctv

import (
	"fmt"
	"testing"
)

func TestXWLB(t *testing.T) {
	datas, err := XWLB("20220103")
	fmt.Println(err, len(datas), datas[0])
}
