/*
Configure ssh server with
	  ssh-keygen -A	            # generate host keys
	  /usr/bin/sshd -p 2222     # Start daemon.
	  kill sshd                 # Stop daemon.

RFC 4251, 4252, 4253, 4254
Architecture             https://www.rfc-editor.org/rfc/rfc4251
Authentication           https://www.rfc-editor.org/rfc/rfc4252
Transport Layer Protocol https://www.rfc-editor.org/rfc/rfc4253
Connection Protocol      https://www.rfc-editor.org/rfc/rfc4254
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type sshServerInfo struct {
	identification string
}

type SshConnection struct {
	conn    net.Conn
	scanner *bufio.Scanner
}

func SshConnect(addr string) SshConnection {
	var err error
	sshConnection := SshConnection{}
	sshConnection.conn, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TCP dial failed: ", err.Error())
		os.Exit(1)
	}
	sshConnection.scanner = bufio.NewScanner(sshConnection.conn)
	return sshConnection
}

func SshIdentificationHandshake(ctx SshConnection, clientId string) string {
	fmt.Fprintf(ctx.conn, clientId+"\r\n")

	serverId := ""
	for ctx.scanner.Scan() {
		if strings.HasPrefix(ctx.scanner.Text(), "SSH-") {
			serverId = ctx.scanner.Text()
			break
		}
	}
	if err := ctx.scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "TCP read failed: ", err.Error())
		os.Exit(1)
	}
	return serverId
}

func SshCloseConnection(ctx SshConnection) {
	ctx.conn.Close()
}

func main() {
	addr := "localhost:2222"
	ctx := SshConnect(addr)

	clientIdentification := "SSH-2.0-MySSH_0.1"
	serverIdentification := SshIdentificationHandshake(ctx, clientIdentification)
	println("Client Identification: ", clientIdentification)
	println("Server Identification: ", serverIdentification)

	SshCloseConnection(ctx)
}
