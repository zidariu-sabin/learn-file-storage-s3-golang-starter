package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	//max number of bytes to story in memory per part in MultiPartForm of 10 MB
	const maxMemory = 10 << 20 //equal to 10 * 2^20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("thumbnail")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
	}

	mediaType := header.Header.Get("Content-Type")

	imageData, err := io.ReadAll(file)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to to read file data", err)
	}

	videoMetadata, err := cfg.db.GetVideo(videoID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch video metadata", err)
	}

	if userID != videoMetadata.UserID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
	}

	var thumbnail thumbnail = // see how to parse map here

	respondWithJSON(w, http.StatusOK, struct{}{})
}
