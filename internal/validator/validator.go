package validator

type Validator struct {
	errors map[string]string
}

type Errors map[string]string

func New() *Validator {
	return &Validator{make(Errors)}
}

func (v *Validator) AddError(key, errorMsg string) {
	if _, ok := v.errors[key]; !ok {
		v.errors[key] = errorMsg
	}
}

func (v *Validator) Check(condition bool, key, errorMsg string) {
	if !condition {
		v.AddError(key, errorMsg)
		return
	}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() Errors {
	return v.errors
}
