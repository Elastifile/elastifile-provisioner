package eurl

import (
	"net/url"
	"strconv"

	"github.com/elastifile/errors"
)

type URL struct {
	*url.URL
}

func (u *URL) UnmarshalJSON(data []byte) error {
	if u == nil || u.URL == nil {
		return nil
	}

	rawurl, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	u1, err := url.Parse(rawurl)
	if err != nil {
		return errors.New(err)
	}

	*u = URL{u1}
	return nil
}

func (u *URL) MarshalJSON() ([]byte, error) {
	var rawurl string

	if u.URL != nil {
		rawurl = u.String()
	}

	return []byte(strconv.Quote(rawurl)), nil
}
