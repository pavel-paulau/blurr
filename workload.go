package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	batchSize    int = 100
	sizeOverhead int = 450
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

type Workload struct {
	Config       workloadConfig
	DeletedItems int64
}

func (w *Workload) generateNewKey(currentDocuments int64) string {
	return fmt.Sprintf("%012d", currentDocuments)
}

func (w *Workload) generateExistingKey(currentDocuments int64) string {
	randRecord := 1 + rand.Int63n(currentDocuments-w.DeletedItems)
	randRecord += w.DeletedItems
	strRandRecord := strconv.FormatInt(randRecord, 10)
	return hash(strRandRecord)
}

func (w *Workload) generateKeyForRemoval() string {
	w.DeletedItems++
	return fmt.Sprintf("%012d", w.DeletedItems)
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

func (w *Workload) generateValue(key string, size int) map[string]interface{} {
	if size < sizeOverhead {
		log.Fatalf("Wrong workload configuration: minimal value size is %v", sizeOverhead)
	}

	alphabet := buildAlphabet(key)

	return map[string]interface{}{
		"name":         buildName(alphabet),
		"email":        buildEmail(alphabet),
		"street":       buildStreet(alphabet),
		"city":         buildCity(alphabet),
		"county":       buildCounty(alphabet),
		"country":      buildCountry(alphabet),
		"state":        buildState(alphabet),
		"full_state":   buildFullState(alphabet),
		"realm":        buildRealm(alphabet),
		"coins":        buildCoins(alphabet),
		"category":     buildCategory(alphabet),
		"achievements": buildAchievements(alphabet),
		"gmtime":       buildGMTime(alphabet),
		"year":         buildYear(alphabet),
		"body":         randString(key, size-sizeOverhead),
	}
}

func (w *Workload) prepareBatch() []string {
	operations := make([]string, 0, batchSize)
	for i := 0; i < w.Config.CreatePercentage; i++ {
		operations = append(operations, "c")
	}
	for i := 0; i < w.Config.ReadPercentage; i++ {
		operations = append(operations, "r")
	}
	for i := 0; i < w.Config.UpdatePercentage; i++ {
		operations = append(operations, "u")
	}
	for i := 0; i < w.Config.DeletePercentage; i++ {
		operations = append(operations, "d")
	}
	if len(operations) != batchSize {
		log.Fatal("Wrong workload configuration: sum of percentages is not equal 100")
	}
	return operations
}

func (w *Workload) prepareSeq(size int64) chan string {
	operations := w.prepareBatch()
	seq := make(chan string, batchSize)
	go func() {
		for i := int64(0); i < size; i += int64(batchSize) {
			for _, randI := range rand.Perm(batchSize) {
				seq <- operations[randI]
			}
		}
	}()
	return seq
}

func (w *Workload) doBatch(client *Client, state *State, seq chan string) {
	for i := 0; i < batchSize; i++ {
		op := <-seq
		if state.Operations < w.Config.Operations {
			var err error
			state.Operations++
			switch op {
			case "c":
				state.Documents++
				key := w.generateNewKey(state.Documents)
				value := w.generateValue(key, w.Config.ValueSize)
				err = client.create(key, value)
			case "r":
				key := w.generateExistingKey(state.Documents)
				err = client.read(key)
			case "u":
				key := w.generateExistingKey(state.Documents)
				value := w.generateValue(key, w.Config.ValueSize)
				err = client.update(key, value)
			case "d":
				key := w.generateKeyForRemoval()
				err = client.delete(key)
			}
			if err != nil {
				state.Errors[op]++
				state.Errors["total"]++
			}
		}
	}
}

func (w *Workload) runWorkload(client *Client, state *State, wg *sync.WaitGroup, targetBatchTimeF float64, seq chan string) {
	for state.Operations < w.Config.Operations {
		t0 := time.Now()
		w.doBatch(client, state, seq)
		t1 := time.Now()

		if !math.IsInf(targetBatchTimeF, 0) {
			targetBatchTime := time.Duration(targetBatchTimeF * math.Pow10(9))
			actualBatchTime := t1.Sub(t0)
			sleepTime := (targetBatchTime - actualBatchTime)
			if sleepTime > 0 {
				time.Sleep(time.Duration(sleepTime))
			}
		}
	}
}

func (w *Workload) runCRUDWorkload(client *Client, state *State, wg *sync.WaitGroup) {
	defer wg.Done()

	seq := w.prepareSeq(w.Config.Operations)
	targetBatchTimeF := float64(batchSize) / float64(w.Config.Throughput)
	w.runWorkload(client, state, wg, targetBatchTimeF, seq)
}
