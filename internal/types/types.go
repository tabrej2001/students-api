package types

type Students struct{
	Name string  `validate:"required"`
	Email string `validate:"required"`
	Age int 	 `validate:"required"`
}