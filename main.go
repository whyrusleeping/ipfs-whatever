package main

import (
	"fmt"
	"log"
	"time"

	randbo "github.com/dustin/randbo"
	api "github.com/ipfs/go-ipfs-api"

	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

func main() {
	r := randbo.New()

	rpath, err := fsrepo.BestKnownPath()
	if err != nil {
		log.Fatal(err)
	}

	apiaddr, err := fsrepo.APIAddr(rpath)
	if err != nil {
		log.Fatal(err)
	}

	sh := api.NewShell(apiaddr)

	count := 2000

	basedata := make([]byte, 100)
	r.Read(basedata)
	base := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

	cur, err := sh.PatchData(base, true, basedata)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cur)

	before := time.Now()
	for i := 0; i < count; i++ {
		out, err := sh.PatchLink(base, "next-link-in-chain", cur, false)
		if err != nil {
			fmt.Println("error: ", base, cur)
			log.Fatal(err)
		}

		cur = out
	}
	took := time.Now().Sub(before)
	fmt.Println(cur)

	opss := float64(count) / took.Seconds()
	fmt.Printf("%f op/s\n", opss)
}
