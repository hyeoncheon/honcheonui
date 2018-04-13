package utils_test

import (
	"testing"

	"github.com/hyeoncheon/honcheonui/utils"
	"github.com/stretchr/testify/require"
)

// ToInterface ---
func Test_ToInterface(t *testing.T) {
	r := require.New(t)
	sa := []string{"apple", "banana"}
	ia := []interface{}{"apple", "banana"}

	ret, err := utils.ToInterface(sa)
	r.NoError(err)
	r.NotNilf(ret, "returns '%v' (%T)", ret, ret)
	r.Equal(ia, ret)
}

func Test_ToInterface_ErrorNotArrayForStringOrInt(t *testing.T) {
	r := require.New(t)

	ret, err := utils.ToInterface("apple")
	r.Nil(ret)
	r.Error(err)
	r.EqualError(err, "not array or slice")

	ret, err = utils.ToInterface(1)
	r.Nil(ret)
	r.Error(err)
	r.EqualError(err, "not array or slice")
}

// Has
func Test_Has(t *testing.T) {
	r := require.New(t)
	ia := []interface{}{"apple", "banana", 1}

	ret := utils.Has(ia, "apple")
	r.True(ret)

	ret = utils.Has(ia, 1)
	r.True(ret)

	ret = utils.Has(ia, "pineapple")
	r.False(ret)
}

// Remove
func Test_Remove(t *testing.T) {
	r := require.New(t)
	ia := []interface{}{"apple", "banana", 1}

	na := utils.Remove(ia, "banana")
	r.Equal([]interface{}{"apple", 1}, na)
	// ia is changed. see comment on function body.
	r.Equal([]interface{}{"apple", 1, 1}, ia)

	na = utils.Remove(ia, "pineapple")
	r.Equal([]interface{}{"apple", 1, 1}, na)
}

// Cleaner
func Test_Cleaner(t *testing.T) {
	r := require.New(t)

	na := utils.Cleaner([]string{"apple", "banana"})
	r.Equal([]string{"apple", "banana"}, na)

	na = utils.Cleaner([]string{"   apple", "banana  ", " o range "})
	r.Equal([]string{"apple", "banana", "o range"}, na)
}
