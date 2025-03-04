package pkg

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

func NewInternalServerError(err error) Error {
	return Error{
		Message: "internal server error",
		Error:   err.Error(),
		Status:  500,
	}
}

func NewInvalidCredentialsError() Error {
	return Error{
		Message: "Invalid Credentials",
		Status:  401,
	}
}

func NewBadRequestError(err error) Error {
	return Error{
		Message: "Bad Request",
		Error:   err.Error(),
		Status:  400,
	}
}

func NewEntityAlreadyExistsError(entity string) Error {
	return Error{
		Message: "Unprocessable entity",
		Status:  409,
		Error:   entity + " already exists!",
	}
}

func NewNotFoundError(entity string) Error {
	return Error{
		Message: entity + " not found",
		Status:  404,
	}
}

func NewMissingFieldError(field string) Error {
	return Error{
		Message: field + " is required",
		Status:  400,
	}
}
