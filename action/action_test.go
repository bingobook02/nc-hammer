package action

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/damianoneill/nc-hammer/result"
	"github.com/damianoneill/nc-hammer/suite"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	ts1, _ := suite.NewTestSuite("../suite/testdata/test-suite.yml")           // testsuite with netconf actions
	ts2, _ := suite.NewTestSuite("../suite/testdata/testsuite-sleep.yml")      // testsuite with sleep actions
	ts3, _ := suite.NewTestSuite("../suite/testdata/testsuite-failblocks.yml") // test sleep with no sleep or netconf actions
	var buff bytes.Buffer
	var want string

	assertCheck := func(ts *suite.TestSuite) {
		t.Helper()
		start := time.Now()
		resultChannel := make(chan result.NetconfResult)
		handleResultsFinished := make(chan bool)
		go result.HandleResults(resultChannel, handleResultsFinished, ts)
		block := ts.GetInitBlock()
		for _, a := range block.Actions {
			if ts == ts1 {
				Execute(start, 0, ts, a, resultChannel)
				assert.NotNil(t, a.Netconf)
			} else if ts == ts2 {
				Execute(start, 0, ts, a, resultChannel)
				assert.NotNil(t, a.Sleep)
			} else {
				log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
				log.SetOutput(&buff)
				Execute(start, 0, ts, a, resultChannel)
				got := buff.String()
				want += "\n ** Problem with your Testsuite, an action in a block section has incorrect YAML indentation for its body, ensure that anything after netconf or sleep is properly indented **\n\n"
				assert.Equal(t, want, got)
			}
		}
	}
	t.Run("test an action passed is netconf", func(t *testing.T) {
		assertCheck(ts1)
	})
	t.Run("test an action passed is sleep", func(t *testing.T) {
		assertCheck(ts2)
	})
	t.Run("test bad indent or nil netconf&sleep", func(t *testing.T) {
		assertCheck(ts3)
	})

}
