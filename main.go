package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Item struct {
	Id   string    `json:"id"`
	DB   string    `json:"dB"`
	Date time.Time `json:"date"`
}
type System struct {
	Items map[string][]Item
}

var system = System{}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func main() {
	system.Init()
	l, err := net.Listen("tcp", ":8080")
	if nil != err {
		log.Println(err)
	}
	defer l.Close()

	go (func() {
		for {
			conn, err := l.Accept()
			defer conn.Close()
			if nil != err {
				log.Println(err)
				continue
			}
			go ConnHandler(conn)
		}
	})()

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	for {
		command := ""
		input := ""
		fmt.Print("> ")
		fmt.Scanf("%s %s\n", &command, &input)
		if command == "live" {
			live := true
			go (func() {
				for {
					if !live {
						break
					}
					tbl := table.New("ID", "dB", "Date")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
					for k, v := range system.Items {
						tbl.AddRow(k, v[len(v)-1].DB, v[len(v)-1].Date)
					}

					clear()
					tbl.Print()
					time.Sleep(1 * time.Second)
				}
			})()
			for {
				input := ""
				fmt.Print("> ")
				fmt.Scanln(&input)
				if input == "q" {
					live = false
					break
				}
			}
		} else if command == "his" {
			tbl := table.New("ID", "dB", "Date")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
			for k, v := range system.Items[input] {
				tbl.AddRow(k, v.DB, v.Date)
			}
			tbl.Print()
		} else if command == "clear" {
			clear()
		} else if command == "q" {
			break
		}
	}
}
func (s *System) Init() {
	s.Items = make(map[string][]Item)
}

func (s *System) AddItem(item Item) {
	s.Items[item.Id] = append(s.Items[item.Id], item)
}

func ConnHandler(conn net.Conn) {
	clear()
	fmt.Println("New Client")
	fmt.Print("> ")
	recvBuf := make([]byte, 4096)
	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}
		if 0 < n {
			data := recvBuf[:n]
			dd := Item{}
			json.Unmarshal(data, &dd)
			dd.Date = time.Now()
			system.AddItem(dd)
		}
	}
}
