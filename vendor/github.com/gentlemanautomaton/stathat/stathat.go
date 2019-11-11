package stathat

import "errors"

// StatHat is a StatHat which does StatHat things
type StatHat struct {
	token string
	ezkey string
	noop  bool
}

// ErrMissingToken means the `StatHat` was missing an `token` value.
// That value is set usually set like this: `stathat.New().Token("sometoken").EZKey("somekey")` (both `Token` and `EZKey` are optional, but various methods need one or the other to have been set).
var ErrMissingToken = errors.New("missing token")

// ErrMissingEZKey means the `StatHat` was missing an `ezkey` value.
// That value is set usually set like this: `stathat.New().Token("sometoken").EZKey("somekey")` (both `Token` and `EZKey` are optional, but various methods need one or the other to have been set).
var ErrMissingEZKey = errors.New("missing ezkey")

// New returns a StatHat.
// You'll want to set either/both of the `Token` or `EZKey` values.
func New() StatHat {
	return StatHat{}
}

// Noop causes any stats to be silently dropped instead of being recorded.
// This is useful at times when you want to opt out of recording stats but not want to unwire the stat calls.
// Writes won't happen (deleting and stat value sending), but any reads will still fetch.
// Writes will return nil for error, since they didn't fail but were simply ignored.
func (s StatHat) Noop() StatHat {
	s.noop = true
	return s
}

// Token sets the access token (see https://www.stathat.com/access).
// It returns an updated StatHat with the token added.  It does not modify the original StatHat.
func (s StatHat) Token(token string) StatHat {
	s.token = token
	return s
}

// EZKey sets the ezkey value.
// It returns an updated StatHat with the token added.  It does not modify the original StatHat.
func (s StatHat) EZKey(ezkey string) StatHat {
	s.ezkey = ezkey
	return s
}

func (s StatHat) apiPrefix() string {
	return `https://www.stathat.com/x/` + s.token
}

func (s StatHat) ezPrefix() string {
	return `https://api.stathat.com/ez`
}

func (s StatHat) classicPrefix() string {
	return `https://api.stathat.com/c`
}
