package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

// https://datatracker.ietf.org/doc/html/rfc1928

func Auth(client net.Conn) (err error) {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return fmt.Errorf("invalid version: %d", ver)
	}

	// read NMETHODS
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	// X'00' NO AUTHENTICATION REQUIRED
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp: " + err.Error())
	}

	return nil
}

// Socks5Connect receive reqeust
// The SOCKS request is formed as follows:
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
func Request(client net.Conn) (string, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return "", errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return "", errors.New("invalid ver/cmd")
	}

	addr := ""
	// ATYP   address type of following address
	// o  IP V4 address: X'01'
	// o  DOMAINNAME: X'03'
	// o  IP V6 address: X'04'
	switch atyp {
	case 1: // ipv4
		n, err = io.ReadFull(client, buf[:4])
		if n != 4 {
			return "", errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3: // domain name
		// the address field contains a fully-qualified domain name. The first
		// octet of the address field contains the number of octets of name that
		// follow, there is no terminating NUL octet.

		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return "", errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	case 4: //ip v6
		return "", errors.New("IPv6: no supported yet")
	default:
		return "", errors.New("invalid atyp")
	}

	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return "", errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	_, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return "", errors.New("write rsp: " + err.Error())
	}

	return destAddrPort, nil
}
