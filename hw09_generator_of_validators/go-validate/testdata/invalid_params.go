package testdata

type InvalidUser struct {
	Phones []UserPhone `validate:"len:"`
}
