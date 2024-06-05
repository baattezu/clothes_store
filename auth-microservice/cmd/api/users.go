package main

import (
	"errors"
	"github.com/nurtikaga/internal/data"
	"github.com/nurtikaga/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &data.UserInfo{
		Name:      input.Name,
		Surname:   input.Surname,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Add the "movies:read" permission for the new user.

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")

			app.failedValidationResponse(w, r, v.Errors)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request URL query parameters.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Call the GetUserByID method to fetch the user information from the database.
	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send the user information as a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Surname string
		Email   string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	input.Email = app.readString(qs, "email", "")
	// Send the users information as a JSON response.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "surname", "email", "-id", "-name", "-surname", "email"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user_info, err := app.models.Users.GetAll(input.Name, input.Surname, input.Email, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user_info": user_info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) editUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request URL query parameters.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Parse the request body into the updated user struct.
	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	type passwordInput struct {
		Plaintext *string `json:"plaintext"`
		Hash      []byte  `json:"-"`
	}

	var input struct {
		Name     *string        `json:"name"`
		Surname  *string        `json:"surname"`
		Email    *string        `json:"email"`
		Password *passwordInput `json:"password"`
	}

	// Set the ID of the updated user to match the ID extracted from the URL.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Call the UpdateUser method to update the user information in the database.
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Surname != nil {
		user.Surname = *input.Surname
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		if input.Password.Plaintext != nil {
			err := user.Password.Set(*input.Password.Plaintext)
			if err != nil {
				// Handle error, perhaps return an error response
			}
		}
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		// Send a JSON response indicating that the user information was successfully updated.
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "User updated successfully"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
func (app *application) deleteUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request URL query parameters.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Call the DeleteUser method to delete the user from the database.
	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a JSON response indicating that the user was successfully deleted.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "User deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
