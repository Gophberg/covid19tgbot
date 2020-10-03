package bot

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeNums ...
func ScrapeNums(config *Config) (int, error) {
	log.Println("Scraper start")
	res, err := http.Get(config.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Responce error: %d %s", res.StatusCode, res.Status)
	}
	log.Printf("GET status: %v", res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var strs []string
	doc.Find(".block_banner_b .block_banner_b_info").Each(func(i int, s *goquery.Selection) {
		band := s.Find(".block_banner_b_info_t").Text()
		strs = append(strs, band)
	})
	var nums []int
	for i := range strs {
		nums = append(nums, extr(strs[i]))
	}
	casesNum := Nums{
		Cases:                nums[0],
		Deaths:               nums[1],
		Recovered:            nums[2],
		Isolation:            nums[3],
		Isolationathome:      nums[4],
		Observation:          nums[5],
		Quarantineathospital: nums[6],
		Unquarantined:        nums[7],
		Peopleathospital:     nums[8],
		Ambulanced:           nums[9],
	}
	log.Printf("Scrapped data: %v", casesNum)
	return casesNum.Cases, err
}

func extr(s string) int {
	start := strings.Index(s, ":")
	end := strings.Index(s, "(")
	t := strings.TrimSpace(s[start+1 : end-1])
	sInt, err := strconv.Atoi(t)
	if err != nil {
		log.Panic(err)
	}
	return sInt
}
