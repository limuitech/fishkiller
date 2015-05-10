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
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const ADD = "http://jfh.10086yux.com/submit.asp"
const REG = "location.href='(down.asp)'"

var procs int
var timeout int
var duration int

func init() {
	flag.IntVar(&procs, "proc", 1, "Start n processes.")
	flag.IntVar(&timeout, "timeout", 10, "Set timeout")
	flag.IntVar(&duration, "freq", 1000, "Set access frequency (ms)")
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan int, 1024)
	client := new(http.Client)
	client.Timeout = time.Second * time.Duration(timeout)
	for i := 0; i < procs; i++ {
		go whoop(client, ch, i, duration)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
	counter := 0
	for {
		select {
		case <-ch:
			counter++
			fmt.Printf("已轰炸%d次\r", counter)
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
		vals.Set("idType", "1")
		vals.Set(
			"cnName",
			base64.StdEncoding.EncodeToString(
				[]byte(
					strconv.Itoa(
						rand.Int(),
					),
				)[0:3],
			),
		)
		vals.Set(
			"sec_val",
			base64.StdEncoding.EncodeToString(
				[]byte(
					strconv.Itoa(
						rand.Int(),
					),
				)[0:10],
			),
		)
		vals.Set(
			"idcard",
			fmt.Sprintf(
				"%04d%04d%05d",
				strconv.Itoa(
					rand.Int()%10000,
				),
				strconv.Itoa(
					rand.Int()%10000,
				),
				strconv.Itoa(
					rand.Int()%100000,
				),
			),
		)
		vals.Set(
			"idcard1",
			fmt.Sprintf(
				"%06d",
				strconv.Itoa(
					rand.Int()%1000000,
				),
			),
		)
		vals.Set(
			"idNo1",
			fmt.Sprintf(
				"%06d%08d%04d",
				strconv.Itoa(
					rand.Int()%1000000,
				),
				strconv.Itoa(
					rand.Int()%10000000,
				),
				strconv.Itoa(
					rand.Int()%10000,
				),
			),
		)
		vals.Set(
			"shouji",
			fmt.Sprintf(
				"1%05d%05d",
				strconv.Itoa(
					rand.Int()%100000,
				),
				strconv.Itoa(
					rand.Int()%100000,
				),
			),
		)
		vals.Set("ssName", fmt.Sprintf("%03d", strconv.Itoa(rand.Int()%1000)))
		vals.Set("sja", "01")
		vals.Set(
			"sja",
			//fmt.Sprintf("%04d", strconv.Itoa(rand.Int()%10000)),
			"2018",
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
