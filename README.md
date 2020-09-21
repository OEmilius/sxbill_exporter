softX3000 raw CDR prometheus exporter


## Usage

sxbill_exporter.exe --help
Usage of sxbill_exporter.exe:
  -BatFilePath string
        path to BatFile (default "rcopy.bat")
  -CdrDir string
        path to CDRs must be same with robocopy.bat file (default "d:\tmp\bill")
  -IntervalMin int
        Time interval minutes between copy CDRs (default 15)
  -ListenAddress string
        listen address (default ":8080")
  -StatsFilePath string
        path to store Stats (default "stats.gob")
        
$sxbill_exporter.exe -ListenAddress=":7777"


## prometheus.yml
scrape_configs:
  - job_name: 'sxbill'
    scrape_interval: 5m
    static_configs:
    - targets: ['192.168.0.200:8080']
	
	
## PromQL examples
delta(TotalNormCdr[1h])    
delta(TotalMinutes[1h]) 
topk(5,rate(FromTgFailCount[1h]))
rate(LocalNormCount[1h])
topk(10, rate(ToTgSec[1h]))