package workloads

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

	"github.com/pavel-paulau/nb/databases"
)

const BatchSize int = 100

func Hash(inString string) string {
	h := md5.New()
	h.Write([]byte(inString))
	return hex.EncodeToString(h.Sum(nil))
}

func RandString(key string, expectedLength int) string {
	var randString string
	if expectedLength > 64 {
		baseString := RandString(key, expectedLength/2)
		randString = baseString + baseString
	} else {
		randString = (Hash(key) + Hash(key[:len(key)-1]))[:expectedLength]
	}
	return randString
}

type N1QL struct {
	Config       Config
	DeletedItems int64
	Zipf         rand.Zipf
	i            Workload
}

func (w *N1QL) SetImplementation(i Workload) {
	w.i = i
}

func (w *N1QL) GenerateNewKey(currentRecords int64) string {
	return fmt.Sprintf("%012d", currentRecords)
}

func (w *N1QL) GenerateExistingKey(currentRecords int64) string {
	randRecord := 1 + rand.Int63n(currentRecords-w.DeletedItems)
	randRecord += w.DeletedItems
	strRandRecord := strconv.FormatInt(randRecord, 10)
	return Hash(strRandRecord)
}

func (w *N1QL) GenerateKeyForRemoval() string {
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
	return Hash(key) + Hash(reverse(key))
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
	idx := strings.Index(alphabet, "7") % NUM_STATES
	if idx == -1 {
		idx = 56
	}
	return STATES[idx][0]
}

func buildFullState(alphabet string) string {
	idx := strings.Index(alphabet, "8") % NUM_STATES
	if idx == -1 {
		idx = 56
	}
	return STATES[idx][1]
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

var OVERHEAD = int(450)

func (w *N1QL) GenerateValue(key string, size int) map[string]interface{} {
	if size < OVERHEAD {
		log.Fatalf("Wrong workload configuration: minimal value size is %v", OVERHEAD)
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
		"body":         RandString(key, size-OVERHEAD),
	}
}

func (w *N1QL) PrepareBatch() []string {
	operations := make([]string, 0, BatchSize)
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
	if len(operations) != BatchSize {
		log.Fatal("Wrong workload configuration: sum of percentages is not equal 100")
	}
	return operations
}

func (w *N1QL) PrepareSeq(size int64) chan string {
	operations := w.PrepareBatch()
	seq := make(chan string, BatchSize)
	go func() {
		for i := int64(0); i < size; i += int64(BatchSize) {
			for _, randI := range rand.Perm(BatchSize) {
				seq <- operations[randI]
			}
		}
	}()
	return seq
}

func (w *N1QL) DoBatch(db databases.Database, state *State, seq chan string) {
	for i := 0; i < BatchSize; i++ {
		op := <-seq
		if state.Operations < w.Config.Operations {
			var err error
			state.Operations++
			switch op {
			case "c":
				state.Records++
				key := w.i.GenerateNewKey(state.Records)
				value := w.i.GenerateValue(key, w.Config.ValueSize)
				err = db.Create(key, value)
			case "r":
				key := w.i.GenerateExistingKey(state.Records)
				err = db.Read(key)
			case "u":
				key := w.i.GenerateExistingKey(state.Records)
				value := w.i.GenerateValue(key, w.Config.ValueSize)
				err = db.Update(key, value)
			case "d":
				key := w.i.GenerateKeyForRemoval()
				err = db.Delete(key)
			}
			if err != nil {
				fmt.Println(err)
				state.Errors[op]++
				state.Errors["total"]++
			}
		}
	}
}

func (w *N1QL) runWorkload(database databases.Database,
	state *State, wg *sync.WaitGroup, targetBatchTimeF float64, seq chan string) {

	for state.Operations < w.Config.Operations {
		t0 := time.Now()
		w.i.DoBatch(database, state, seq)
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

func (w *N1QL) RunCRUDWorkload(database databases.Database,
	state *State, wg *sync.WaitGroup) {
	defer wg.Done()

	seq := w.PrepareSeq(w.Config.Operations)
	targetBatchTimeF := float64(BatchSize) / float64(w.Config.Throughput)
	w.runWorkload(database, state, wg, targetBatchTimeF, seq)
}
