package linkpearl

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type linkpearlSuite struct {
	suite.Suite
}

func (s *linkpearlSuite) SetupTest() {
	s.T().Helper()
}

func (s *linkpearlSuite) TearDownTest() {
	s.T().Helper()
}

func TestLinkpearl(t *testing.T) {
	suite.Run(t, &linkpearlSuite{})
}
