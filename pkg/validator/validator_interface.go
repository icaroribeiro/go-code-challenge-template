package validator

// IValidator transport/http is the validator's contract.
type IValidator interface {
	Validate(i interface{}) error
	ValidateWithTags(i interface{}, tags string) error
}
