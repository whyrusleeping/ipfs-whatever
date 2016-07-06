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

func checkAddLink(count int) (float64, float64, error) {
	var times []float64

	r := randbo.New()
	basedata := make([]byte, 100)
	r.Read(basedata)
	base := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

	prev := base
	for i := 0; i < count; i++ {
		start := time.Now()

		cur, err := sh.PatchLink(base, "FIRST", prev, false)
		if err != nil {
			fmt.Println("error: ", err)
			return 0, 0, err
		}

		for j := 0; j < 200; j++ {
			out, err := sh.PatchLink(cur, fmt.Sprintf("link-%d", j), base, false)
			if err != nil {
				fmt.Println("error: ", base, cur)
				return 0, 0, err
			}

			cur = out
		}
		took := float64(time.Now().Sub(start)) / 200 / float64(time.Second)
		times = append(times, 1/took)
	}
	avg, stdev := timeStats(times)
	return avg, stdev, nil
}

func checkAddFile(size int) (float64, float64, error) {
	trials := 10
	buf := new(bytes.Buffer)
	var times []float64

	for i := 0; i < trials; i++ {
		io.CopyN(buf, randbo.New(), int64(size))

		start := time.Now()
		_, err := sh.Add(buf)
		if err != nil {
			return 0, 0, err
		}
		took := float64(time.Now().Sub(start)) / float64(time.Second)
		times = append(times, took)
	}

	av, stdev := timeStats(times)
	return av, stdev, nil
}

func timeStats(ts []float64) (float64, float64) {
	var average float64
	for _, d := range ts {
		average += d / float64(len(ts))
	}

	var stdevsum float64
	for _, d := range ts {
		v := average - d
		stdevsum += (v * v)
	}

	stdev := math.Sqrt(stdevsum / float64(len(ts)))

	return average, stdev
}

type IpfsBenchStats struct {
	PatchOpsPerSec  float64
	Add10MBTime     float64
	Add10MBStdev    float64
	DirAddOpsPerSec float64
	DirAddOpsStdev  float64
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

	diradd, diraddstd, err := checkAddLink(100)
	if err != nil {
		log.Fatal(err)
	}
	stats.DirAddOpsPerSec = diradd
	stats.DirAddOpsStdev = diraddstd

	json.NewEncoder(os.Stdout).Encode(stats)
}
