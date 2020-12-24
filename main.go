package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Data struct {
	Version   string
	Statement []Statement
}
type Statement struct {
	Resource string
}

//func createPolicyVersionInput()

func main() {
	session := session.Must(session.NewSession())

	iamClient := iam.New(session)

	//https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.ListGroups
	//List all the groups
	lgOutput, err := iamClient.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, group := range lgOutput.Groups {
		//List the attached policies per each group
		lagpOutput, err := iamClient.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
			GroupName: group.GroupName,
		})
		if err != nil {
			log.Println(err)
		}
		// Loops through each attached policy
		for _, policy := range lagpOutput.AttachedPolicies {
			gpoutput, err := iamClient.GetPolicy(&iam.GetPolicyInput{
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				log.Println(err)
			}
			pv, _ := iamClient.GetPolicyVersion(&iam.GetPolicyVersionInput{
				VersionId: gpoutput.Policy.DefaultVersionId,
				PolicyArn: gpoutput.Policy.Arn,
			})

			decodedv, err := url.QueryUnescape(aws.StringValue(pv.PolicyVersion.Document))
			if err != nil {
				log.Fatal(err)
			}

			aData := Data{}
			err = json.Unmarshal([]byte(decodedv), &aData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(*policy.PolicyName)
			fmt.Println(aData.Statement)
		}
	}
}
