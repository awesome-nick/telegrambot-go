package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var imgPath string

func UploadHandler(imgPath string) gin.HandlerFunc {
	f := func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			c.JSON(http.StatusAccepted, err)
			return
		}

		err = c.SaveUploadedFile(file, imgPath+file.Filename)
		if err != nil {
			c.JSON(http.StatusAccepted, err)
			return
		}

		c.JSON(http.StatusOK, "Image "+file.Filename+" successfully uploaded")
	}
	return gin.HandlerFunc(f)
}
