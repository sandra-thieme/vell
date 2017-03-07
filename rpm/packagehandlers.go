package rpm

import (
	"net/http"
	"github.com/gorilla/mux"
)

// POST /repositories/{name}/packages
func AddRPM(w http.ResponseWriter, r *http.Request) {
	repo := YumRepository{mux.Vars(r)["name"]}
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		fail(w, err)
	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				fail(w, err)
			}
			defer src.Close()

			repo.add(file.Filename, src)
		}
	}
	repo.update()
	w.WriteHeader(http.StatusCreated)
}
