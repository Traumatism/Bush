package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
)

type Output struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	Data struct {
		Players struct {
			Online int `json:"online"`
			Max    int `json:"max"`
			Sample []struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"sample,omitempty"`
		} `json:"players"`

		Version struct {
			Name     string `json:"name"`
			Protocol int    `json:"protocol"`
		} `json:"version"`

		Favicon           string                 `json:"favicon"`
		Description       map[string]interface{} `json:"description"`
		ParsedDescription string                 `json:"parsed_description,omitempty"`
	} `json:"data"`
	Date string `json:"date"`
}

func (cmd *Output) hashedFavicon() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(cmd.Data.Favicon)))
}

func (cmd *Output) parsedDescription() string {
	description_dict := cmd.Data.Description
	parsed_description := ""

	if description_dict["text"] != nil {
		parsed_description += description_dict["text"].(string)
	}

	if description_dict["translate"] != nil {
		parsed_description += description_dict["translate"].(string)
	}

	if description_dict["extra"] != nil {
		for _, extra_dict := range description_dict["extra"].([]interface{}) {
			extra_dict := extra_dict.(map[string]interface{})

			if extra_dict["text"] != nil {
				parsed_description += extra_dict["text"].(string)
			}

			if extra_dict["translate"] != nil {
				parsed_description += extra_dict["translate"].(string)
			}
		}
	}

	return parsed_description
}

func parseData(data string) {
	output := Output{}

	err := json.Unmarshal([]byte(data), &output)

	if err != nil {
		return
	}

	output.Data.Favicon = output.hashedFavicon()
	output.Data.ParsedDescription = output.parsedDescription()

	json, err := json.Marshal(output)

	if err != nil {
		return
	}

	fmt.Println(string(json))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		parseData(scanner.Text())
	}
}
