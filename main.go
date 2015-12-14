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
	build     string
	buildDate string
)

func main() {
	fmt.Printf("Drone AWS CodeDeply Plugin built at %s\n", buildDate)

	workspace := drone.Workspace{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if len(vargs.AccessKeyID) == 0 {
		fmt.Println("Please provide an access key")

		os.Exit(1)
		return
	}

	if len(vargs.SecretAccessKey) == 0 {
		fmt.Println("Please provide a secret key")

		os.Exit(1)
		return
	}

	if len(vargs.Region) == 0 {
		fmt.Println("Please provide a region")

		os.Exit(1)
		return
	}

	svc := codedeploy.New(
		session.New(&aws.Config{
			Region:      aws.String(vargs.Region),
			Credentials: credentials.NewStaticCredentials(vargs.AccessKeyID, vargs.SecretAccessKey, ""),
		}))

	resp, err := svc.CreateDeployment(
		&codedeploy.CreateDeploymentInput{
			ApplicationName:               aws.String("ApplicationName"),
			DeploymentConfigName:          aws.String("DeploymentConfigName"),
			DeploymentGroupName:           aws.String("DeploymentGroupName"),
			Description:                   aws.String("Description"),
			IgnoreApplicationStopFailures: aws.Bool(true),
			Revision: &codedeploy.RevisionLocation{
				GitHubLocation: &codedeploy.GitHubLocation{
					CommitId:   aws.String("CommitId"),
					Repository: aws.String("Repository"),
				},
				RevisionType: aws.String("RevisionLocationType"),
				S3Location: &codedeploy.S3Location{
					Bucket:     aws.String("S3Bucket"),
					BundleType: aws.String("BundleType"),
					ETag:       aws.String("ETag"),
					Key:        aws.String("S3Key"),
					Version:    aws.String("VersionId"),
				},
			},
		})

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
		return
	}

	fmt.Println(resp)
}
