package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	apiRouter := os.Getenv("API_ROUTER")
	if apiRouter == "" {
		apiRouter = "/upload"
	}

	r.POST(apiRouter, upload)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = r.Run(":" + port) // listen and serve on
}

func upload(c *gin.Context) {
	// Get file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"message": "file error",
		})
		return
	}
	if file == nil {
		c.JSON(500, gin.H{
			"message": "file is nil",
		})
		return
	}
	fmt.Println(header.Filename)
	filename := header.Filename

	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("REGION")),
		Endpoint: aws.String(os.Getenv("ENDPOINT")),
		// accessKey, secretKey
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("ACCESS_KEY_ID"),
			os.Getenv("SECRET_ACCESS_KEY"),
			"",
		),
	}))

	pathPrefix := os.Getenv("PATH_PREFIX")

	key := filename
	if pathPrefix != "" {
		key = pathPrefix + "/" + filename
	}

	fmt.Println("pathPrefix", pathPrefix)

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET")),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"message": "upload error",
		})
		return
	}
	cdnUrl := os.Getenv("CDN_URL")

	c.JSON(200, gin.H{
		"message": "upload",
		"url":     cdnUrl + "/" + key,
	})
}
