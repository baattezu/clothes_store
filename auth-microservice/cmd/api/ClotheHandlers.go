package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nurtikaga/internal/data"
	"github.com/nurtikaga/internal/validator"
)

func (app *application) createClotheHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Module_name     string `json:"module_name"`
		Module_duration int32  `json:"module_duration"`
		ExamType        string `json:"exam_type"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	clothes := &data.ClotheInfo{
		ModuleName:     input.Module_name,
		ModuleDuration: input.Module_duration,
		ExamType:       input.ExamType,
	}
	//v := validator.New()
	/*if data.ValidateModule(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}*/
	err = app.clothes.ClotheInfo.Insert(clothes)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/clothes/%d", clothes.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"clothes info": clothes}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getClotheHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	clothe, err := app.clothes.ClotheInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"clothes": clothe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) editClotheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	clothe, err := app.clothes.ClotheInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Module_name     string `json:"module_name"`
		Module_duration int32  `json:"module_duration"`
		ExamType        string `json:"exam_type"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	clothe.ModuleName = input.Module_name
	clothe.ModuleDuration = input.Module_duration
	clothe.ExamType = input.ExamType

	err = app.clothes.ClotheInfo.Update(clothe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module": clothe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteClotheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.clothes.ClotheInfo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Clothe successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClotheHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Module_name string
		ExamType    string
		Page        int
		PageSize    int
		Sort        string
	}
	v := validator.New()

	qs := r.URL.Query()
	input.Module_name = app.readString(qs, "title", "")
	input.ExamType = app.readString(qs, "title", "")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)

}
