package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
