package users

import (
	"CTFBattle/db"
	"encoding/json"
	"net/http"
	"time"

	"gofr.dev/pkg/gofr"
)

type User struct {
	UserID    string    `json:"userid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// AddUser adds a new user to the database
// AddUser adds a new user to the database, but only if the user does not already exist
func AddUser(c *gofr.Context, user *User) error {
	client := db.GetClient()

	// Check if the user already exists by email
	existingUser, err := GetUser(user.Email)
	if err == nil && existingUser != nil {
		// If the user already exists, do not add them again
		return nil
	}

	// Set the user creation time
	user.CreatedAt = time.Now()

	// Prepare the data to insert
	data := map[string]interface{}{
		"userid":     user.UserID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	// Insert new user into the database
	_, _, err = client.From("users").Insert(data, false, "","","").Execute()
	if err != nil {
		return err
	}

	return nil
}

// GetUser retrieves a user by email
func GetUser(email string) (*User, error) {
	client := db.GetClient()
	var user User
	
	result, _, err := client.From("users").
		Select("*", "", false).
		Eq("email", email).
		Single().
		Execute()
	
	if err != nil {
		return nil, err
	}
	
	// Parse the result into User struct
	err = json.Unmarshal(result, &user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// DeleteUser deletes a user by email
func DeleteUser(email string) error {
	client := db.GetClient()
	
	_, _, err := client.From("users").Delete("*", "").
		Eq("email", email).
		Execute()
	
	if err != nil {
		return err
	}
	
	return nil
}

func AddUserHandler(ctx *gofr.Context) (interface{}, error) {
	var user User

	// Parse the request body into the user struct
	err := ctx.Bind(&user)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "Invalid request body")
		return nil, err
	}

	err = AddUser(ctx, &user)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to add user")
		return nil, err
	}

	return map[string]string{"message": "User added successfully"}, nil
}

// GetUserHandler handles the endpoint to get user details by email
func GetUserHandler(ctx *gofr.Context) (interface{}, error) {
	email := ctx.PathParam("email")

	user, err := GetUser(email)
	if err != nil {
		ctx.Error(http.StatusNotFound, "User not found")
		return nil, err
	}

	return user, nil
}

// DeleteUserHandler handles the endpoint to delete a user by email
func DeleteUserHandler(ctx *gofr.Context) (interface{}, error) {
	email := ctx.PathParam("email")

	err := DeleteUser(email)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to delete user")
		return nil, err
	}

	return map[string]string{"message": "User deleted successfully"}, nil
}
