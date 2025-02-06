package validator

import (
	"strings"
	"testing"
)

// Define a struct for testing
type Person struct {
	Name  string `validate:"required" errormgs:"Name is required"`
	Email string `validate:"required,email" errormgs:"Invalid email"`
	Age   int    `validate:"gte=0,lte=130" errormgs:"Invalid age"`
}

func TestValidateFunc(t *testing.T) {
	validate := New()

	// Create an instance of the struct to validate
	person := Person{
		Name:  "",        // Invalid: Name is required
		Email: "invalid", // Invalid: Invalid email
		Age:   150,       // Invalid: Invalid age
	}

	// Call the ValidateFunc and check for expected errors
	err := ValidateFunc[Person](person, validate)

	// Check if there are validation errors
	if err == nil {
		t.Error("Expected validation errors, but got nil")
		return
	}

	// Check the error messages
	expectedErrs := []string{
		"Name: Name is required (required)",
		"Email: Invalid email (email)",
		"Age: Invalid age (lte)",
	}

	for _, expected := range expectedErrs {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error message '%s' not found", expected)
		}
	}
}

func TestValidateFuncV2(t *testing.T) {
	validate := New()

	// Create an instance of the struct to validate
	person := Person{
		Name:  "John Doe",            // Valid
		Email: "johndoe@example.com", // Valid
		Age:   30,                    // Valid
	}

	// Call the ValidateFunc and check for expected errors
	err := ValidateFunc[Person](person, validate)

	// Check if there are validation errors
	if err != nil {
		t.Errorf("Expected no validation errors, but got error: %s", err)
		return
	}

	// Test case where the validation fails with custom error message
	person.Name = "" // Invalid: Name is required

	err = ValidateFunc[Person](person, validate)

	// Check if there are validation errors
	if err == nil {
		t.Error("Expected validation errors, but got nil")
		return
	}

	// Check the error message
	expectedErr := "Name: Name is required (required)"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
	}

	// Test case where the validation fails without custom error message
	person.Email = "invalid" // Invalid: Invalid email
	person.Name = "a"
	err = ValidateFunc[Person](person, validate)

	// Check if there are validation errors
	if err == nil {
		t.Error("Expected validation errors, but got nil")
		return
	}

	// Check the error message
	expectedErr = "Email: Invalid email (email)"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
	}

	// Test case where the validation fails with panic
	person.Age = 150 // Invalid: Invalid age
	person.Email = "a@gmail.com"
	// Mock a panic during validation
	// defer func() {
	// 	if r := recover(); r == nil {
	// 		t.Error("Expected panic, but no panic occurred")
	// 	}
	// }()

	err = ValidateFunc[Person](person, validate)

	expectedErr = "Age: Invalid age (lte)"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
	}

}
