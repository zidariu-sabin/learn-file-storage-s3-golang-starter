package main

import (
	"encoding/base64"
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

	//max number of bytes to story in memory per part in MultiPartForm of 10 MB
	const maxMemory = 10 << 20 //equal to 10 * 2^20
	err = r.ParseMultipartForm(maxMemory)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse request body", err)
	}
	//parsing thumbnail file from request body
	file, thumbnailHeader, err := r.FormFile("thumbnail")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
	}
	defer file.Close()

	//data media type from thumbnail file
	mediaType := thumbnailHeader.Header.Get("Content-Type")
	imageData, err := io.ReadAll(file)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to to read file data", err)
	}

	//updating video thumbnail url to file source
	videoMetadata, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch video metadata", err)
	}
	//check if the user updating the Thumbnail is the owner of the video
	if userID != videoMetadata.UserID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
	}

	//base64 encoding image so it can be stored in the database
	base64Thumbnail := base64.StdEncoding.EncodeToString(imageData)
	thumbnailURL := fmt.Sprintf("data:%v;base64,%v", mediaType, base64Thumbnail)
	videoMetadata.ThumbnailURL = &thumbnailURL
	err = cfg.db.UpdateVideo(videoMetadata)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Updating video thumbnail", err)
	}

	//returning updated video metadata
	respondWithJSON(w, http.StatusOK, videoMetadata)
}
