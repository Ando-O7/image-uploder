package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Upload upload files.
func Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	// get uuid.
	uuid := c.PostForm("uuid")

	for _, file := range files {
		// Load uuid for file name.
		err := c.SaveUploadedFile(file, "images/"+uuid+".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
	}

	c.JSON(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

// Delete remove file.
func Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	err := os.Remove(fmt.Sprintf("images/%s.png", uuid))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("id: %s is deleted!", uuid)})
}

// file info
type File struct {
	Path string `json:"path"`
	Size int64 	`json:"size"`
}

func dirwalk(dir string) (files []File, err error) {
    err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        path = strings.Replace(path, "images/", "http://localhost:8888/", 1)
        size := info.Size()
        f := File{
            Path: path,
            Size: size,
        }
        files = append(files, f)
        return nil
    })
    if err != nil {
        return
    }
    files = files[1:]
    return
}

// List return url & size list
func List(c *gin.Context) {
    files, err := dirwalk("./images")
    if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, files)
}