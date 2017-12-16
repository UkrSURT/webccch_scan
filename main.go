package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

type Response struct {
	Code string //`json: "code"`
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var conn *websocket.Conn

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cmd := exec.Command("tcpdump", "-l", "-ilo", "-nXs0", "udp", "and", "port", "4729")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	} // start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Println("before output")
	// out, err := cmd.Output()
	// if err != nil{
	// panic(err)	// }
	// fmt.Println(string(out))	// read command's stdout line by line
	in := bufio.NewScanner(stdout)

	// striings := make([]string, 0, 100)

	//f, err := os.Create("out.txt")
	check(err)

	for in.Scan() {
		// log.Printf(in.Text()) // write each line to your log, or anything you need
		// striings = append(striings, in.Text())
		res := &Response{
			Code: in.Text(),
		}

		// jsonRes, _ := json.Marshal(res)

		if err = conn.WriteJSON(res); err != nil {
			return
		}

		//_, _ = f.WriteString(in.Text())
		// _, err := io.Copy(f, in.Text())
		// check(err)
	}

	// fmt.Print(striings)

	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	err = cmd.Wait()
	panic(err.Error())

}

func main() {

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
