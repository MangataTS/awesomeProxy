package ac_automaton

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestACAutomaton(t *testing.T) {
	ac := NewACAutomaton()
	ss := []string{"a", "aa", "aaa"}
	for _, s := range ss {
		ac.Insert(s)
	}
	ac.Build()
	s := `aaaaa`
	res := ac.FindMatches(s)
	t.Log(res)
}

type TestSamp struct {
	Sample []Sample `json:"sample"`
}
type Sample struct {
	Content string   `json:"content"`
	Keyword []string `json:"keyword"`
}

func TestACAutomatonProfile(t *testing.T) {
	sampleFilename := "ac_sample.json"

	file, err := os.Open(sampleFilename)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(file)
	ts := TestSamp{}
	if err != nil {
		t.Error("read json file err: ", err)
	}
	err = json.Unmarshal(data, &ts)
	if err != nil {
		t.Error("Unmarshal json file err: ", err)
		return
	}
	contents := make([]string, 0)
	keywords := make([]string, 0)
	ac := NewACAutomaton()
	for _, sample := range ts.Sample {
		for _, key := range sample.Keyword {
			ac.Insert(key)
			keywords = append(keywords, key)
		}
		contents = append(contents, sample.Content)

	}

	ac.Build()
	for _, content := range contents {
		res := ac.FindMatches(content)
		for _, keyword := range keywords {
			cnt := strings.Count(content, keyword)
			fmt.Println(keyword, ": ", cnt)
			if cnt != res[keyword] {
				t.Error("not equal")
				t.Log(content)
				t.Log(keyword)
				t.Log(cnt)
				t.Log(res[keyword])
				return
			}
		}
	}
}
