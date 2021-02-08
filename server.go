package main

import (
	"bufio"
	"fmt"
	"net"
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

func (c *config) init() {
	vpr := viper.New()
	vpr.SetConfigFile("config.yaml")
	vpr.ReadInConfig()
	portNumber := vpr.GetInt("port")
	message := vpr.GetString("message")
	strPortNumber := fmt.Sprintf(":%d", portNumber)
	c.message = message
	c.port = strPortNumber
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	c := &config{}
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT:
					c.init()
				}
			}
		}
	}()

	do(c)

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
		if err != nil {
			fmt.Println(err)
			return
		}

		s := strings.Split(data, " ")
		if len(s) == 3 {
			printResult(s, message)
		}

	}
	conn.Close()
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
