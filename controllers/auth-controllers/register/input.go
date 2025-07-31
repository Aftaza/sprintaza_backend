package registerAuth

// GoogleOAuthRegisterInput represents the input data from Google OAuth
type GoogleOAuthRegisterInput struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
	Name  string `json:"name" binding:"required" validate:"required,min=2,max=100"`
}

// Validate validates the input data
func (input *GoogleOAuthRegisterInput) Validate() error {
	if input.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if input.Name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if len(input.Name) < 2 {
		return &ValidationError{Field: "name", Message: "Name must be at least 2 characters long"}
	}
	if len(input.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Name must not exceed 100 characters"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}