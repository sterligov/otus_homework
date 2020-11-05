// Code generated by cool valid tool; DO NOT EDIT.
package testdata

import "regexp"

type ValidationError struct {
	Field string
	Err   string
}

func (this User) Validate() ([]ValidationError, error) {
	ve := make([]ValidationError, 0)

	if len(this.ID) != 36 {
		err := ValidationError{
			Field: "ID",
			Err:   "Len must be equal to 36",
		}
		ve = append(ve, err)
	}

	if this.Age < 18 {
		err := ValidationError{
			Field: "Age",
			Err:   "Value must be greater or equal than 18",
		}
		ve = append(ve, err)
	}

	if this.Age > 50 {
		err := ValidationError{
			Field: "Age",
			Err:   "Value must be less or equal to 50",
		}
		ve = append(ve, err)
	}

	if match, err := regexp.MatchString("^\\w+@\\w+\\.\\w+$", this.Email); err != nil {
		return nil, err
	} else if !match {
		err := ValidationError{
			Field: "Email",
			Err:   "Value must satisfy regular expression ^\\w+@\\w+\\.\\w+$",
		}
		ve = append(ve, err)
	}

	if "admin" != this.Role && "stuff" != this.Role {
		err := ValidationError{
			Field: "Role",
			Err:   "Value must be in list [admin,stuff]",
		}
		ve = append(ve, err)
	}

	for _, v := range this.Phones {
		if len(v) != 11 {
			err := ValidationError{
				Field: "Phones",
				Err:   "Len must be equal to 11",
			}
			ve = append(ve, err)
		}
	}

	return ve, nil
}

func (this App) Validate() ([]ValidationError, error) {
	ve := make([]ValidationError, 0)

	if len(this.Version) != 5 {
		err := ValidationError{
			Field: "Version",
			Err:   "Len must be equal to 5",
		}
		ve = append(ve, err)
	}

	return ve, nil
}

func (this Response) Validate() ([]ValidationError, error) {
	ve := make([]ValidationError, 0)

	if 200 != this.Code && 404 != this.Code && 500 != this.Code {
		err := ValidationError{
			Field: "Code",
			Err:   "Value must be in list [200,404,500]",
		}
		ve = append(ve, err)
	}

	return ve, nil
}
