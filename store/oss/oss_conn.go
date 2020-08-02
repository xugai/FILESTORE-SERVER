package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client
// LTAI4GBwyJYXh2NWRpBr9fp8  S0cOobEbehrkpMSCXSoqYzwRqmhCT7
const (
	ossEndpoint = "oss-cn-shenzhen.aliyuncs.com"
	accessKeyID = "LTAI4GBwyJYXh2NWRpBr9fp8"
	accessSecret = "S0cOobEbehrkpMSCXSoqYzwRqmhCT7"
	bucketName = "filestoreserver"
)

func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	client, err := oss.New(ossEndpoint, accessKeyID, accessSecret)
	if err != nil {
		fmt.Printf("Initial oss client failed: %v\n", err)
		return nil
	}
	ossCli = client
	return ossCli
}

func Bucket() *oss.Bucket {
	cli := Client()
	bucket, err := cli.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Get oss bucket failed: %v\n", err)
		return nil
	}
	return bucket
}

func Download(objKey string) string {
	signURL, err := Bucket().SignURL(objKey, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Printf("Get the signed URL failed: %v\n", err)
		return ""
	}
	return signURL
}

// 指定为oss的某个bucket设定过期规则
func SetBucketLifecycle(bucketName string) {
	rule := oss.BuildLifecycleRuleByDays("rule1", "tmpfile/", true, 14)
	rules := []oss.LifecycleRule{rule}
	err := Client().SetBucketLifecycle(bucketName, rules)
	if err != nil {
		fmt.Printf("Set rule for bucket: %v failed, %v\n", bucketName, err)
	}
}
