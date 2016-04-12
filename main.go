package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	randbo "github.com/dustin/randbo"
	api "github.com/ipfs/go-ipfs-api"

	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

var sh *api.Shell

func checkPatchOpsPerSec(count int) (float64, error) {
	r := randbo.New()
	basedata := make([]byte, 100)
	r.Read(basedata)
	base := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

	cur, err := sh.PatchData(base, true, basedata)
	if err != nil {
		log.Fatal(err)
	}

	before := time.Now()
	for i := 0; i < count; i++ {
		out, err := sh.PatchLink(base, "next-link-in-chain", cur, false)
		if err != nil {
			fmt.Println("error: ", base, cur)
			return 0, err
		}

		cur = out
	}
	took := time.Now().Sub(before)

	return float64(count) / took.Seconds(), nil
}

func checkAddFile(size int) (time.Duration, time.Duration, error) {
	trials := 10
	buf := new(bytes.Buffer)
	var times []time.Duration

	for i := 0; i < trials; i++ {
		io.CopyN(buf, randbo.New(), int64(size))

		start := time.Now()
		_, err := sh.Add(buf)
		if err != nil {
			return 0, 0, err
		}
		took := time.Now().Sub(start)
		times = append(times, took)
	}

	av, stdev := timeStats(times)
	return av, stdev, nil
}

func timeStats(ts []time.Duration) (time.Duration, time.Duration) {
	var sum time.Duration
	for _, d := range ts {
		sum += d
	}

	average := sum / time.Duration(len(ts))

	var stdevsum time.Duration
	for _, d := range ts {
		v := average - d
		stdevsum += (v * v)
	}

	stdev := time.Duration(math.Sqrt(float64(stdevsum / time.Duration(len(ts)))))

	return average, stdev
}

type IpfsBenchStats struct {
	PatchOpsPerSec float64
	Add10MBTime    time.Duration
	Add10MBStdev   time.Duration
}

func main() {
	rpath, err := fsrepo.BestKnownPath()
	if err != nil {
		log.Fatal(err)
	}

	apiaddr, err := fsrepo.APIAddr(rpath)
	if err != nil {
		log.Fatal(err)
	}

	sh = api.NewShell(apiaddr)

	stats := new(IpfsBenchStats)

	count, err := checkPatchOpsPerSec(2000)
	if err != nil {
		log.Fatal(err)
	}
	stats.PatchOpsPerSec = count

	av, stdev, err := checkAddFile(10 * 1024 * 1024)
	if err != nil {
		log.Fatal(err)
	}
	stats.Add10MBTime = av
	stats.Add10MBStdev = stdev

	json.NewEncoder(os.Stdout).Encode(stats)
}
