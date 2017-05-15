package lib

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/viper"
)

var cfDir = ".cloudformation"

//it's purpose is to find a cloudformation template from the template dir, create it with the aws api, and return
//the output associated with it (e.g. aws)
type CloudFormationTemplate struct {
	Name      string
	StackName string
	Output    map[string]string
	Bytes     []byte
}

func NewCloudFormationTemplate(name string) *CloudFormationTemplate {
	c := new(CloudFormationTemplate)
	c.ReadFile(name)
	c.Name = name
	return c
}

func (c *CloudFormationTemplate) CreateStack(suffix string) {
	c.LoadCloudFormationTemplate(suffix)
}

func (c *CloudFormationTemplate) LoadCloudFormationTemplate(suffix string) {
	fmt.Printf("\nRunning CloudFormation stack [%s]", c.Name)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	svc := cloudformation.New(sess)

	templateBody := string(c.Bytes)
	c.StackName = c.Name + "-" + suffix
	cloudFormationValues := make(map[string]string)

	createStackParams := &cloudformation.CreateStackInput{
		StackName:    aws.String(c.StackName),
		TemplateBody: aws.String(templateBody),
	}

	describeStacksParams := &cloudformation.DescribeStacksInput{
		StackName: aws.String(c.StackName),
	}

	out, err := svc.CreateStack(createStackParams)
	if err != nil {
		if strings.Contains(err.Error(), "AlreadyExistsException") {
			fmt.Printf("\nCloudFormation stack [%s] already exists. Skipping...", c.Name)
			descOut, err := svc.DescribeStacks(describeStacksParams)
			if err != nil {
				panic(err)
			}
			c.ParseOutput(descOut, cloudFormationValues)
		}
		fmt.Printf("%s", err)
		panic(err)
	} else {
		fmt.Printf("%s", out)
	}

	stackReady := false

	for stackReady != true {

		descOut, err := svc.DescribeStacks(describeStacksParams)
		if err != nil {
			fmt.Printf("%s", err)
			panic(err)
		} else {
			fmt.Printf("\nCloudFormation stack [%s] is creating...", c.Name)
		}

		if *descOut.Stacks[0].StackStatus == "CREATE_COMPLETE" {
			stackReady = true
			fmt.Printf("\nCloudFormation stack [%s] ready...\n", c.Name)
			c.ParseOutput(descOut, cloudFormationValues)
		}

		time.Sleep(time.Second * 7)
	}

}

func (c *CloudFormationTemplate) ParseOutput(descOut *cloudformation.DescribeStacksOutput, cloudFormationValues map[string]string) {
	stack := descOut.Stacks[0]
	for _, cfOutput := range stack.Outputs {
		trimKey := strings.TrimSpace(*cfOutput.OutputKey)
		trimVal := strings.TrimSpace(*cfOutput.OutputValue)
		cloudFormationValues[trimKey] = trimVal
	}
	c.Output = cloudFormationValues
}

func (c *CloudFormationTemplate) ReadFile(name string) {
	path := viper.GetString("release")
	cfTemplate := path + "/" + cfDir + "/" + name + ".yaml"
	yamlBytes, err := ioutil.ReadFile(cfTemplate)
	if err != nil {
		panic(err)
	}
	c.Bytes = yamlBytes

}
