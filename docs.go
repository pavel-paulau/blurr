package qb

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	sizeOverhead = 450
)

func hash(inString string) string {
	h := md5.New()
	h.Write([]byte(inString))
	return hex.EncodeToString(h.Sum(nil))
}

func randString(key string, expectedLength int) string {
	var newString string
	if expectedLength > 64 {
		baseString := randString(key, expectedLength/2)
		newString = baseString + baseString
	} else {
		newString = (hash(key) + hash(key[:len(key)-1]))[:expectedLength]
	}
	return newString
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func buildAlphabet(key string) string {
	return hash(key) + hash(reverse(key))
}

func buildName(alphabet string) string {
	return fmt.Sprintf("%s %s", alphabet[:6], alphabet[6:12])
}

func buildEmail(alphabet string) string {
	return fmt.Sprintf("%s@%s.com", alphabet[12:18], alphabet[18:24])
}

func buildStreet(alphabet string) string {
	return alphabet[54:62]
}

func buildCity(alphabet string) string {
	return alphabet[24:30]
}

func buildCounty(alphabet string) string {
	return alphabet[48:54]
}

func buildCountry(alphabet string) string {
	return alphabet[42:48]
}

func buildRealm(alphabet string) string {
	return alphabet[30:36]
}

func buildCoins(alphabet string) float64 {
	var coins, _ = strconv.ParseInt(alphabet[36:40], 16, 0)
	return math.Max(0.1, float64(coins)/100.0)
}

func buildCategory(alphabet string) int16 {
	var category, _ = strconv.ParseInt(string(alphabet[41]), 16, 0)
	return int16(category % 3)
}

func buildYear(alphabet string) int16 {
	var year, _ = strconv.ParseInt(string(alphabet[62]), 16, 0)
	return int16(1985 + year)
}

func buildState(alphabet string) string {
	idx := strings.Index(alphabet, "7") % len(unitedStates)
	if idx == -1 {
		idx = 56
	}
	return unitedStates[idx][0]
}

func buildFullState(alphabet string) string {
	idx := strings.Index(alphabet, "8") % len(unitedStates)
	if idx == -1 {
		idx = 56
	}
	return unitedStates[idx][1]
}

func buildGMTime(alphabet string) []int16 {
	var id, _ = strconv.ParseInt(string(alphabet[63]), 16, 0)
	seconds := 396 * 24 * 3600 * (id % 12)
	d := time.Duration(seconds) * time.Second
	t := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC).Add(d)

	return []int16{
		int16(t.Year()),
		int16(t.Month()),
		int16(t.Day()),
		int16(t.Hour()),
		int16(t.Minute()),
		int16(t.Second()),
		int16(t.Weekday() - 1),
		int16(t.YearDay()),
		int16(0),
	}
}

func buildAchievements(alphabet string) (achievements []int16) {
	achievement := int16(256)
	for i, char := range alphabet[42:58] {
		var id, _ = strconv.ParseInt(string(char), 16, 0)
		achievement = (achievement + int16(id)*int16(i)) % 512
		if achievement < 256 {
			achievements = append(achievements, achievement)
		}
	}
	return
}

type doc struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Street       string  `json:"street"`
	City         string  `json:"city"`
	County       string  `json:"county"`
	Country      string  `json:"country"`
	State        string  `json:"state"`
	FullState    string  `json:"fullState"`
	Realm        string  `json:"realm"`
	Coins        float64 `json:"coins"`
	Category     int16   `json:"category"`
	Achievements []int16 `json:"achievements"`
	GMTime       []int16 `json:"gmtime"`
	Year         int16   `json:"year"`
	Body         string  `json:"body"`
}

func newDoc(key string, size int) doc {
	alphabet := buildAlphabet(key)

	return doc{
		Name:         buildName(alphabet),
		Email:        buildEmail(alphabet),
		Street:       buildStreet(alphabet),
		City:         buildCity(alphabet),
		County:       buildCounty(alphabet),
		Country:      buildCountry(alphabet),
		State:        buildState(alphabet),
		FullState:    buildFullState(alphabet),
		Realm:        buildRealm(alphabet),
		Coins:        buildCoins(alphabet),
		Category:     buildCategory(alphabet),
		Achievements: buildAchievements(alphabet),
		GMTime:       buildGMTime(alphabet),
		Year:         buildYear(alphabet),
		Body:         randString(key, size-sizeOverhead),
	}
}
