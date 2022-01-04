package cctv

import (
	"fmt"
	"testing"
)

func TestJingJi(t *testing.T) {
	datas, err := JingJi("ARTIcyLY0Zw5mHtzaPRIx3FX220104")
	fmt.Println(err, len(datas), datas[0])
}
