// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sterligov/otus_homework/hw10_program_optimization/mock"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("invalid email", func(t *testing.T) {
		data := `{"Email":"invalid_email.com"}`
		actual, err := GetDomainStat(bytes.NewBufferString(data), "com")

		require.Empty(t, actual)
		require.Error(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		data := `{"Email":"`
		actual, err := GetDomainStat(bytes.NewBufferString(data), "com")

		require.Empty(t, actual)
		require.Error(t, err)
	})

	t.Run("empty data", func(t *testing.T) {
		actual, err := GetDomainStat(bytes.NewBufferString(""), "com")

		require.Empty(t, actual)
		require.Nil(t, err)
	})

	t.Run("reading error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mock.NewMockReader(ctrl)
		m.EXPECT().Read(gomock.Any()).Return(0, fmt.Errorf("error"))

		actual, err := GetDomainStat(m, "com")

		require.Empty(t, actual)
		require.Error(t, err)
	})
}
