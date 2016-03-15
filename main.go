package main

import (
	"encoding/base32"
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

	api, err := fsrepo.APIAddr(rpath)
	if err != nil {
		log.Fatal(err)
	}

	sh := api.NewShell(api)

	count := 1000

	var names []string
	buf := make([]byte, 10)
	for i := 0; i < count; i++ {
		r.Read(buf)

		names = append(names, base32.StdEncoding.EncodeToString(buf))
	}

	before := time.Now()
	base := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
	cur := base
	for i := 0; i < count; i++ {
		out, err := sh.PatchLink(base, names[i], cur, false)
		if err != nil {
			fmt.Println("error: ", base, names[i], cur)
			fmt.Println(err)
			return
		}

		cur = out
	}
	took := time.Now().Sub(before)

	opss := float64(count) / took.Seconds()
	fmt.Printf("%f op/s\n", opss)
}
