package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	Gender         string     `validate:"oneof=male female prefer_not_to"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

type SignupRequest struct {
	Email           string `validate:"required,email"`
	Password        string `validate:"required"`
	PasswordConfrim string `validate:"eqfield=Password"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {

	validate = validator.New(validator.WithRequiredStructEnabled())

	validateStruct()
	validateVariable()
	validateSignup()
}

func validateStruct() {

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Gender:         "male",
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses:      []*Address{address},
	}

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
}

func validateVariable() {

	myEmail := "joeybloggs.gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}

	// email ok, move on
}

// validateSignup demonstrates the process of validating a signup request
// using a predefined struct validation. It creates a sample request with
// invalid data to showcase error handling and a second request with valid
// data to ensure successful validation.
func validateSignup() {
	fmt.Println("Validate Signup request")

	// Create a SignupRequest instance with invalid data
	req := &SignupRequest{
		Email:           "test@test.com",
		Password:        "Password123!",
		PasswordConfrim: "badpassword",
	}

	// Validate the SignupRequest instance against defined struct validation rules
	err := validate.Struct(req)
	if err != nil {
		// Log the validation failure and the specific error details
		fmt.Printf("Signup request validation failed: %v\n", err)
	}

	// Create a new SignupRequest instance with corrected data
	req = &SignupRequest{
		Email:           "test@test.com",
		Password:        "Password123!",
		PasswordConfrim: "Password123!",
	}

	// Revalidate the corrected SignupRequest instance
	err = validate.Struct(req)
	if err != nil {
		// If this code path executes, there is an unexpected error, so panic
		panic(0) // Should not reach here in normal circumstances
	}

	// Log successful validation of the signup request
	fmt.Println("Signup request has been validated")
}
