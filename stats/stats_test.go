package stats

import (
	"fmt"
)

func ExampleStatsSaveToFile() {
	s := NewStats()
	s.TotalFailCdr = 10
	s.NatFailCount = 20
	s.FromTgFailCount[99] = 777
	err := s.SaveToFile("stats.gob")
	fmt.Println(err)
	//Output: <nil>
}

func ExampleLoadFromFile() {
	s, err := LoadFromFile("stats.gob")
	fmt.Println(s, err)
	//Output: {10 0 0 map[] map[99:777] map[] map[] map[] map[] 0 0 20 0 0 0} <nil>
}
