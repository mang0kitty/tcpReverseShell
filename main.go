package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mang0kitty/tcpReverseShell/rsh"

	"github.com/urfave/cli/v2"
)

/**
S: start listening
S: wait for a client

-- s.Listen()

C: connect to the server
C: send "hello" to the server

-- c.Connect()
-- c.Send("hello")

S: receive "hello" from the client
S: send "hello" back to the client

-- d := s.Receive()
-- s.Send(d)

C: close
S: close
-- c.Close()
-- s.Close()
*/

func startServer(c *cli.Context) error {
	addr := c.String("addr")
	fmt.Printf("Listening on %s\n", addr)
	tcp, err := rsh.NewTCPServer(addr)
	if err != nil {
		return err
	}

	defer tcp.Close()

	data, err := tcp.Receive()
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return tcp.Send(data)
}

func startClient(c *cli.Context) error {
	addr := c.String("addr")
	fmt.Printf("Client connecting to %s\n", addr)
	tcp, err := rsh.NewTCPClient(addr)
	if err != nil {
		return err
	}
	ps := rsh.Powershell()
	err = ps.Execute("tasklist")

	if err != nil {
		log.Fatal(err)
	}

	// cmd := exec.Command("powershell.exe", "-NoExit", "-Command", "-")
	// stdin, _ := cmd.StdinPipe()
	// stdout, _ := cmd.StdoutPipe()
	// cmd.Start()
	defer tcp.Close()

	err = tcp.Send([]byte("this is a test"))
	if err != nil {
		return err
	}

	response, err := tcp.Receive()
	if err != nil {
		return err
	}

	fmt.Printf(string(response))

	return nil
}

func main() {
	app := &cli.App{
		Name: "rsh",

		Commands: []*cli.Command{
			&cli.Command{
				Name:  "server",
				Usage: "Runs an rsh server.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: "0.0.0.0:5555",
						Usage: "The address to listen on.",
					},
					&cli.StringFlag{
						Name:  "run",
						Value: "",
						Usage: "Specify a program to run on clients which connect to the server.",
					},
				},
				Action: startServer,
			},
			&cli.Command{
				Name:  "client",
				Usage: "Runs an rsh client.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: "127.0.0.1:5555",
						Usage: "the server to connect to",
					},
				},
				Action: startClient,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
