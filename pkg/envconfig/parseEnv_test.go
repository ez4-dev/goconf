package envconfig_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trinhpham/go-conf/pkg/envconfig"
)

type SimpleConfig struct {
	StrValue   string  `env:"STR"`
	IntValue   int     `env:"INT"`
	UIntValue  uint    `env:"UINT"`
	FloatValue float32 `env:"FLOAT"`
	BoolValue  bool    `env:"BOOL"`
}

type NestedConfig struct {
	StrValue   string                  `env:"STR"`
	MapInt     map[string]int8         `env:"MINT"`
	SliceFloat []float64               `env:"SFL"`
	MapValue   map[string]SimpleConfig `env:"MAP"`
	SliceVal   []SimpleConfig          `env:"SLICE"`
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type ParseEnvSuite struct {
	suite.Suite
}

func TestParseEnvSuite(t *testing.T) {
	suite.Run(t, new(ParseEnvSuite))
}

// Clear env before test
func (suite *ParseEnvSuite) SetupTest() {
	os.Clearenv()
}

func (suite *ParseEnvSuite) TestGetEnv() {
	os.Setenv("CC", "DD")
	assert.New(suite.T())
	suite.Equal("BB", envconfig.GetEnv("AA", "BB"))
	suite.Equal("DD", envconfig.GetEnv("CC", "EE"))
}

func (suite *ParseEnvSuite) TestParseSimpleConfigFromEnv() {
	// Given
	os.Setenv("PREFIX_STR", "abc")
	os.Setenv("PREFIX_INT", "-123")
	os.Setenv("PREFIX_UINT", "12")
	os.Setenv("PREFIX_FLOAT", "12.34")
	os.Setenv("PREFIX_BOOL", "true")

	config := SimpleConfig{}
	envconfig.Load("PREFIX", &config)

	suite.Equal("abc", config.StrValue)
	suite.Equal(-123, config.IntValue)
	suite.EqualValues(12, config.UIntValue)
	suite.EqualValues(12.34, config.FloatValue)
}

func (suite *ParseEnvSuite) TestParseNestedConfigFromEnv() {
	// Given
	os.Setenv("PREFIX_STR", "abc")
	os.Setenv("PREFIX_MINT_ITEM1", "1")
	os.Setenv("PREFIX_MINT_ITEM2", "2")
	os.Setenv("PREFIX_SFL_1", "12.34")
	os.Setenv("PREFIX_SFL_4", "56.78")
	os.Setenv("PREFIX_MAP_S1_STR", "abc")
	os.Setenv("PREFIX_SLICE_5_STR", "xyz")
	os.Setenv("PREFIX_SLICE_5_BOOL", "n")

	config := NestedConfig{}
	envconfig.Load("PREFIX", &config)

	suite.EqualValues("abc", config.StrValue)
	suite.EqualValues(1, config.MapInt["item1"])
	suite.EqualValues(2, config.MapInt["item2"])
	suite.EqualValues(12.34, config.SliceFloat[1])
	suite.EqualValues(56.78, config.SliceFloat[4])
	suite.EqualValues("xyz", config.SliceVal[5].StrValue)
	suite.EqualValues(false, config.SliceVal[5].BoolValue)
}

func (suite *ParseEnvSuite) TestParseToMap() {
	os.Setenv("PREFIX_ITEM1", "1")
	os.Setenv("PREFIX_ITEM2", "2")

	config := make(map[string]int64)
	envconfig.Load("PREFIX", &config)
	suite.EqualValues(1, config["item1"])
	suite.EqualValues(2, config["item2"])
}

func (suite *ParseEnvSuite) TestParseToBool() {
	os.Setenv("P", "f")
	bVal := true
	envconfig.Load("P", &bVal)
	suite.EqualValues(false, bVal)
}

func (suite *ParseEnvSuite) TestPanicOnInvalidType() {
	var tChan chan string
	var tBool bool
	var tInt int
	var tUint uint
	var tFloat float32
	var tFloat64 float64
	var tSlice []string
	os.Setenv("P_X", "x")

	suite.Panics(func() { envconfig.Load("", nil) })
	suite.Panics(func() { envconfig.Load("P", &tChan) })
	suite.Panics(func() { envconfig.Load("P_X", &tBool) })
	suite.Panics(func() { envconfig.Load("P_X", &tInt) })
	suite.Panics(func() { envconfig.Load("P_X", &tUint) })
	suite.Panics(func() { envconfig.Load("P_X", &tFloat) })
	suite.Panics(func() { envconfig.Load("P_X", &tFloat64) })
	suite.Panics(func() { envconfig.Load("P", &tFloat64) })
	suite.Panics(func() { envconfig.Load("P", &tSlice) })
}
