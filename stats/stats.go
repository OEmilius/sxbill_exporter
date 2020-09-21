package stats

import (
	"encoding/gob"
	"log"
	"os"

	//	"sort"
	"strings"
	"sxbill_exporter/call"
	"sync"
)

type Stats struct {
	mux               sync.Mutex
	LastProccedFile   string
	TotalFailCdr      int64
	TotalNormCdr      int64
	TotalMinutes      float64
	TotalSec          float64
	FromTgNormCount   map[uint16]int64
	FromTgFailCount   map[uint16]int64
	FromTgSec         map[uint16]float64
	ToTgNormCount     map[uint16]int64
	ToTgFailCount     map[uint16]int64
	ToTgSec           map[uint16]float64
	LocalFailCount    int64
	LocalNormCount    int64
	NatFailCount      int64
	NatNormCount      int64
	InternatFailCount int64
	InternatNormCount int64
}

func NewStats() *Stats {
	s := Stats{}
	s.FromTgNormCount = make(map[uint16]int64)
	s.FromTgFailCount = make(map[uint16]int64)
	s.FromTgSec = make(map[uint16]float64)
	s.ToTgNormCount = make(map[uint16]int64)
	s.ToTgFailCount = make(map[uint16]int64)
	s.ToTgSec = make(map[uint16]float64)
	return &s
}

func (s *Stats) Get_lastProccedFile() string {
	s.mux.Lock()
	last := s.LastProccedFile
	defer s.mux.Unlock()
	return last
}

func (s *Stats) Set_lastProccedFile(fname string) {
	s.mux.Lock()
	s.LastProccedFile = fname
	s.mux.Unlock()
}

func (s *Stats) Copy() Stats {
	s.mux.Lock()
	s2 := *s
	defer s.mux.Unlock()
	return s2
}

func (s *Stats) appendStats(failCdr, normCdr []call.Call) {
	log.Println("start Stats) appendStats")
	s.TotalFailCdr = s.TotalFailCdr + int64(len(failCdr))
	s.TotalNormCdr = s.TotalNormCdr + int64(len(normCdr))
	for _, c := range normCdr {
		s.TotalMinutes = s.TotalMinutes + float64(c.Conversation_time)/100.0/60.0
		s.TotalSec = s.TotalSec + float64(c.Conversation_time)/100.0
		s.FromTgNormCount[c.Trunk_group_in] = s.FromTgNormCount[c.Trunk_group_in] + 1
		s.FromTgSec[c.Trunk_group_in] = s.FromTgSec[c.Trunk_group_in] + (float64(c.Conversation_time) / 100.0)
		s.ToTgNormCount[c.Trunk_group_out] = s.ToTgNormCount[c.Trunk_group_out] + 1
		s.ToTgSec[c.Trunk_group_out] = s.ToTgSec[c.Trunk_group_out] + (float64(c.Conversation_time) / 100.0)
		if strings.HasPrefix(c.Called_number, "495") || strings.HasPrefix(c.Called_number, "499") {
			s.LocalNormCount = s.LocalNormCount + 1
		} else if strings.HasPrefix(c.Called_number, "810") {
			s.InternatNormCount = s.InternatNormCount + 1
		} else {
			s.NatNormCount = s.NatNormCount + 1
		}

	}
	for _, c := range failCdr {
		s.FromTgFailCount[c.Trunk_group_in] = s.FromTgFailCount[c.Trunk_group_in] + 1
		s.ToTgFailCount[c.Trunk_group_out] = s.ToTgFailCount[c.Trunk_group_out] + 1
		if strings.HasPrefix(c.Called_number, "495") || strings.HasPrefix(c.Called_number, "499") {
			s.LocalFailCount = s.LocalFailCount + 1
		} else if strings.HasPrefix(c.Called_number, "810") {
			s.InternatFailCount = s.InternatFailCount + 1
		} else {
			s.NatFailCount = s.NatFailCount + 1
		}
	}
	log.Println("end Stats) appendStats")
}

func (s *Stats) AppendStats(failCdr, normCdr []call.Call) {
	s.mux.Lock()
	s.appendStats(failCdr, normCdr)
	s.mux.Unlock()
}

func (s *Stats) SaveToFile(path string) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	log.Println("start Stats) SaveToFile")
	log.Println("last proeccess file=", s.LastProccedFile)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	dataEncoder := gob.NewEncoder(file)
	err = dataEncoder.Encode(s)
	defer file.Close()
	if err != nil {
		return err
	}
	log.Println("lastProccedFile=", s.LastProccedFile)
	log.Println("end Stats) SaveToFile")
	return nil
}

func LoadFromFile(path string) (*Stats, error) {
	log.Println("start stats LoadFromFile ", path)
	s := NewStats()
	file, err := os.Open(path)
	if err != nil {
		log.Println("end stats LoadFromFile ", err)
		return s, err
	}
	defer file.Close()
	dataDecoder := gob.NewDecoder(file)
	err = dataDecoder.Decode(&s)
	if err != nil {
		log.Println("end stats LoadFromFile ", err)
		return s, err
	}
	log.Println("end stats LoadFromFile successfully")
	log.Println("last process file=", s.LastProccedFile)
	return s, nil
}

// func GetStats(failCdr, normCdr []call.Call) Stats {
// 	s := NewStats()
// 	s.TotalFailCdr = len(failCdr)
// 	s.TotalNormCdr = len(normCdr)
// 	for _, c := range normCdr {
// 		s.FromTgNormCount[c.Trunk_group_in] = s.FromTgNormCount[c.Trunk_group_in] + 1
// 		s.FromTgSec[c.Trunk_group_in] = s.FromTgSec[c.Trunk_group_in] + (float64(c.Conversation_time) / 100)
// 		s.ToTgNormCount[c.Trunk_group_out] = s.ToTgNormCount[c.Trunk_group_out] + 100
// 		s.ToTgSec[c.Trunk_group_out] = s.ToTgSec[c.Trunk_group_out] + (float64(c.Conversation_time) / 100)
// 	}
// 	return s
// }

// type TgCount struct {
// 	TgNumb uint16
// 	Count  int
// }

// type TgCountList []TgCount

// func SortByValueInt(m map[uint16]int) TgCountList {
// 	pl := make(TgCountList, len(m))
// 	i := 0
// 	for k, v := range m {
// 		pl[i] = TgCount{k, v}
// 		i++
// 	}
// 	sort.Sort(sort.Reverse(pl))
// 	return pl
// }

// func (p TgCountList) Len() int           { return len(p) }
// func (p TgCountList) Less(i, j int) bool { return p[i].Count < p[j].Count }
// func (p TgCountList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// /*           */

// type TgSec struct {
// 	TgNumb  uint16
// 	Seconds float64
// }

// type TgSecList []TgSec

// func SortByValueFloat(m map[uint16]float64) TgSecList {
// 	pl := make(TgSecList, len(m))
// 	i := 0
// 	for k, v := range m {
// 		pl[i] = TgSec{k, v}
// 		i++
// 	}
// 	sort.Sort(sort.Reverse(pl))
// 	return pl
// }

// func (p TgSecList) Len() int           { return len(p) }
// func (p TgSecList) Less(i, j int) bool { return p[i].Seconds < p[j].Seconds }
// func (p TgSecList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
