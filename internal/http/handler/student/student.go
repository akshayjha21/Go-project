package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/akshayjha21/Student-Api/internal/storage"
	"github.com/akshayjha21/Student-Api/internal/types"
	response "github.com/akshayjha21/Student-Api/internal/utils"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		//so we are reading the response and decoding it

		//->it might also throw an error
		slog.Info("Creating a student")
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			//now we will return a json response
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//REQUEST VALIDATOR
		if err := validator.New().Struct((student)); err != nil {
			validateErrs := err.(validator.ValidationErrors) //we are type casting it to the given argument type present in the error func in response.go i.e validator.ValidationErrors,if not then it will throw an error

			response.WriteJson(w, http.StatusBadRequest, response.ValidatorError(validateErrs))
			return
		}

		lastid, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("User created successfully", slog.String("UserId", fmt.Sprint(lastid)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastid})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intid)
		if err != nil {
			slog.Error("error getting user ", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
