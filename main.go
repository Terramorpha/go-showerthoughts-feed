package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/*
reddit feed template


go see https://www.reddit.com/dev/api

*/

var (
	Duration     time.Duration
	BufferLength int
)

func init() {
	flag.DurationVar(&Duration, "d", 7*time.Second, "usage")
	flag.IntVar(&BufferLength, "b", 100, "buffer length")
	flag.Parse()
}

const subreddit = "showerthoughts"

func main() {

	clientObject := http.Client{}
	//prepares the request
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/"+subreddit+"/new.json?sort=new&limit=1", nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	//req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Set("User-Agent", "redditFeedBot")
	//reddit doesn't accept invalid user-agents
	for {
		resp, err := clientObject.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(body)
		resp.Body.Close()
		//fmt.Println(body)

		var baseFeed redditFeed

		//fmt.Println(body)

		err = json.Unmarshal(bodyByte, &baseFeed)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//for _, v := range baseFeed.Data.Children {
		//	fmt.Println("")
		//	feed(v.Data.Title, 50)
		//}
		if len(baseFeed.Data.Children) != 1 {
			fmt.Fprintln(os.Stderr, baseFeed)
			time.Sleep(Duration)
		}
		feed("     "+baseFeed.Data.Children[0].Data.Title+"     ", BufferLength, Duration)
	}

}

type redditFeed struct {
	Kind string
	Data data
}

type data struct {
	Modhash          string
	WhitelistSstatus string
	Children         []item
}

type item struct {
	Kind string
	Data data2
}

type data2 struct {
	Domain string
	Title  string
}

func feed(text string, bufferlength int, timeToWait time.Duration) { //tuimetowait est en secondes
	//e monde ceci est un long message de plus de vingt caractères!!!
	buffer := make([]rune, bufferlength)

	iterationNum := len(text) - bufferlength
	if iterationNum <= 0 {
		fmt.Println(text)
		time.Sleep(timeToWait)
		return
	}
	//splitTimeToWait := (timeToWait / time.Duration(iterations)) - 2*time.Second
	EdgeTime := timeToWait / 5
	//fmt.Println("edgeTime:", EdgeTime)
	splitTime := (timeToWait - (EdgeTime * 2)) / time.Duration(iterationNum)
	EdgeTime += timeToWait % splitTime

	//fmt.Println("pre-edgeTime:", EdgeTime)
	//fmt.Println("split:", splitTime)
	//fmt.Println("edgeTime:", EdgeTime)
	for pointer := 0; pointer < iterationNum; pointer++ {
		for i := range buffer { //
			buffer[i] = []rune(text)[i+pointer] // on remplie le buffer
		} //
		fmt.Println(string(buffer[:]))
		if pointer == 0 {
			time.Sleep(EdgeTime)
		} else {
			time.Sleep(splitTime)
		}
	}
	time.Sleep(EdgeTime)
}
