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
	PortNumber int    `json:"portNumber"`
	Message    string `json:"message"`
}

var c config
var activeConnections []net.Conn
var isPortChanged bool

/*
	updates the config struct for current config values
*/
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

/*
	opens a new listener for current port number
*/
func openNewStream() {
	newStream, err := net.Listen("tcp", c.port)
	if err != nil {
		return
	}
	go handleConnections(newStream)
}

/*
	gets the data from active connection
*/
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

/*
	determining the result for given string array and message
	after that prints the result
*/
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

/*
	closes all active connections
*/
func closeConnections() {
	for i := 0; i < len(activeConnections); i++ {
		activeConnections[i].Close()
	}
	activeConnections = []net.Conn{}

}

/*
	handle the connections for given listener
*/
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

/*
	opens a connection for the first page load
*/
func openInitialConnection() {
	dstream, err := net.Listen("tcp", c.port)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstream.Close()
	handleConnections(dstream)

}

/*
	creating a channel for signal handling
*/
func createSignalChannel() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT:
					updateConfig()
					if isPortChanged {
						closeConnections()
						openNewStream()
					}
				}
			}
		}
	}()
}

/*
	apply needed operations for given request from frontend
*/
func apply(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	switch r.Method {
	case "GET":
		strPort := c.port
		ind := strings.LastIndex(strPort, ":")
		intPort, err := strconv.Atoi(string(strPort[ind+1:]))

		if err != nil {
			return
		}

		var data vueData
		data = vueData{PortNumber: intPort, Message: c.message}

		openInitialConnection()
		json.NewEncoder(w).Encode(data)
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var comingData vueData

		decoder.Decode(&comingData)
		viper.Set("message", comingData.Message)
		viper.Set("port", comingData.PortNumber)
		viper.WriteConfig()

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

/*
	main function of server.go
*/
func main() {

	c = config{}

	createSignalChannel()

	// setting viper fon config reading
	viper.SetConfigFile("config.yaml")
	updateConfig()

	http.HandleFunc("/", apply)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}

}
