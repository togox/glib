package uploader

import(
  	"github.com/aws/aws-sdk-go/aws"
  	"github.com/aws/aws-sdk-go/aws/session"
  	"github.com/aws/aws-sdk-go/service/s3/s3manager"
  	"github.com/aws/aws-sdk-go/aws/credentials"
  	"log"
	  "os"
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
  FileName string
  HtmlBucket string
}
func Upload(options Options) (*s3manager.UploadOutput) {
	bucket := options.HtmlBucket
	creds := credentials.NewStaticCredentials(options.AwsKey, options.AwsSecret, "")
	uploader := s3manager.NewUploader(session.New(&aws.Config{Credentials: 	  creds, Region: aws.String("us-west-1")}))
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
