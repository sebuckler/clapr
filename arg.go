package clapr

import (
	"fmt"
	"strconv"
	"strings"
)

type Arg struct {
	Binder     ArgBinder
	IsHelp     bool
	Name       string
	ShortName  rune
	Repeatable bool
	Required   bool
	Usage      string
}

type ArgBinder interface {
	Bind(arg string, val string) error
}

type boolArgBinder struct {
	val *bool
}

type float64Binder struct {
	val *float64
}

type float64ListBinder struct {
	val *[]float64
}

type intBinder struct {
	val *int
}

type intListBinder struct {
	val *[]int
}

type int64Binder struct {
	val *int64
}

type int64ListBinder struct {
	val *[]int64
}

type stringBinder struct {
	val *string
}

type stringListBinder struct {
	val *[]string
}

type uintBinder struct {
	val *uint
}

type uintListBinder struct {
	val *[]uint
}

type uint64Binder struct {
	val *uint64
}

type uint64ListBinder struct {
	val *[]uint64
}

func NewBoolArgBinder(p *bool) ArgBinder {
	return &boolArgBinder{val: p}
}

func (b *boolArgBinder) Bind(arg string, val string) error {
	if val != "" {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	}

	*(b.val) = true

	return nil
}

func NewFloat64ArgBinder(p *float64) ArgBinder {
	return &float64Binder{val: p}
}

func (b *float64Binder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	if f64val, err := strconv.ParseFloat(val, 64); err != nil {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	} else {
		*(b.val) = f64val
		return nil
	}
}

func NewFloat64ListArgBinder(p *[]float64) ArgBinder {
	return &float64ListBinder{val: p}
}

func (b *float64ListBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	var f64vals []float64

	for _, listval := range strings.Split(val, ",") {
		if f64val, err := strconv.ParseFloat(strings.TrimSpace(listval), 64); err != nil || listval == "" {
			return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
		} else {
			f64vals = append(f64vals, f64val)
		}
	}

	*(b.val) = f64vals

	return nil
}

func NewIntArgBinder(p *int) ArgBinder {
	return &intBinder{val: p}
}

func (b *intBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	if intval, err := strconv.Atoi(val); err != nil {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	} else {
		*(b.val) = intval
		return nil
	}
}

func NewIntListArgBinder(p *[]int) ArgBinder {
	return &intListBinder{val: p}
}

func (b *intListBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	var intvals []int

	for _, listval := range strings.Split(val, ",") {
		if intval, err := strconv.Atoi(listval); err != nil {
			return fmt.Errorf("invalid option-argument: '%s' for option: %s", listval, arg)
		} else {
			intvals = append(intvals, intval)
		}
	}

	*(b.val) = intvals

	return nil
}

func NewInt64ArgBinder(p *int64) ArgBinder {
	return &int64Binder{val: p}
}

func (b *int64Binder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	if i64val, err := strconv.ParseInt(val, 10, 64); err != nil {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	} else {
		*(b.val) = i64val
		return nil
	}
}

func NewInt64ListArgBinder(p *[]int64) ArgBinder {
	return &int64ListBinder{val: p}
}

func (b *int64ListBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	var i64vals []int64

	for _, listval := range strings.Split(val, ",") {
		if i64val, err := strconv.ParseInt(listval, 10, 64); err != nil {
			return fmt.Errorf("invalid option-argument: '%s' for option: %s", listval, arg)
		} else {
			i64vals = append(i64vals, i64val)
		}
	}

	*(b.val) = i64vals

	return nil
}

func NewStringArgBinder(p *string) ArgBinder {
	return &stringBinder{val: p}
}

func (b *stringBinder) Bind(_ string, val string) error {
	if val == "" {
		return nil
	}

	*(b.val) = val

	return nil
}

func NewStringListArgBinder(p *[]string) ArgBinder {
	return &stringListBinder{val: p}
}

func (b *stringListBinder) Bind(_ string, val string) error {
	if val == "" {
		return nil
	}

	var strvals []string

	for _, strval := range strings.Split(val, ",") {
		strvals = append(strvals, strval)
	}

	*(b.val) = strvals

	return nil
}

func NewUintArgBinder(p *uint) ArgBinder {
	return &uintBinder{val: p}
}

func (b *uintBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	if uintval, err := strconv.ParseUint(val, 10, 0); err != nil {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	} else {
		*(b.val) = uint(uintval)
		return nil
	}
}

func NewUintListArgBinder(p *[]uint) ArgBinder {
	return &uintListBinder{val: p}
}

func (b *uintListBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	var uintvals []uint

	for _, listval := range strings.Split(val, ",") {
		if uintval, err := strconv.ParseUint(listval, 10, 0); err != nil {
			return fmt.Errorf("invalid option-argument: '%s' for option: %s", listval, arg)
		} else {
			uintvals = append(uintvals, uint(uintval))
		}
	}

	*(b.val) = uintvals

	return nil
}

func NewUint64ArgBinder(p *uint64) ArgBinder {
	return &uint64Binder{val: p}
}

func (b *uint64Binder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	if u64val, err := strconv.ParseUint(val, 10, 64); err != nil {
		return fmt.Errorf("invalid option-argument: '%s' for option: %s", val, arg)
	} else {
		*(b.val) = u64val
		return nil
	}
}

func NewUint64ListArgBinder(p *[]uint64) ArgBinder {
	return &uint64ListBinder{val: p}
}

func (b *uint64ListBinder) Bind(arg string, val string) error {
	if val == "" {
		return nil
	}

	var u64vals []uint64

	for _, listval := range strings.Split(val, ",") {
		if u64val, err := strconv.ParseUint(listval, 10, 64); err != nil {
			return fmt.Errorf("invalid option-argument: '%s' for option: %s", listval, arg)
		} else {
			u64vals = append(u64vals, u64val)
		}
	}

	*(b.val) = u64vals

	return nil
}
