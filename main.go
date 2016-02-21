package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone AWS CodeDeploy Plugin built from %s\n", buildCommit)

	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if vargs.Application == "" {
		vargs.Application = repo.Name
	}

	if vargs.RevisionType == "" {
		vargs.RevisionType = codedeploy.RevisionLocationTypeGitHub
	}

	if vargs.AccessKey == "" {
		fmt.Println("Please provide an access key id")
		os.Exit(1)
	}

	if vargs.SecretKey == "" {
		fmt.Println("Please provide a secret access key")
		os.Exit(1)
	}

	if vargs.Region == "" {
		fmt.Println("Please provide a region")
		os.Exit(1)
	}

	if vargs.DeploymentGroup == "" {
		fmt.Println("Please provide a deployment group")
		os.Exit(1)
	}

	var location *codedeploy.RevisionLocation

	switch vargs.RevisionType {
	case codedeploy.RevisionLocationTypeGitHub:
		location = &codedeploy.RevisionLocation{
			RevisionType: aws.String(vargs.RevisionType),
			GitHubLocation: &codedeploy.GitHubLocation{
				CommitId:   aws.String(build.Commit),
				Repository: aws.String(repo.FullName),
			},
		}
	case codedeploy.RevisionLocationTypeS3:
		if vargs.BundleType == "" {
			fmt.Println("Please provide a bundle type")
			os.Exit(1)
		}

		if vargs.BucketName == "" {
			fmt.Println("Please provide a bucket name")
			os.Exit(1)
		}

		if vargs.BucketKey == "" {
			fmt.Println("Please provide a bucket key")
			os.Exit(1)
		}

		switch vargs.BundleType {
		case codedeploy.BundleTypeTar:
		case codedeploy.BundleTypeTgz:
		case codedeploy.BundleTypeZip:
		default:
			fmt.Println("Invalid bundle type")
			os.Exit(1)
		}

		location = &codedeploy.RevisionLocation{
			RevisionType: aws.String(vargs.RevisionType),
			S3Location: &codedeploy.S3Location{
				BundleType: aws.String(vargs.BundleType),
				Bucket:     aws.String(vargs.BucketName),
				Key:        aws.String(vargs.BucketKey),
				ETag:       aws.String(vargs.BucketEtag),
				Version:    aws.String(vargs.BucketVersion),
			},
		}
	default:
		fmt.Println("Invalid revision type")
		os.Exit(1)
	}

	svc := codedeploy.New(
		session.New(&aws.Config{
			Region: aws.String(vargs.Region),
			Credentials: credentials.NewStaticCredentials(
				vargs.AccessKey,
				vargs.SecretKey,
				"",
			),
		}),
	)

	_, err := svc.CreateDeployment(
		&codedeploy.CreateDeploymentInput{
			ApplicationName:               aws.String(vargs.Application),
			DeploymentConfigName:          aws.String(vargs.DeploymentConfig),
			DeploymentGroupName:           aws.String(vargs.DeploymentGroup),
			Description:                   aws.String(vargs.Description),
			IgnoreApplicationStopFailures: aws.Bool(vargs.IgnoreStopFailures),
			Revision:                      location,
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully deployed")
}
