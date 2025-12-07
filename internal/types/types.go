package types

type Student struct {
	Id    int64  `json:"Id"`
	Name  string `json:"Name" validate:"required"`
	Email string `json:"Email" validate:"required"`
	Age   int    `json:"Age" validate:"required"`
}

//no validation require here
type StudentPatch struct {
	Name  *string `json:"Name"`
	Email *string `json:"Email"`
	Age   *int    `json:"Age"`
}
