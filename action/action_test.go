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
	var buff bytes.Buffer
	var want string
	ts1, _ := suite.NewTestSuite("../suite/testdata/test-suite.yml")             // testsuite with netconf actions
	ts2, _ := suite.NewTestSuite("../suite/testdata/testsuite-fail-actions.yml") // testsuite with no netconf or sleep actions
	myTests := []*suite.TestSuite{ts1, ts2}
	start := time.Now()
	resultChannel := make(chan result.NetconfResult)
	handleResultsFinished := make(chan bool)
	go result.HandleResults(resultChannel, handleResultsFinished, ts1)
	for _, myTs := range myTests {
		for _, b := range myTs.Blocks {
			for _, a := range b.Actions {
				log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
				log.SetOutput(&buff)
				if myTs == ts1 {
					Execute(start, 0, ts1, a, resultChannel)
					assert.True(t, (a.Sleep != nil) || (a.Netconf != nil)) // checks for netconf or sleep actions
				} else {
					Execute(start, 0, ts2, a, resultChannel)
					got := buff.String()
					want += "\n ** Problem with your Testsuite, an action in a block section has incorrect YAML indentation for its body, ensure that anything after netconf or sleep is properly indented **\n\n"
					assert.Equal(t, want, got)
				}
			}
		}
	}
}
