package forms

type errors map[string][]string

func (e *errors) Add(field, message string) {
	(*e)[field] = append((*e)[field], message)
}

// Return the first error message
func (e *errors) Get(field string) string {
	es := (*e)[field]
	if len(es) == 0 {
		return ""
	}
	return (*e)[field][0]
}
