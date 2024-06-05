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
		Cloth_name string `json:"cloth_name"`
		Cloth_cost int32  `json:"cloth_cost"`
		Cloth_size string `json:"cloth_size"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	clothes := &data.ClotheInfo{
		ClothName: input.Cloth_name,
		ClothCost: input.Cloth_cost,
		ClothSize: input.Cloth_size,
	}

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
		Cloth_name string `json:"cloth_name"`
		Cloth_cost int32  `json:"cloth_cost"`
		Cloth_size string `json:"cloth_size"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	clothe.ClothName = input.Cloth_name
	clothe.ClothCost = input.Cloth_cost
	clothe.ClothSize = input.Cloth_size

	err = app.clothes.ClotheInfo.Update(clothe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"cloth": clothe}, nil)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Clothe successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClotheHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Cloth_name string
		Cloth_cost string
		Page       int
		PageSize   int
		Sort       string
	}
	v := validator.New()

	qs := r.URL.Query()
	input.Cloth_name = app.readString(qs, "cloth_name", "")
	input.Cloth_cost = app.readString(qs, "cloth_cost", "")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}
