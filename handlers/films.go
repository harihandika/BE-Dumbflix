package handlers

import (
	filmsdto "dumbflix/dto/films"
	dto "dumbflix/dto/result"
	"dumbflix/models"
	"dumbflix/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"

	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type handlerFilm struct {
	FilmRepository repositories.FilmRepository
}

func HandlerFilm(FilmRepository repositories.FilmRepository) *handlerFilm {
	return &handlerFilm{FilmRepository}
}

func (h *handlerFilm) FindFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	films, err := h.FilmRepository.FindFilm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: films}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilm) GetFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	film, err := h.FilmRepository.GetFilm(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilm(film)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilm) CreateFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	fmt.Println(userId)
	dataContex := r.Context().Value("dataFile")
	filepath := dataContex.(string)

	category_id, _ := strconv.Atoi(r.FormValue("category_id"))

	request := filmsdto.CreateFilmRequest{
		Title:      r.FormValue("title"),
		Year:       r.FormValue("year"),
		CategoryID: category_id,
		Thumbnail:  filepath,
		Desc:       r.FormValue("desc"),
		LinkFilm:   r.FormValue("link"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbflix"})

	if err != nil {
		fmt.Println(err.Error())
	}

	film := models.Film{
		Title:         request.Title,
		ThumbnailFilm: resp.SecureURL,
		Year:          request.Year,
		CategoryID:    request.CategoryID,
		Desc:          request.Desc,
		LinkFilm:      request.LinkFilm,
	}

	data, err := h.FilmRepository.CreateFilm(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	film, _ = h.FilmRepository.GetFilm(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilm(film)}
	json.NewEncoder(w).Encode(response)
}

// func (h *handlerFilm) UpdateFilm(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	dataContex := r.Context().Value("dataFile")
// 	filepath := dataContex.(string)

// 	var ctx = context.Background()
// 	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
// 	var API_KEY = os.Getenv("API_KEY")
// 	var API_SECRET = os.Getenv("API_SECRET")

// 	// Add your Cloudinary credentials ...
// 	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

// 	// Upload file to Cloudinary ...
// 	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbflix"})

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	request := filmsdto.UpdateFilmRequest{
// 		Title:     r.FormValue("title"),
// 		Year:      r.FormValue("year"),
// 		Desc:      r.FormValue("desc"),
// 		Thumbnail: resp.SecureURL,
// 		LinkFilm:  r.FormValue("link"),
// 	}
// 	// request := new(filmsdto.UpdateFilmRequest)
// 	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 	// 	w.WriteHeader(http.StatusBadRequest)
// 	// 	response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
// 	// 	json.NewEncoder(w).Encode(response)
// 	// 	return
// 	// }

// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])
// 	film, err := h.FilmRepository.GetFilm(int(id))
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	if request.Title != "" {
// 		film.Title = request.Title
// 	}

// 	if request.Thumbnail != "" {
// 		film.ThumbnailFilm = request.Thumbnail
// 	}

// 	if request.Year != "" {
// 		film.Year = request.Year
// 	}

// 	if request.Desc != "" {
// 		film.Desc = request.Desc
// 	}

// 	if request.LinkFilm != "" {
// 		film.LinkFilm = request.LinkFilm
// 	}

// 	data, err := h.FilmRepository.UpdateFilm(film)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilm(data)}
// 	json.NewEncoder(w).Encode(response)
// }

func (h *handlerFilm) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	film, err := h.FilmRepository.GetFilm(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.FilmRepository.DeleteFilm(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilm(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseFilm(u models.Film) filmsdto.FilmResponse {
	return filmsdto.FilmResponse{
		ID:         u.ID,
		Title:      u.Title,
		Thumbnail:  u.ThumbnailFilm,
		Category:   u.Category,
		CategoryID: u.CategoryID,
		Year:       u.Year,
		Desc:       u.Desc,
		LinkFilm:   u.LinkFilm,
	}
}
