package date

import (
	"errors"
	"time"
)

// Date overrides the Marshaler and Unmarshaler interfaces
type Date time.Time

// MarshalJSON ...
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New(
			"Time.MarshalJSON: year outside of range [0,9999]",
		)
	}

	layout := "2006-01-02"
	b := make([]byte, 0, len(layout)+2)

	b = append(b, '"')
	b = t.AppendFormat(b, layout)
	b = append(b, '"')

	return b, nil
}

// UnmarshalJSON ...
func (d Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	layout := "2006-01-02"
	t, err := time.Parse(`"`+layout+`"`, string(data))
	d = Date(t)

	return err
}

func (d Date) String() string {
	t := time.Time(d)
	return t.String()
}
