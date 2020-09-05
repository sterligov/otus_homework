package hw02_unpack_string //nolint:golint,stylecheck

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpackWithRandomString(t *testing.T) {
	rand.Seed(time.Now().Unix())

	n := 1000
	var expected, str strings.Builder

	for i := 0; i < n; i++ {
		ch := byte('a' + rand.Intn(26))
		if ch < ('a'+'z')/2 {
			r := rand.Intn(10)
			str.WriteByte(ch)
			str.WriteByte(byte(r) + '0')
			expected.WriteString(strings.Repeat(string(ch), r))
		} else {
			str.WriteByte(ch)
			expected.WriteByte(byte(ch))
		}
	}

	actual, err := Unpack(str.String())

	require.NoError(t, err)
	require.Equal(t, expected.String(), actual)
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
func TestUnpackWithDifferentRunes(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `п 4字2 -! z2Ф0Р1\\u3ß1`,
			expected: `п    字字 -! zzР\uuuß`,
		},
		{
			input:    `\10汉0\20ß0`,
			expected: "",
		},
		{
			input:    "\t2\n0\r1",
			expected: "\t\t\r",
		},
		{
			input:    `й.\\4,2`,
			expected: `й.\\\\,,`,
		},
		{
			input:    `a\字`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `\фb`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `b\п`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `字р\`,
			expected: "",
			err:      ErrInvalidString,
		},
	} {
		tst := tst
		t.Run(tst.input, func(t *testing.T) {
			result, err := Unpack(tst.input)
			require.Equal(t, tst.err, err)
			require.Equal(t, tst.expected, result)
		})
	}
}
