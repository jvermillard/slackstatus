package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func setStatus(token string, status string, emoji string) {
	fmt.Println("set status", status, emoji)

	profile := url.QueryEscape("{\"status_text\": \"" + status + "\",\"status_emoji\": \"" + emoji + "\"}")

	resp, err := http.Post("https://slack.com/api/users.profile.set?token="+token+"&profile="+profile, "", nil)
	if err != nil {
		log.Fatal("Error calling slack API", err)
	}

	fmt.Println("slack API status:", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading slack API return", err)
	}

	answer := make(map[string]interface{})

	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Fatal("Error decoding slack API json", err)
	}

	pretty, err := json.MarshalIndent(answer, "", "  ")
	if err != nil {
		log.Fatal("Error pretty printing slack API json", err)
	}

	fmt.Println("return:\n", string(pretty))

}

func main() {
	var token string
	var filename string
	flag.StringVar(&token, "token", "", "slack API token")
	flag.StringVar(&filename, "file", "", "csv file containing locations")

	flag.Parse()

	if len(token) <= 0 || len(filename) <= 0 {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	locations := make(map[string][]string, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		locations[record[0]] = []string{record[1], record[2]}
	}

	fmt.Println("check SSID")

	var ssid string
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
		out, err := cmd.Output()
		if err != nil {
			if exerr, ok := err.(*exec.ExitError); ok {
				fmt.Println("no wifi", exerr)
				out = make([]byte, 0)
			} else {
				log.Fatal(err)
			}
		}
		r, _ := regexp.Compile("SSID: [a-z]+")
		ssid = r.FindString(string(out))[6:]
	} else {
		cmd := exec.Command("iwgetid", "-r")
		out, err := cmd.Output()
		if err != nil {
			if exerr, ok := err.(*exec.ExitError); ok {
				fmt.Println("no wifi", exerr)
				out = make([]byte, 0)
			} else {
				log.Fatal("You must install `iwgetid`", err)
			}
		}
		ssid = strings.Replace(string(out), "\n", "", -1)
	}

	fmt.Println("ssid", ssid)

	// look at outside IP
	ip := ""
	countryCode := ""

	resp, err := http.Get("http://ip-api.com/json")

	if err != nil {
		fmt.Println("can't get IP", err)
	} else {
		loc := struct {
			City        string
			Country     string
			CountryCode string
			Query       string
		}{}
		data, err :=
			ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("can't read IP", err)
		} else {
			err = json.Unmarshal(data, &loc)
			if err != nil {
				fmt.Println("can't parse IP", err)
			} else {
				fmt.Println(loc)
				ip = loc.Query
				countryCode = loc.CountryCode
				fmt.Println(countryCode, ip)
			}
		}
	}

	// look if we know the SSID, then the outside IP, and finally try to put in travel mode
	wifiloc := locations[ssid]
	iploc := locations[ip]
	if wifiloc != nil {
		setStatus(token, wifiloc[0], wifiloc[1])
	} else if iploc != nil {
		setStatus(token, iploc[0], iploc[1])
	} else {
		setStatus(token, "traveling :flag-"+strings.ToLower(countryCode)+":", ":airplane:")
	}
}
