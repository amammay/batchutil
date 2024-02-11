package batchutil_test

import (
	"errors"
	"testing"

	"github.com/amammay/batchutil"
	"github.com/stretchr/testify/require"
)

func testSlice(size int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	return slice
}

func Test(t *testing.T) {

	var ranges [][]int
	err := batchutil.All(testSlice(100), 10, func(vals []int) error {
		ranges = append(ranges, vals)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, ranges, 10)
	require.Len(t, ranges[0], 10)
	require.Len(t, ranges[1], 10)
	require.Len(t, ranges[2], 10)
	require.Len(t, ranges[3], 10)
	require.Len(t, ranges[4], 10)
	require.Len(t, ranges[5], 10)
	require.Len(t, ranges[6], 10)
	require.Len(t, ranges[7], 10)
	require.Len(t, ranges[8], 10)
	require.Len(t, ranges[9], 10)

	require.Equal(t, ranges[0][0], 0)
	require.Equal(t, ranges[0][9], 9)
	require.Equal(t, ranges[1][0], 10)
	require.Equal(t, ranges[1][9], 19)
	require.Equal(t, ranges[2][0], 20)
	require.Equal(t, ranges[2][9], 29)
	require.Equal(t, ranges[3][0], 30)
	require.Equal(t, ranges[3][9], 39)
	require.Equal(t, ranges[4][0], 40)
	require.Equal(t, ranges[4][9], 49)
	require.Equal(t, ranges[5][0], 50)
	require.Equal(t, ranges[5][9], 59)
	require.Equal(t, ranges[6][0], 60)
	require.Equal(t, ranges[6][9], 69)
	require.Equal(t, ranges[7][0], 70)
	require.Equal(t, ranges[7][9], 79)
	require.Equal(t, ranges[8][0], 80)
	require.Equal(t, ranges[8][9], 89)
	require.Equal(t, ranges[9][0], 90)
	require.Equal(t, ranges[9][9], 99)

}

func TestHalfPages(t *testing.T) {
	var ranges [][]int
	err := batchutil.All(testSlice(15), 10, func(vals []int) error {
		ranges = append(ranges, vals)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, ranges, 2)

	require.Len(t, ranges[0], 10)
	require.Len(t, ranges[1], 5)

	require.Equal(t, ranges[0][0], 0)
	require.Equal(t, ranges[0][9], 9)
	require.Equal(t, ranges[1][0], 10)
	require.Equal(t, ranges[1][4], 14)

}

func TestTinyPages(t *testing.T) {

	var ranges [][]int
	err := batchutil.All(testSlice(1), 10, func(vals []int) error {
		ranges = append(ranges, vals)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, ranges, 1)
	require.Len(t, ranges[0], 1)
	require.Equal(t, ranges[0][0], 0)
}

func TestSliceOfStrings(t *testing.T) {

	var got [][]string
	stringSlice := []string{"one", "two", "three", "four", "five"}
	err := batchutil.All(stringSlice, 2, func(vals []string) error {
		got = append(got, vals)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, got, 3)
	require.Len(t, got[0], 2)
	require.Len(t, got[1], 2)
	require.Len(t, got[2], 1)

	require.Equal(t, got[0][0], "one")
	require.Equal(t, got[0][1], "two")
	require.Equal(t, got[1][0], "three")
	require.Equal(t, got[1][1], "four")
	require.Equal(t, got[2][0], "five")
}

func TestAbort(t *testing.T) {

	var s []string
	err := batchutil.All(s, 10, func(vals []string) error {
		return batchutil.ErrAbort
	})
	require.NoError(t, err)
}

func TestErr(t *testing.T) {
	errTest := errors.New("something went wrong")
	err := batchutil.All(testSlice(1), 10, func(val []int) error {
		return errTest
	})
	require.ErrorIs(t, err, errTest)
}
