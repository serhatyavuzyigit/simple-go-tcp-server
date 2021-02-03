package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func main() {

	vpr := viper.New()
	vpr.SetConfigFile("config.yaml")
	vpr.ReadInConfig()
	portNumber := vpr.GetInt("port")
	message := vpr.GetString("message")
	strPortNumber := fmt.Sprintf(":%d", portNumber)

	dstream, err := net.Listen("tcp", strPortNumber)

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
		go handle(conn, message)

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
		switch operator {
		case "+":
			result := operand1 + operand2
			fmt.Println(message, result)
		case "-":
			result := operand1 - operand2
			fmt.Println(message, result)
		case "*":
			result := operand1 * operand2
			fmt.Println(message, result)
		case "/":
			result := operand1 / operand2
			fmt.Println(message, result)
		}
	}
}
