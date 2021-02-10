package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/viper"
)

type config struct {
	port    string
	message string
}

type vueData struct {
	portNumber int    `json:"portNumber"`
	message    string `json:"message"`
}

var activeConnections []net.Conn
var isPortChanged bool

func (c *config) init() {
	vpr := viper.New()
	vpr.SetConfigFile("config.yaml")
	vpr.ReadInConfig()
	portNumber := vpr.GetInt("port")
	message := vpr.GetString("message")
	strPortNumber := fmt.Sprintf(":%d", portNumber)

	if c.port != "" && c.port != strPortNumber {
		isPortChanged = true
	} else {
		isPortChanged = false
	}

	c.message = message
	c.port = strPortNumber
}

func main() {

	http.HandleFunc("/", apply)
	http.ListenAndServe(":8090", nil)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP)

	c := &config{}
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT:
					c.init()
					if isPortChanged {
						closeConnections()
					}
				}
			}
		}
	}()

	do(c)

}

func apply(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	fmt.Println("girdi hocam")
	decoder := json.NewDecoder(r.Body)
	var reqValues vueData
	fmt.Println(decoder)
	decoder.Decode(&reqValues)
	fmt.Println(reqValues.message)
	fmt.Println(reqValues.portNumber)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func closeConnections() {
	for i := 0; i < len(activeConnections); i++ {
		activeConnections[i].Close()
	}
	activeConnections = []net.Conn{}
}

func do(c *config) {
	c.init()
	dstream, err := net.Listen("tcp", c.port)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer dstream.Close()
	for {

		conn, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handle(conn, c.message)

	}
}

func handle(conn net.Conn, message string) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		activeConnections = append(activeConnections, conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		s := strings.Split(data, " ")
		if len(s) == 3 {
			printResult(s, message)
		}

	}

}

func printResult(s []string, message string) {
	operand1, err1 := strconv.Atoi(s[0])
	operator := s[1]
	operand2, err2 := strconv.Atoi(strings.TrimSuffix(s[2], "\n"))
	if err1 == nil && err2 == nil {
		var result int = 0
		switch operator {
		case "+":
			result = operand1 + operand2
		case "-":
			result = operand1 - operand2
		case "*":
			result = operand1 * operand2
		case "/":
			result = operand1 / operand2
		}
		fmt.Println(message, result)
	}
}
