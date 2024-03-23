package repos

import (
	"fmt"
	"time"
)

type Int struct {
	value string
}

type Time struct {
	value string
}

type IsValue struct {
	value string
}

func Public() IsValue  { return IsValue{"public"} }
func Private() IsValue { return IsValue{"private"} }

func (Int) Range(from, to int) Int {
	return Int{fmt.Sprintf("%d..%d", from, to)}
}

func (Time) Range(from, to time.Time) Time {
	return Time{fmt.Sprintf("%s..%s", from.Format("2006-01-02"), to.Format("2006-01-02"))}
}

func (Int) Min(min int) Int {
	return Int{fmt.Sprintf(">=%d", min)}
}

func (Time) Min(min time.Time) Time {
	return Time{fmt.Sprintf(">=%s", min.Format("2006-01-02"))}
}

func (Int) MinEx(min int) Int {
	return Int{fmt.Sprintf(">%d", min)}
}

func (Time) MinEx(min time.Time) Time {
	return Time{fmt.Sprintf(">%s", min.Format("2006-01-02"))}
}

func (Int) Max(max int) Int {
	return Int{fmt.Sprintf("<=%d", max)}
}

func (Time) Max(max time.Time) Time {
	return Time{fmt.Sprintf("<=%s", max.Format("2006-01-02"))}
}

func (Int) MaxEx(max int) Int {
	return Int{fmt.Sprintf("<%d", max)}
}

func (Time) MaxEx(max time.Time) Int {
	return Int{fmt.Sprintf("<%s", max.Format("2006-01-02"))}
}

func (Int) Eq(val int) Int {
	return Int{fmt.Sprintf("=%d", val)}
}

func (Time) Eq(val time.Time) Int {
	return Int{fmt.Sprintf("=%s", val.Format("2006-01-02"))}
}

type Order struct {
	Value string
}

func Desc() *Order { return &Order{"desc"} }
func Asc() *Order  { return &Order{"asc"} }

type Sort struct {
	Value string
}

func SortByStars() *Sort {
	return &Sort{"stars"}
}
