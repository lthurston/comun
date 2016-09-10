package main

import (
	cli "gopkg.in/urfave/cli.v2"
	"os"
	"sync"
	"log"
	"bufio"
	"fmt"
)

type counter struct {
	sync.Mutex
	strings    map[string]int
	max int
}

func newCounter() *counter {
	counter := &counter{
		strings: make(map[string]int),
		max: 0,
	}
	return counter
}

func (c *counter) add(string string) (int, int) {
	var i int
	var ok bool
	if i, ok = c.strings[string]; ok == true {
		c.strings[string] = i + 1
	} else {
		c.strings[string] = 1
	}
	if c.max < i + 1 {
		c.max++
	}
	return c.max, c.strings[string]
}

func main() {
	app := &cli.App {
		Name: "comun",
		Usage: "identifies common lines from multiple files, returns the first found",
		Action: run,
	}
	app.Run(os.Args)
}

var wg sync.WaitGroup
var count *counter
var resultChan chan string

func run(c *cli.Context) error {
	wg = sync.WaitGroup{}
	count = newCounter()
	target := c.NArg()
	resultChan = make(chan string, 9)

	for i := 0; i < target; i++ {
		filename := c.Args().Get(i)
		wg.Add(1)
		go func() error {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				count.Lock()
				text := scanner.Text()
				max, cnt := count.add(text)
				count.Unlock()
				if max >= target {
					if max == cnt {
						resultChan <- text
					}
					wg.Done()
					return nil
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			wg.Done()
			return nil
		}()
	}

	wg.Wait()

	var result string
	select {
	case result = <-resultChan:
	default:
	}

	close(resultChan)
	if len(result) > 0 {
		fmt.Println(result)
		os.Exit(0)
	}
	fmt.Println("No common lines found")
	os.Exit(1)

	return nil
}
