// Package size defines a standard Size type to be used throughout the code for sizes of storage units.
// It should prevent confusion as to what units a size is in.
package size

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// Size is a signed type to avoid overflow problems when doing arithmetic and
// conversions to other signed types.
type Size int64

const MaxSize = math.MaxInt64

const (
	Bytes Size = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB

	Byte = Bytes
)

const (
	KB Size = 1000 * Byte
	MB Size = 1000 * KB
	GB Size = 1000 * MB
	TB Size = 1000 * GB
	PB Size = 1000 * TB
)

const (
	BlockSize = 4 * KiB // ELFS block size
)

func (s Size) String() string {
	var suffix string
	var value float64

	value = float64(s)

	switch {
	case s >= EiB:
		suffix, value = "EiB", value/float64(EiB)
	case s >= PiB:
		suffix, value = "PiB", value/float64(PiB)
	case s >= TiB:
		suffix, value = "TiB", value/float64(TiB)
	case s >= GiB:
		suffix, value = "GiB", value/float64(GiB)
	case s >= MiB:
		suffix, value = "MiB", value/float64(MiB)
	case s >= KiB:
		suffix, value = "KiB", value/float64(KiB)
	case s == 1:
		suffix = "byte"
	case s == 0:
		suffix = ""
	default:
		suffix = "bytes"
	}

	formatted := fmt.Sprintf("%.1f", value)
	formatted = strings.TrimSuffix(formatted, ".0")

	if suffix != "" {
		formatted = formatted + " " + suffix
	}

	return formatted
}

func Parse(s string) (Size, error) {
	var (
		suffix  string
		value   float64
		err     error
		newSize Size
	)

	if s == "" {
		return 0, fmt.Errorf("size: invalid size '%v'", s)
	}
	re := regexp.MustCompile(`(\d+(:?\.\d+)?)\s*(B|(:?[EPTGMK](:?i?B)?))?`)
	groups := re.FindAllStringSubmatch(s, -1)

	if len(groups) == 0 {
		return 0, fmt.Errorf("size: Incorrect format: %s", s)
	}
	value, err = strconv.ParseFloat(groups[0][1], 64)
	if err != nil {
		return 0, fmt.Errorf(
			"size: failed converting '%v' to float64, error=%v",
			groups[0][1], err,
		)
	}

	if len(groups[0]) > 3 {
		suffix = groups[0][3]
	}

	switch suffix {
	case "EiB", "E":
		newSize = Size(value * float64(EiB))
	case "PiB", "P":
		newSize = Size(value * float64(PiB))
	case "TiB", "T":
		newSize = Size(value * float64(TiB))
	case "GiB", "G":
		newSize = Size(value * float64(GiB))
	case "MiB", "M":
		newSize = Size(value * float64(MiB))
	case "KiB", "K":
		newSize = Size(value * float64(KiB))

	case "PB":
		newSize = Size(value * float64(PB))
	case "TB":
		newSize = Size(value * float64(TB))
	case "GB":
		newSize = Size(value * float64(GB))
	case "MB":
		newSize = Size(value * float64(MB))
	case "KB":
		newSize = Size(value * float64(KB))

	default:
		newSize = Size(value)
	}

	return newSize, err
}

func (s *Size) Unmarshal(buf []byte) error {
	r, err := Parse(string(buf))
	if err != nil {
		return err
	}
	*s = r
	return nil
}

// Blocks returns the number of blocks that the Size represents
func (s Size) Blocks() int {
	return int(s / BlockSize)
}

// FromBlocks creates a size that is equal to the specified amount of blocks
func FromBlocks(blocks int) Size {
	return Size(blocks) * BlockSize
}

// Max of two sizes
func Max(x, y Size) Size {
	if x > y {
		return x
	}
	return y
}

// Min of two sizes
func Min(x, y Size) Size {
	if x < y {
		return x
	}
	return y
}

// Abs value of size
func Abs(x Size) Size {
	if x > 0 {
		return x
	}
	return -x
}

// Randn returns a random number in the range [0,n).
func Randn(n Size) Size {
	if n == 0 {
		return Size(0)
	}
	return Size(rand.Int63n(int64(n)))
}

// Similar determines whether the provided sizes are "close enough".
func Similar(x, y Size, epsilon float64) bool {
	if math.Abs(1-float64(x)/float64(y)) > epsilon {
		return false
	}
	return true
}
