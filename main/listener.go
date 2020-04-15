package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)
import "github.com/valyala/fasthttp"

type Data struct {
	Text string `json:"text"`
	ContentId int `json:"content_id"`
	ClientId int `json:"client_id"`
	Timestamp int `json:"timestamp"`
}

func main() {


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_KEY")

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3BucketName := s3Bucket
	s3Client, err_ := minio.New("s3.amazonaws.com", accessKeyId, secretKey, true)
	if err_ != nil {
		log.Fatalln(err_)
	}

	requestHandler := func(ctx *fasthttp.RequestCtx)  {
		// set some headers and status code first
		ctx.SetStatusCode(fasthttp.StatusOK)

		var value = ctx.PostBody()
		s := string(value[:len(value)])

		fmt.Println(s)
		in := value[:len(value)]
		var data *Data;
		_err := json.Unmarshal(in, &data)
		if (_err != nil) {
			panic(_err)
		}

		if err_ == nil {
			go putFile(s3BucketName, s3Client, data, in)
		}

		// then override already written body
		ctx.SetBody([]byte("msg OK"))

		// then update status code
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

	print("Start listening on fixed port 8082")
	fasthttp.ListenAndServe(":8082", requestHandler)
}

func putFile(bucketName string, s3Client *minio.Client, data *Data, content []byte)  {
	obj := bytes.NewReader(content)
	fileName := "content_logs_" + _convertTimestampToISO(int64(data.Timestamp)) + strconv.FormatInt(int64(data.ClientId), 10)
	n, err := s3Client.PutObject(bucketName, fileName,
		obj, obj.Size(), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			ContentEncoding: "gzip",
		})
	if err != nil {
		log.Println(err)
	}

	log.Println("Uploaded", fileName, " of size: ", n, "Successfully.")
}

func _convertTimestampToISO(value int64) string  {
	i, err := strconv.ParseInt(strconv.Itoa(int(value)), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02")
}