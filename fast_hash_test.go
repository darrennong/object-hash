package labs

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestObjectHash(t *testing.T) {
	data0 := obj{1, true, 3, 4, 5, 6, 7, 8}
	data1 := obj{1, true, 3, 4, 5, 6, 7, 8}
	data2 := obj{1, true, 3, 4, 5, 6, 7, 8}
	data3 := obj{1, true, 3, 4, 5, 6, 7, 8}
	data4 := obj{1, true, 3, 4, 5, 6, 7, 8}
	data5 := obj{1, true, 3, 4, 5, 6, 7, 8}
	dataA := testObj{data0, []int{9, 10, 11, 12, 13, 14, 15}, "test string", make(map[string]*obj)}
	dataB := testObj{data1, []int{9, 10, 11, 12, 13, 14, 15}, "test string", make(map[string]*obj)}
	dataA.K["data0"] = &data0
	dataA.K["data2"] = &data2
	dataA.K["data4"] = &data4
	dataB.K["data1"] = &data1
	dataB.K["data3"] = &data3
	dataB.K["data4"] = &data5
	fmt.Println("simple object:")
	b := time.Now()
	for i := 0; i < 10000; i++ {
		ObjectHash(&data0)
	}
	fmt.Printf("1000 operations take                   %6d microseconds.\n", time.Since(b).Microseconds())
	b = time.Now()
	for i := 0; i < 10000; i++ {
		reflect.DeepEqual(data0, data1)
	}
	fmt.Printf("1000 reflect.DeepEqual operations take %6d microseconds.\n", time.Since(b).Microseconds())
	b = time.Now()
	for i := 0; i < 10000; i++ {
		s := fmt.Sprintf("%v", data0)
		md5.New().Sum([]byte(s))
	}
	fmt.Printf("1000 normal hash operations take       %6d microseconds.\n", time.Since(b).Microseconds())
	fmt.Println("complex object:")
	b = time.Now()
	for i := 0; i < 10000; i++ {
		ObjectHash(&dataA)
	}
	fmt.Printf("1000 operations take                   %6d microseconds.\n", time.Since(b).Microseconds())
	b = time.Now()
	for i := 0; i < 10000; i++ {
		reflect.DeepEqual(dataA, dataB)
	}
	fmt.Printf("1000 reflect.DeepEqual operations take %6d microseconds.\n", time.Since(b).Microseconds())
	b = time.Now()
	for i := 0; i < 10000; i++ {
		s := fmt.Sprintf("%v", dataA)
		md5.New().Sum([]byte(s))
	}
	fmt.Printf("1000 normal hash operations take       %6d microseconds.\n", time.Since(b).Microseconds())
	fmt.Printf("Before change:\nHash of DataA:  %s\nHash of DataB:  %s\n", ObjectHash(&dataA), ObjectHash(&dataB))
	data0.A = 2
	data1.B = false
	fmt.Printf("After change:\nHash of DataA:  %s    Hash of DataB:  %s\n", ObjectHash(&dataA), ObjectHash(&dataB))
	data2.C++
	data3.D++
	fmt.Printf("After change:\nHash of DataA:  %s    Hash of DataB:  %s\n", ObjectHash(&dataA), ObjectHash(&dataB))
	data4.E++
	data5.F++
	fmt.Printf("After change:\nHash of DataA:  %s    Hash of DataB:  %s\n", ObjectHash(&dataA), ObjectHash(&dataB))
	dataA.H++
	dataB.I[3]++
	fmt.Printf("After change:\nHash of DataA:  %s    Hash of DataB:  %s\n", ObjectHash(&dataA), ObjectHash(&dataB))
}

type testObj struct {
	obj
	I []int
	J string
	K map[string]*obj
}
type obj struct {
	A int8
	B bool
	C int16
	D uint16
	E int32
	F uint32
	G int64
	H uint64
}
