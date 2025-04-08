package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func PrintStats() {
	for {
		debug.FreeOSMemory()
		time.Sleep(5 * time.Second)
		/*time.Sleep(5 * time.Second)
		var oldStats runtime.MemStats
		runtime.ReadMemStats(&oldStats)
		s := time.Now()



		var newStats runtime.MemStats
		runtime.ReadMemStats(&newStats)

		log.Printf(
			"free mem took %s BEFORE: %4dmb(%4dmb/%5dmb gc:%v) NOW: %4dmb(%4dmb/%5dmb gc:%v)",
			time.Since(s).String(),
			bToMb(oldStats.Alloc), bToMb(oldStats.HeapAlloc), bToMb(oldStats.Sys), oldStats.NumGC,
			bToMb(newStats.Alloc), bToMb(newStats.HeapAlloc), bToMb(newStats.Sys), newStats.NumGC,
		)*/
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func targeter(typ string, endpoint string, key string) func(tgt *vegeta.Target) error {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}
		*tgt = vegeta.Target{
			Method: "POST",
			URL:    endpoint,
			Header: http.Header{
				"content-type": []string{"application/json"},
			},
			Body: []byte(GenerateData(typ, key)),
		}
		return nil
	}
}
func main() {
	go PrintStats()

	if os.Getenv("OTLP_ENDPOINT") == "" {
		log.Fatalf("set all env carefully.")
	}

	rate := vegeta.Rate{
		Freq: IntEnv("RATE", 0),
		Per:  time.Second,
	}

	//duration := time.Duration(IntEnv("DURATION", 60)) * time.Second
	target := targeter(os.Getenv("OTLP_TYPE"), os.Getenv("OTLP_ENDPOINT")+"/v1/"+os.Getenv("OTLP_TYPE"), os.Getenv("ACCOUNT_KEY"))

	log.Printf("starting attack")

	log.Printf("RATE= %d", uint64(IntEnv("RATE", 0)))
	log.Printf("DURATION= %d", uint64(IntEnv("DURATION", 0)))
	//for {
	//start := time.Now()
	attacker := vegeta.NewAttacker(
		vegeta.MaxWorkers(uint64(IntEnv("RATE", 0))),
		vegeta.Workers(uint64(IntEnv("RATE", 0))),
		vegeta.MaxBody(1000000),
		//vegeta.DNSCaching(120*time.Second),
		vegeta.KeepAlive(true),
		//vegeta.H2C(true),
		//vegeta.H2C(true),
		vegeta.HTTP2(true),
		vegeta.Timeout(60*time.Second),
	)

	var metrics vegeta.Metrics
	var lock sync.Mutex
	go func() {
		for res := range attacker.Attack(target, rate, 0, "Big Bang!") {
			lock.Lock()
			metrics.Add(res)
			lock.Unlock()
			if res.Error != "" {
				log.Printf("response code: %v error %s body:%v", res.Code, res.Error, string(res.Body))
			} /* else if res.Latency >= 20*time.Second {
				log.Printf("response %v resp %s", res.Code, res.Latency.String())
			}*/
			res.End()
		}
	}()

	var pointPerRequest float64 = 0
	var sizePerRequest float64 = 0
	switch os.Getenv("OTLP_TYPE") {
	case "logs":
		logs, _ := strconv.Atoi(os.Getenv("LOG_COUNT"))
		pointPerRequest = float64(logs)
		sizePerRequest = (float64(logs) / 1000) * (0.33 * 1024 * 1024)
		break
	}

	for {
		lock.Lock()
		mtrs := metrics
		mtrs.Close()
		metrics = vegeta.Metrics{}
		lock.Unlock()

		fmt.Printf("TOTAL %6s (%2.2f%%) (THROUGHPUT: %6s/s) (POINTS: %6s/s) (SIZE: %6s/s) RANGE (%6s to %6s) 	(MEAN: %6s) P9x (%6s  %5s  %6s)\n",
			NumFormat(int64(mtrs.Requests)),
			mtrs.Success*100,
			NumFormat(int64(mtrs.Throughput)),
			NumFormat(int64(mtrs.Throughput*pointPerRequest)),
			ByteFormat(int64(mtrs.Throughput*sizePerRequest)),

			RoundDuration(mtrs.Latencies.Min),
			RoundDuration(mtrs.Latencies.Max), RoundDuration(mtrs.Latencies.Mean),
			RoundDuration(mtrs.Latencies.P90), RoundDuration(mtrs.Latencies.P95), RoundDuration(mtrs.Latencies.P99),
			//RoundDuration(time.Since(start)),
		)
		time.Sleep(10 * time.Second)
		debug.FreeOSMemory()
	}
	//time.Sleep(5 * time.Second)
	//log.Printf("starting again.")
	//debug.FreeOSMemory()
	//}
}

func NumFormat(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d", b)
	}
	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c",
		float64(b)/float64(div), "kMBTPE"[exp])
}
func ByteFormat(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
func RoundDuration(d time.Duration) time.Duration {
	switch {
	case d > time.Hour:
		d = d.Round(time.Hour / 10)
		break
	case d > time.Minute:
		d = d.Round(time.Minute / 10)
		break
	case d > time.Second:
		d = d.Round(time.Second / 100)
		break
	case d > time.Millisecond:
		d = d.Round(time.Millisecond)
		break
	case d > time.Microsecond:
		d = d.Round(time.Microsecond)
		break
	}
	return d
}

func IntEnv(name string, def int) int {
	if os.Getenv(name) == "" {
		return def
	}
	int, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		log.Fatalf("error parsing %s as number %v", name, err)
	}
	return int
}
