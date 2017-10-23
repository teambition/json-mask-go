package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) TestStringsContains() {
	s.True(stringsContains([]string{"1", "2"}, "1"))
	s.False(stringsContains([]string{"1", "2"}, "3"))
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
