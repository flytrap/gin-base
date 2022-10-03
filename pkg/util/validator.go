package util

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Mobile 验证手机号码
func MobileValidator(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^(13|14|15|17|18|19)[0-9]{9}$`, fl.Field().String())
	return ok
}

// IdCard 验证身份证号码
func IdCardValidator(fl validator.FieldLevel) bool {
	id := fl.Field().String()

	var a1Map = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}

	var idStr = strings.ToUpper(string(id))
	var reg, err = regexp.Compile(`^[0-9]{17}[0-9X]$`)
	if err != nil {
		return false
	}
	if !reg.Match([]byte(idStr)) {
		return false
	}
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, err := strconv.Atoi(string(c)); err == nil {
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = a1Map[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}
	return a1Str == signChar
}

// 检测生日
func CheckBirthDate(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if ok {
		// 当前时间应该大于生日时间
		if time.Now().After(t) {
			return true
		}
	}
	return false
}
