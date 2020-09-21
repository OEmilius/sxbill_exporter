package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sxbill_exporter/call"
	"sxbill_exporter/controldir"
	"sxbill_exporter/rcopy"
	"sxbill_exporter/stats"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	ListenAddress string
	BatFile       string
	CdrDir        string
	StatsFile     string
	IntervalMin   int
}

var cfg = Config{}
var STATS *stats.Stats

var (
	TotalFailCdr      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "TotalFailCdr"})
	TotalNormCdr      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "TotalNormCdr"})
	TotalMinutes      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "TotalMinutes"})
	TotalSec          = prometheus.NewGauge(prometheus.GaugeOpts{Name: "TotalSec"})
	FromTgNormCount   = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "FromTgNormCount"}, []string{"tg"})
	FromTgFailCount   = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "FromTgFailCount"}, []string{"tg"})
	FromTgSec         = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "FromTgSec"}, []string{"tg"})
	ToTgNormCount     = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "ToTgNormCount"}, []string{"tg"})
	ToTgFailCount     = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "ToTgFailCount"}, []string{"tg"})
	ToTgSec           = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "ToTgSec"}, []string{"tg"})
	LocalFailCount    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "LocalFailCount"})
	LocalNormCount    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "LocalNormCount"})
	NatFailCount      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "NatFailCount"})
	NatNormCount      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "NatNormCount"})
	InternatFailCount = prometheus.NewGauge(prometheus.GaugeOpts{Name: "InternatFailCount"})
	InternatNormCount = prometheus.NewGauge(prometheus.GaugeOpts{Name: "InternatNormCount"})
)

func init() {
	flag.StringVar(&cfg.ListenAddress, "ListenAddress", ":8080", "listen address")
	flag.StringVar(&cfg.BatFile, "BatFilePath", "rcopy.bat", "path to BatFile")
	flag.StringVar(&cfg.CdrDir, "CdrDir", "d:\tmp\bill", "path to CDRs must be same with robocopy.bat file")
	flag.StringVar(&cfg.StatsFile, "StatsFilePath", "stats.gob", "path to store Stats")
	flag.IntVar(&cfg.IntervalMin, "IntervalMin", 15, "Time interval minutes between copy CDRs")
	flag.Parse()
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(TotalFailCdr)
	prometheus.MustRegister(TotalNormCdr)
	prometheus.MustRegister(TotalMinutes)
	prometheus.MustRegister(TotalSec)
	prometheus.MustRegister(FromTgNormCount)
	prometheus.MustRegister(FromTgFailCount)
	prometheus.MustRegister(FromTgSec)
	prometheus.MustRegister(ToTgNormCount)
	prometheus.MustRegister(ToTgFailCount)
	prometheus.MustRegister(ToTgSec)
	prometheus.MustRegister(LocalFailCount)
	prometheus.MustRegister(LocalNormCount)
	prometheus.MustRegister(NatFailCount)
	prometheus.MustRegister(NatNormCount)
	prometheus.MustRegister(InternatFailCount)
	prometheus.MustRegister(InternatNormCount)
}

func main() {
	log.Println("start")
	log.Printf("%#v\r\n", cfg)
	// cfg.BatFile = `rcopy.bat`
	// cfg.CdrDir = `d:\sxbill_exporter\bill`
	// cfg.StatsFile = "stats.gob"
	// cfg.IntervalMin = 60
	var err error
	STATS, err = stats.LoadFromFile(cfg.StatsFile)
	if err != nil {
		log.Println(err)
		STATS = stats.NewStats()
	}
	fmt.Println("lastProcessedFile=", STATS.Get_lastProccedFile())
	time.Sleep(30 * time.Second)
	//setPromMetrics(STATS)
	go StartTicker(cfg.IntervalMin, STATS)
	defer STATS.SaveToFile(cfg.StatsFile)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(cfg.ListenAddress, nil)
	fmt.Scanln()

}

func StartTicker(min int, s *stats.Stats) {
	tick := time.NewTicker(time.Duration(min) * time.Minute)
	for {
		log.Println(rcopy.Execute(cfg.BatFile))
		flist, err := controldir.WalkDirectory(cfg.CdrDir)
		if err != nil {
			panic(err)
		}
		npFiles := controldir.GetNotProcessedFiles(s.Get_lastProccedFile(), flist)
		log.Println("found new files=", npFiles)
		for _, fn := range npFiles {
			failCdr, normCdr := call.ProcessFile(fn)
			s.AppendStats(failCdr, normCdr)
			s.Set_lastProccedFile(fn)
			log.Println("in for lastProccedFile=", s.Get_lastProccedFile())
			setPromMetrics(s)
		}
		// if len(npFiles) != 0 {
		// 	s.Set_lastProccedFile(npFiles[len(npFiles)-1])
		// }
		<-tick.C
	}
}

func setPromMetrics(s *stats.Stats) {
	//s2 := STATS.Copy()
	TotalFailCdr.Set(float64(s.TotalFailCdr))
	TotalNormCdr.Set(float64(s.TotalNormCdr))
	TotalMinutes.Set(s.TotalMinutes)
	TotalSec.Set(s.TotalSec)
	for k, v := range s.FromTgNormCount {
		FromTgNormCount.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(float64(v))
	}
	for k, v := range s.FromTgFailCount {
		FromTgFailCount.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(float64(v))
	}
	for k, v := range s.FromTgSec {
		FromTgSec.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(v)
	}
	for k, v := range s.ToTgNormCount {
		ToTgNormCount.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(float64(v))
	}
	for k, v := range s.ToTgFailCount {
		ToTgFailCount.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(float64(v))
	}
	for k, v := range s.ToTgSec {
		ToTgSec.With(prometheus.Labels{"tg": fmt.Sprintf("%d", k)}).Set(v)
	}
	LocalFailCount.Set(float64(s.LocalFailCount))
	LocalNormCount.Set(float64(s.LocalNormCount))
	NatFailCount.Set(float64(s.NatFailCount))
	NatNormCount.Set(float64(s.NatNormCount))
	InternatFailCount.Set(float64(s.InternatFailCount))
	InternatNormCount.Set(float64(s.InternatNormCount))
}
