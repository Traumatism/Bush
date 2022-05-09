package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Description struct {
	/* Minecraft description (Chat component) structure */
	Text  string `json:"text,omitempty"`
	Extra []struct {
		Text string `json:"text,omitempty"`
	} `json:"extra,omitempty"`
}

type Players struct {
	/* Minecraft player structure */
	Max    int `json:"max"`
	Online int `json:"online"`
	Sample []struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	} `json:"sample,omitempty"`
}

type Response struct {
	/* Minecraft SLP response structure */
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players     Players     `json:"players"`
	Description Description `json:"description"`
}

type ResponseB struct {
	/* Same as `Response` but with Description as a string */
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players     Players `json:"players"`
	Description string  `json:"description"`
}

type Output struct {
	Host        string  `json:"host"`
	Port        uint16  `json:"port"`
	Version     string  `json:"version"`
	Protocol    int     `json:"protocol"`
	Players     Players `json:"players"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (output *Output) Human() string {
	return fmt.Sprintf("(%s:%d)(%d/%d)(%s)(%s)", output.Host, output.Port, output.Players.Online, output.Players.Max, output.Version, output.Description)
}

func (output *Output) Json() string {
	json, _ := json.Marshal(output)
	return string(json)
}

type RawOutput struct {
	Host string                 `json:"host"`
	Port uint16                 `json:"port"`
	Data map[string]interface{} `json:"data"`
	Date string                 `json:"date"`
}

func (output *RawOutput) Json() string {
	json, _ := json.Marshal(output)
	return string(json)
}

type byteReaderWrap struct {
	reader io.Reader
}

func (wrapper *byteReaderWrap) ReadByte() (byte, error) {
	/* Read one byte from an I/O stream */
	buf := make([]byte, 1)

	if _, err := wrapper.reader.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

type Stats struct {
	sent    int
	total   int
	errors  int
	success int
}

func (s *Stats) Inc() {
	s.sent++
}

func (s *Stats) IncErrors() {
	s.errors++
}

func (s *Stats) IncSuccess() {
	s.success++
}

func (s *Stats) Display() string {
	return fmt.Sprintf(
		"%d sent / %d total (%d%% done) %d success / %d errors (%d%% hit)",
		s.sent, s.total, int(float64(s.sent)/float64(s.total)*100),
		s.success, s.errors, int(float64(s.success)/float64(s.sent)*100),
	)
}

func ReadVarint(r io.Reader) (uint32, error) {
	/* Read a varint from an I/O stream */
	value, err := binary.ReadUvarint(&byteReaderWrap{r})
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func ScanTarget(target string, timeout time.Duration, format string) int {
	/* Scan a target */
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return 0
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

	// read total packet lenght (pkt ID + pkt data lenght + pkt data)
	total_lenght, err := ReadVarint(conn)

	if err != nil {
		return 0
	}

	buf_total := bytes.NewBuffer(nil)

	if _, err = io.CopyN(buf_total, conn, int64(total_lenght)); err != nil {
		return 0
	}

	// read pkt ID (should be 0x00 => handshake response)
	if packet_id, err := ReadVarint(buf_total); err != nil || uint32(packet_id) != uint32(0x00) {
		return 0
	}

	// read pkt data lenght
	lenght, err := ReadVarint(buf_total)

	if err != nil {
		return 0
	}

	// read pkt data
	buf_data := make([]byte, lenght)
	max, err := buf_total.Read(buf_data)

	if err != nil {
		return 0
	}

	conn.Close()

	data := buf_data[:max]
	parts := strings.Split(target, ":")
	port_int, _ := strconv.Atoi(parts[1])

	if format == "raw" {

		output := RawOutput{
			Host: parts[0],
			Port: uint16(port_int),
			Date: time.Now().Format("2006-01-02 15:04:05"),
		}

		if err := json.Unmarshal(data, &output.Data); err != nil {
			return 0
		}

		fmt.Println(output.Json())
		return 1
	}

	output := Output{
		Host: parts[0],
		Port: uint16(port_int),
		Date: time.Now().Format(time.RFC3339),
	}

	var response Response

	if json.Unmarshal(data, &response) != nil {
		var responseB ResponseB

		if json.Unmarshal(data, &responseB) != nil {
			return 0
		}

		output.Version = responseB.Version.Name
		output.Protocol = responseB.Version.Protocol
		output.Players = responseB.Players
	} else {
		output.Version = response.Version.Name
		output.Protocol = response.Version.Protocol
		output.Players = response.Players
	}

	description := response.Description

	if description.Text != "" {
		output.Description = description.Text
	} else {
		output.Description = ""
	}

	if len(description.Extra) > 0 {
		for _, extra := range description.Extra {
			output.Description += extra.Text
		}
	}

	fmt.Println(func() string {
		if format == "json" {
			return output.Json()
		}
		return output.Human()
	}())

	return 1
}

func ShowStats(stats *Stats) {
	for stats.sent < stats.total {
		fmt.Fprintf(os.Stderr, "%s\n", stats.Display())
		time.Sleep(1 * time.Second)
	}
}

func Wait(running *int, threads *int) {
	for {

		if *running <= *threads {
			break
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func main() {

	threads := flag.Int("threads", 100, "Number of threads")
	timeout := flag.Int("timeout", 5000, "Timeout in milliseconds")
	format := flag.String("format", "human", "Output format (json, human, raw)")
	stats := flag.Bool("stats", false, "Show stats every second (note: it wille disable the stream mode)")

	flag.Parse()

	if *format != "json" && *format != "human" && *format != "raw" {
		fmt.Fprintln(os.Stderr, "Invalid format")
		os.Exit(1)
	}

	running := 0

	scanner := bufio.NewScanner(os.Stdin)

	if *stats {
		targets := []string{}

		for scanner.Scan() {
			targets = append(targets, scanner.Text())
		}

		stats_cls := Stats{
			total: len(targets),
		}

		go ShowStats(&stats_cls)

		for _, target := range targets {

			Wait(&running, threads)

			running++

			go func(target string) {

				defer func() {
					running--
					stats_cls.Inc()
				}()

				switch ScanTarget(target, time.Millisecond*time.Duration(*timeout), *format) {
				case 1:
					stats_cls.success++
				case 0:
					stats_cls.errors++
				}
			}(target)
		}
	} else {
		for scanner.Scan() {

			target := scanner.Text()

			Wait(&running, threads)

			running++

			go func(target string) {

				defer func() {
					running--
				}()

				ScanTarget(target, time.Millisecond*time.Duration(*timeout), *format)
			}(target)
		}
	}

	for running > 0 {
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Fprintln(os.Stderr, "Done")
}
