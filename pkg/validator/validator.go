package validator

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator

func InitTrans() (err error) {
	// 获取gin的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() // 中文翻译器
		uni := ut.New(zhT, zhT)
		Trans, _ = uni.GetTranslator("zh")
		// 注册翻译器
		err = zh_translations.RegisterDefaultTranslations(v, Trans)

		// 注册一个函数，获取struct tag中的字段备注作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label") // 优先使用label标签
			if name == "" {
				name = fld.Tag.Get("json") // 没有label标签就用json标签
			}
			return name
		})

		return
	}
	return
}

// 翻译错误信息
func TranslateErr(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		return validationErrs[0].Translate(Trans)
	}
	return err.Error()
}
