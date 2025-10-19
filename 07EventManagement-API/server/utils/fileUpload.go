package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUpload(c*gin.Context)(string, error){
	// get as "file" name from frontend
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	} 

	// create uploads folder
	err = os.MkdirAll("Uploads", os.ModePerm)
	if err != nil {
		return "", err
	}

	// create filename nd file path
	fileName := fmt.Sprint(time.Now().Unix()) + "_" + file.Filename
	filePath := filepath.Join("Uploads", fileName)

	// save changes 
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		return "" , err
	}

	return  filePath, err
}