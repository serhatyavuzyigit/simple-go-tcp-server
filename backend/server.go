package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	PortNumber int    `json:"portNumber"`
	Message    string `json:"message"`
}

var c config
var activeConnections []net.Conn
var isPortChanged bool

func updateConfig() {
	viper.ReadInConfig()
	portNumber := viper.GetInt("port")
	message := viper.GetString("message")
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
	viper.SetConfigFile("config.yaml")
	viper.ReadInConfig()

	//http.HandleFunc("/", apply)
	//http.ListenAndServe(":8092", nil)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	c = config{}
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT:
					updateConfig()
					if isPortChanged {
						closeConnections()
					}
				}
			}
		}
	}()

	do()

}

func apply(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var comingData vueData

	decoder.Decode(&comingData)
	viper.Set("message", comingData.Message)
	viper.Set("port", comingData.PortNumber)
	viper.WriteConfig()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func closeConnections() {
	for i := 0; i < len(activeConnections); i++ {
		activeConnections[i].Close()
	}
	activeConnections = []net.Conn{}
	newStream, err := net.Listen("tcp", c.port)
	if err != nil {
		return
	}
	go handleConnections(newStream)
}

func do() {
	updateConfig()
	//initializeTcpListener()
	dstream, err := net.Listen("tcp", c.port)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstream.Close()

	handleConnections(dstream)

}

func handleConnections(l net.Listener) {
	for {

		conn, err := l.Accept()
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
