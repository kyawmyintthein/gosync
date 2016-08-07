package s3

import(
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awsutil"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/session"
    "gosync/logger"
    "fmt"
)

type Config struct{
    awsConfig aws.Config
    BucketName  string 
}

var (
    config  
    s3Client s3.S3
)

func Init(config Config){
   s3Client = S3.New(&confg.awsConfig)
}

func Sync(filepath string){
    var(
        fileBytes []byte
        err error
        fileType string
        size int64
    )
    fileBytes, size, fileType, err = readFile(filepath)
    if err != nil{
        logger.Error(fmt.Printf("ERROR %s", err))
    }
    absPath, _ := filepath.Abs(filepath)
    params := &s3.PutObjectInput{
        Bucket:        aws.String(s3Client.BucketName),
        Key:           aws.String(absPath),
        Body:          fileBytes,
        ContentLength: aws.Int64(size),
        ContentType:   aws.String(fileType),
    }

    updateToS3(params); err != nil{
}

// read file as bytes array
func readFile(filepath string) ([]byte, string, int64, error){
    var(
        bytes []byte
        err error
        fileType string = ""
        size  int64
    )

    file, err = os.Open(filepath)
    if err != nil {
        return bytes, size, fileType, err
    }
    defer file.Close()

    fileInfo, _ := file.Stat()
    size = fileInfo.Size()
    buffer := make([]byte, size)
    file.Read(buffer)
    bytes := bytes.NewReader(buffer)
    fileType = http.DetectContentType(buffer)
    return bytes, size, fileType, err
}

// upload to s3
func updateToS3(params s3.PutObjectInput){
    resp, err := s3Client.PutObject(params)
    logger.Info(fmt.Printf("response %s", awsutil.StringValue(resp)))
    if err != nil{
        logger.Warn(fmt.Printf("response %s", awsutil.StringValue(resp)))
        logger.Error(fmt.Printf("ERROR %s", err))
    }
}

