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
	r.POST("/upload", func(c *gin.Context) {
		// Get file
		file, header, _ := c.Request.FormFile("file")
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
		svc := s3.New(sess)
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET")),
			Key:    aws.String(filename),
			Body:   file,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		cdnUrl := os.Getenv("CDN_URL")

		c.JSON(200, gin.H{
			"message": "upload",
			"cdnUrl":  cdnUrl + "/" + filename,
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port) // listen and serve on
}
