package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const ParamError = 1000  // 参数参数错误
const UnKnowError = 9999 // 未知错误

var ErrorCodeMap = map[string]int{
	"verifyCodeError": 1001, "userNotFound": 1002, "passwordError": 1003, "phoneLimit": 1004, "unauthorized": 1005,
}

var ErrorMsgMap = map[int]string{
	1001: "验证码错误", 1002: "用户没找到", 1003: "用户名或密码错误", 1004: "被限制的手机号", 1005: "未经授权的",
}

var ErrorMsgMapEn = map[int]string{
	1001: "verify code error", 1002: "user not found", 1003: "user or password error", 1004: "phone limit", 1005: "unauthorized",
}

func ErrorWarper(err error) gin.H {
	if err != nil {
		s := err.Error()
		code, ok := ErrorCodeMap[s] // 自定义异常信息
		if ok {
			return gin.H{"error": ErrorMsgMap[code], "code": code}
		}
		if es, ok := err.(validator.ValidationErrors); ok {
			return gin.H{"error": warpValidatorError(es), "code": ParamError}
		}
		return gin.H{"error": s, "code": UnKnowError}
	}
	return gin.H{"code": 0, "error": ""}
}

func warpValidatorError(es validator.ValidationErrors) map[string]string {
	errs := map[string]string{}
	for _, v := range es {
		errs[v.Field()] = v.Param()
	}
	return errs
}
