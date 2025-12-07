package types

type Student struct {
	Id    int64  `json:"Id"`
	Name  string `json:"Name" validate:"required"`
	Email string `json:"Email" validate:"required"`
	Age   int    `json:"Age" validate:"required"`
}
