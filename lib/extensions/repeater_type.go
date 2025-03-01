package extensions

type RepeaterType struct {
	display       string
	hasParameters bool
}

func (repeaterType RepeaterType) IsEqualWith(compared RepeaterType) bool {
	return repeaterType.display == compared.display
}

type RepeaterTypes []RepeaterType

func (types RepeaterTypes) Some(checker func(value RepeaterType) bool) bool {
	for _, repeaterType := range types {
		if checker(repeaterType) {
			return true
		}
	}

	return false
}
