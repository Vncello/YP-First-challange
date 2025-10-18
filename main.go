package main

import (
	"fmt"
	"net/http"
	"time"
	"io"
	"strings"
	"strconv"
)

const (
	TrashLoadAvarage = 30
	TrashPercentRamUsed = 0.8
	TrashPercentDiskUsed = 0.9
	TrashPercentNetUsed = 0.9
	BytesInMB = 1048576
)

var counter = 0

type StatsServe struct {
	LoadAvarage int
	RAMTotal int
	RAMUsed int
	DiskTotal int
	DiskUsed int
	NetTotal int
	NetUsed int
}

func (s *StatsServe) CheckStatistic() {

	if s.LoadAvarage > TrashLoadAvarage {
		fmt.Printf("Load Average is too high: %d\n",s.LoadAvarage)
	}

	if percentRAM:=(float32(s.RAMUsed)/float32(s.RAMTotal)); percentRAM > TrashPercentRamUsed {
		fmt.Printf("Memory usage too high: %.0f%\n", percentRAM*100)
	}

	if FreeDisk := -float32(s.DiskUsed)+float32(s.DiskTotal); FreeDisk/float32(s.DiskTotal) < (1-TrashPercentDiskUsed) {
		fmt.Printf("Free disk space is too low: %.0f Mb left\n", FreeDisk/BytesInMB)
	}

	if FreeNet := float32(s.NetTotal) - float32(s.NetUsed); FreeNet/float32(s.NetTotal) < (1-TrashPercentNetUsed) {
		fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", FreeNet/BytesInMB)
	}
}

func main() {
	for i:=0; i<5; i++ {


		response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")

		if err!=nil {
			fmt.Println(err)
			return
		}

		if response.StatusCode != http.StatusOK {
			counter++
			continue
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		//body := "32,21474836481934,1073741824,5497558138880,4398046511104,104857600,6291456"
		parts := strings.Split(string(body), ",")


		var my_slice = make([]int, 7)

		for i, part := range parts {

			value, err := strconv.Atoi(part)

			if err != nil{
				fmt.Println(err)
				counter++
				break
			}
			my_slice[i] = value
		}

		if counter == 3 {
			fmt.Println("Unable to fetch server statistic")
			return
		}

		Stats := StatsServe{my_slice[0], my_slice[1], my_slice[2], my_slice[3], my_slice[4], my_slice[5], my_slice[6]}

		//fmt.Println(Stats)
		Stats.CheckStatistic()


		time.Sleep(2 * time.Second)
	}
}
