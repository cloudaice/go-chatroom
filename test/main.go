package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	speed int64 = 0
)

func work(ch chan struct{}, idx int) {
	defer func() {
		<-ch
	}()
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Println(err)
		return
	}
	buf := bufio.NewWriter(conn)
	go func() {
		io.Copy(ioutil.Discard, conn)
	}()

	for {
		atomic.AddInt64(&speed, 1)
		time.Sleep(200 * time.Millisecond)
		room := strconv.Itoa(((idx-1)/100+1)*10 - 9 + rand.Int()&10)
		msg := room + " " + strconv.Itoa(rand.Int()) + "\n"
		_, err := buf.Write([]byte(msg))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	go func() {
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			fmt.Println(atomic.LoadInt64(&speed))
			atomic.StoreInt64(&speed, 0)
		}
	}()
	idx := 1
	ch := make(chan struct{}, 40000)
	for {
		ch <- struct{}{}
		go work(ch, idx)
		time.Sleep(2 * time.Millisecond)
		idx++
	}
}
