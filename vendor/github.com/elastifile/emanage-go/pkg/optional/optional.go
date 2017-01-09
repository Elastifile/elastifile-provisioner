// The optional package is used for building structs with optional values that can be marshaled to JSON.
// This is necessary to distinguish between zero values, missing values, and actual nil (JSON "null") values.
//
// This issue is a tough problem, and various approaches exist.
// Some interesting references:
// - https://golang.org/pkg/database/sql/#NullString
// - https://willnorris.com/2014/05/go-rest-apis-and-pointers
//
// TODO: A more elegant way to handle this would be to use opaque data types that implement public interfaces,
// like sql.NullString does: Use interfaces for basic types that have methods for marshaling and unmarshaling.
// Similar to what we have started to do for emanage.DedupStatus etc.
//
// NOTE: The emanage REST API, by which this package is used, supports optional values.
// In addition, according to Amos, if a value is sent as nil/null, it is considered by emanage as if it isn't sent at all.
// Therefore, there is no need for us to explicitly support sending nil values, which simplifies matters quite a bit.
package optional

type String *string

func NewString(s string) String {
	return &s
}

type Int *int

func NewInt(i int) Int {
	return &i
}

type Float64 *float64

func NewFloat64(f float64) Float64 {
	return &f
}

type Bool *bool

func NewBool(b bool) Bool {
	return &b
}

func True() Bool {
	return NewBool(true)
}

func False() Bool {
	return NewBool(false)
}
