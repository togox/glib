package parser

import(
  	"github.com/aws/aws-sdk-go/aws"
  	"github.com/aws/aws-sdk-go/aws/session"
  	"github.com/aws/aws-sdk-go/service/s3/s3manager"
  	"github.com/aws/aws-sdk-go/aws/credentials"
  	"log"
	  "os"
    "github.com/go-resty/resty"
    "io/ioutil"
)
var (
	localPath string
	bucket    string
	prefix    string
)


type Options struct {
	AwsKey string
	AwsSecret    string
  File *os.File
  HtmlBucket string
  Url string
  FileName string
  ParserServerURL string
}
func Parse(options Options) (*s3manager.UploadOutput) {
   resp, err := resty.R().
   SetHeader("Content-Type", "application/json").
   SetBody(map[string]interface{}{"url": options.Url, "password": options.FileName}).
   Post(options.ParserServerURL)
   log.Println("resp:", resp)
   if(err != nil) {
     log.Println("Error parser html:", err)
     return nil
   } else {
      json := []byte(resp.String())
 		  writeErr := ioutil.WriteFile(options.FileName, json, 0644)
      if(writeErr != nil) {
        return nil
      } else {
        file, err := os.Open(options.FileName)
        bucket := options.HtmlBucket
       	creds := credentials.NewStaticCredentials(options.AwsKey, options.AwsSecret, "")
       	uploader := s3manager.NewUploader(session.New(&aws.Config{Credentials: 	  creds, Region: aws.String("us-west-1")}))
       	result, err := uploader.Upload(&s3manager.UploadInput{
       		Body:   file,
       		Bucket: aws.String(bucket),
       		Key:    aws.String(options.FileName),
       		ACL:    aws.String("public-read"),
       	})
       	if(err != nil) {
       		log.Println("Error upload file to S3 with message:", err)
       		return nil
       	}
         var removeErr = os.Remove(options.FileName)
         if (removeErr != nil) {
           log.Println("Error remove file ", options.FileName)
         }
       	return result
      }

   }
}

func Upload(options Options) (*s3manager.UploadOutput) {
	bucket := options.HtmlBucket
	creds := credentials.NewStaticCredentials(options.AwsKey, options.AwsSecret, "")
	uploader := s3manager.NewUploader(session.New(&aws.Config{Credentials: 	  creds, Region: aws.String("us-east-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   options.File,
		Bucket: aws.String(bucket),
		Key:    aws.String(options.FileName),
		ACL:    aws.String("public-read"),
	})
	if(err != nil) {
		log.Println("Error upload file to S3 with message:", err)
		return nil
	}
  var removeErr = os.Remove(options.FileName)
  if (removeErr != nil) {
    log.Println("Error remove file ", options.FileName)
  }
	return result
}
