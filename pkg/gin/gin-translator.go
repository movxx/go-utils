package gin

import (
	"os"
	"reflect"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/movxx/go-utils/pkg/mapslice"
	"github.com/pkg/errors"
)

var trans ut.Translator

func init() {
	locale := strings.ToLower(os.Getenv("LOCALE"))
	if locale == "" {
		locale = "en"
	}
	InitTrans(locale)
}

func CombineErrors(err error) error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return err
	}
	// we join the multi fields validate msg to one
	validateMsg := slice.Join(mapslice.MapValue2Slice(removeTopStruct(errs.Translate(trans))), ",")
	return errors.New(validateMsg)
}

func addValueToMap(fields map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		if len(fieldArr) > 1 {
			NewFields := map[string]string{fieldArr[1]: err}
			returnMap := addValueToMap(NewFields)
			if res[fieldArr[0]] != nil {
				for k, v := range returnMap {
					res[fieldArr[0]].(map[string]interface{})[k] = v
				}
			} else {
				res[fieldArr[0]] = returnMap
			}
			continue
		} else {
			res[field] = err
			continue
		}
	}
	return res
}

// removeTopStruct remove the struct name of fields.
func removeTopStruct(fields map[string]string) map[string]interface{} {
	lowerMap := map[string]string{}
	for field, err := range fields {
		fieldArr := strings.SplitN(field, ".", 2)
		lowerMap[fieldArr[1]] = err
	}
	res := addValueToMap(lowerMap)
	return res
}

// InitTrans init the translator,
func InitTrans(locale string) {
	// modify customized validator of gin framework
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	// register a function for get json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	zhT := zh.New()
	enT := en.New()

	// we add multi translator and set the fallback translator with en
	uni := ut.New(enT, zhT, enT)

	// we should extract locale from http header of 'Accept-Language',
	// find the translator by FindTranslator with multi locale also be ok.
	trans, _ = uni.GetTranslator(locale)

	// register translator with locale, we don't care of the error.
	switch locale {
	case "en":
		enTrans.RegisterDefaultTranslations(v, trans)
	case "zh":
		zhTrans.RegisterDefaultTranslations(v, trans)
	default:
		enTrans.RegisterDefaultTranslations(v, trans)
	}
}
