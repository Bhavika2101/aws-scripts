package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("ASIA4WFI752EY7NCYRBC", "UHPtc+zFrYd4RyrZHQJ5ZgBN6apUBdo1cXTxCKzP", "IQoJb3JpZ2luX2VjEB4aCmFwLXNvdXRoLTEiSDBGAiEAjm+4RNA0UCFu68BQTds2vQLNcKCiyVBisgr19NVMDb0CIQDHS3c0LQwrucswkVVcjLVTLc7dGZEa0e0wvTDjdPzRzyrrAQhXEAEaDDg3MjIzMjc3NTMwNSIME4jSjxETSUP0PUf/KsgB2AkrzxKM08voaCKQUG/7Ulz0p4ovR6wOb4ZO0BnH/hN9rPjeUT48qkFX3g6TPYWLRjj25N9WUphnt/wuueQgJFMMo+wYlBIQ/11chwnKG3loZG655qpoHlFP/72WqESc+ypDP5DpJMfIhEimzTYFOME8V9RaC9a6L82yLW3gZlBGISxEeqZ9euiZ+6goHwH1gdhyBdt0vM6tS2u1aTsBrk93kayOIZME+U9JT0kSnECBlNI4i3ZAx1/SvE83GxmWXdZpWpOm0C4w9rjojwY6lwEGAcsvkV44W3EfhSHJP2vvrHAKHd44bCbTTrBqmBRuG4x+cKDRIDRqlrzpSJf1gPmhX7KLQ+isFK2ch0A8voI/gBYRs6E1pczUv04ul05BrkKF7AjCECxCnYUSgAneZGpFj+RSkIKbJiCSdYtS2jXFyrKlZsq/q4eLZIeaubxjDY5IX94M0wYYL8RLOqKTUPm7+XxEjVC6"),
	})
	if err != nil {
		panic(err)
	}
	svc := ec2.New(sess)
	result, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String("testing5"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidKeyPair.Duplicate" {
			exitErrorf("Keypair %q already exists.", "testing5")
		}
		exitErrorf("Unable to create key pair: %s, %v.", "testing5", err)
	}

	fil, _ := os.Create("testing5.txt")
	data := []byte(*result.KeyMaterial)
	_, er := fil.Write(data)
	if er != nil {
		panic(er)
	}
	fil.Close()

	content, err := ioutil.ReadFile("testing5.txt")
	if err != nil {
		panic(err)
	}
	keypair := string(content)
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(100),
				},
			},
		},
		ImageId:      aws.String("ami-0c1a7eb9b30ab60cc"),
		InstanceType: aws.String("t2.micro"),
		KeyName:      aws.String(keypair),
		MaxCount:     aws.Int64(2),
		MinCount:     aws.Int64(1),
		SecurityGroupIds: []*string{
			aws.String("sg-01c8e7ad387140b9a"),
		},
		SubnetId: aws.String("subnet-e70d0b8f"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Purpose"),
						Value: aws.String("test"),
					},
				},
			},
		},
	}

	result1, err := svc.RunInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result1)

	vr := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String("i-0f94c91ae8974ff1f"),
		},
	}
	defer 
	ter, _ := svc.TerminateInstances(
		&ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				aws.String("i-0f94c91ae8974ff1f"),
			},
		},
	)


	fmt.Println(ter)
	res, _ := svc.DescribeInstances(vr)
	fmt.Println(res)

}
