package manor_test

import (
	"fmt"
	"testing"

	"time"

	manor "github.com/monchickey/manor-go/v2"
)

// go test -v file_util_test.go

func TestPathStatus(t *testing.T) {
	fmt.Println(manor.PathIsExist("/usr/aaa"), manor.PathIsExist("/usr/bin"))
	fmt.Println(manor.PathIsFile("/etc/hosts"), manor.PathIsFile("/usr/local/"))
	fmt.Println(manor.PathIsDir("/usr/local"), manor.PathIsDir("/usr/bin/python"))
}

func TestGobSerialize(t *testing.T) {
	var a = 3
	fmt.Println(manor.GobSerialize("/root/a.gob", a))
	a += 2
	fmt.Println(manor.GobSerialize("/root/a.gob", a))

	type Aaa struct {
		A int
		B string
	}
	as := Aaa{A: 3, B: "hello"}
	fmt.Println(manor.GobSerialize("/root/aaa.gob", as))

	type Bbb struct {
		A time.Time
		B time.Duration
	}
	bb := Bbb{A: time.Now(), B: 10 * time.Second}
	fmt.Println(manor.GobSerialize("/root/bbb.gob", bb))

	var b int
	err := manor.GobDeserialize("/root/a.gob", &b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(b)
	}
	aaa := Aaa{}
	err = manor.GobDeserialize("/root/aaa.gob", &aaa)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(aaa)
	}

	bbb := Bbb{}
	err = manor.GobDeserialize("/root/bbb.gob", &bbb)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(bbb)
	}
}
