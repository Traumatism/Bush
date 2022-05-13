package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Output struct {
	Host string                 `json:"host"`
	Port uint16                 `json:"port"`
	Data map[string]interface{} `json:"data"`
	Date string                 `json:"date"`
}

type byteReaderWrap struct {
	reader io.Reader
}

func (wrapper *byteReaderWrap) ReadByte() (byte, error) {
	buf := make([]byte, 1)

	if _, err := wrapper.reader.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

func ReadVarint(reader io.Reader) (uint32, error) {
	value, err := binary.ReadUvarint(&byteReaderWrap{reader})
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func ScanTarget(target string, timeout time.Duration) {
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return
	}

	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	conn.Write([]byte{
		/*
			https://wiki.vg/Server_List_Ping#Handshake
		*/
		0x07, // pkt lenght
		0x00, // pkt ID
		0x2f, // protocol version
		0x01, // dst. hostname lenght
		0x5f, // dst. hostname
		0x00, // dst. port lenght
		0x01, // dst. port
		0x01, // next state
		/*
			https://wiki.vg/Protocol#Request
		*/
		0x01,
		0x00,
	})

	pkt_len, err := ReadVarint(conn)

	if err != nil {
		return
	}

	pkt_buf := bytes.NewBuffer(nil)

	if _, err = io.CopyN(pkt_buf, conn, int64(pkt_len)); err != nil {
		return
	}

	if packet_id, err := ReadVarint(pkt_buf); err != nil || uint32(packet_id) != uint32(0x00) {
		return
	}

	pkt_data_len, err := ReadVarint(pkt_buf)

	if err != nil {
		return
	}

	pkt_data_buf := make([]byte, pkt_data_len)
	max, err := pkt_buf.Read(pkt_data_buf)

	if err != nil {
		return
	}

	parts := strings.Split(target, ":")
	port_int, _ := strconv.Atoi(parts[1])

	output := &Output{
		Host: parts[0],
		Port: uint16(port_int),
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}

	if json.Unmarshal(pkt_data_buf[:max], &output.Data) != nil {
		return
	}

	json, _ := json.Marshal(output)

	fmt.Println(string(json))
}

func main() {

	threads := flag.Int("threads", 100, "Number of threads")
	timeout := flag.Int("timeout", 5000, "Timeout in milliseconds")

	flag.Parse()

	running := 0

	log.Println("scanning...")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		target := scanner.Text()

		for {

			if running <= *threads {
				break
			}

			time.Sleep(time.Millisecond * 500)
		}

		running++

		go func(target string) {

			defer func() { running-- }()

			ScanTarget(target, time.Millisecond*time.Duration(*timeout))
		}(target)

	}

	log.Println("done.")
}
