package echo_handle

import (
	"errors"
	"silo/pkg/errs"
	"silo/pkg/utils"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

// NewValidator NewValidator
func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())

	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")

	validate = validator.New()
	err := entranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil
	}

	return &Validator{
		translator: trans,
		validator:  validate,
	}
}

// Validator .
type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// Validate Validate
func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		var e validator.ValidationErrors
		errors.As(err, &e)
		t := e.Translate(v.translator)
		return errs.NewBadRequestError(strings.Join(utils.MapValues(t), ";"))
	}
	return nil
}
