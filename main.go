package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TrashLoadAvarage     = 30
	TrashPercentRAMUsed  = 0.8
	TrashPercentDiskUsed = 0.9
	TrashPercentNetUsed  = 0.9
	BytesInMB            = 1048576
)

var counter = 0

type StatsServe struct {
	LoadAvarage int
	RAMTotal    int
	RAMUsed     int
	DiskTotal   int
	DiskUsed    int
	NetTotal    int
	NetUsed     int
}

func (s *StatsServe) CheckStatistic() {

	if s.LoadAvarage > TrashLoadAvarage {
		fmt.Printf("Load Average is too high: %d\n", s.LoadAvarage)
	}

	if percentRAM := (float32(s.RAMUsed) / float32(s.RAMTotal)); percentRAM > TrashPercentRAMUsed {
		fmt.Printf("Memory usage too high: %.0f%%\n", percentRAM*100)
	}

	if FreeDisk := -float32(s.DiskUsed) + float32(s.DiskTotal); FreeDisk/float32(s.DiskTotal) < (1 - TrashPercentDiskUsed) {
		fmt.Printf("Free disk space is too low: %d Mb left\n", int(FreeDisk/BytesInMB))
	}

	if FreeNet := float32(s.NetTotal) - float32(s.NetUsed); FreeNet/float32(s.NetTotal) < (1 - TrashPercentNetUsed) {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(FreeNet/1000000))
	}
}

func main() {
	for i := 0; i < 100; i++ {

		response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")

		if err != nil {
			fmt.Println(err)
			return
		}

		if response.StatusCode != http.StatusOK {
			counter++
			continue
		} else {
			counter = 0
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		//body := "9,4707314194,1909918678,371916994759,100338391010,7917216201,1162050106"
		parts := strings.Split(string(body), ",")

		var mySlice = make([]int, 7)

		for i, part := range parts {

			value, err := strconv.Atoi(part)

			if err != nil {
				fmt.Println(err)
				counter++
				break
			}
			mySlice[i] = value
		}

		if counter == 3 {
			fmt.Println("Unable to fetch server statistic")
			return
		}

		Stats := StatsServe{mySlice[0], mySlice[1], mySlice[2], mySlice[3], mySlice[4], mySlice[5], mySlice[6]}

		//fmt.Println(Stats)
		Stats.CheckStatistic()

		time.Sleep(2 * time.Second)
	}
}
