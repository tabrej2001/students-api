package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/tabrej2001/student-api/internal/storage"
	"github.com/tabrej2001/student-api/internal/types"
	"github.com/tabrej2001/student-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Students

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err.Error())
		}

		// request validation

		if err := validator.New().Struct(student); err != nil {
			validationError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationError))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("User created successfully", slog.String("userId", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetStudentByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err)) 
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("error gettig user", slog.String("id", id ))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
