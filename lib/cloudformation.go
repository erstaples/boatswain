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

type CloudFormationTemplate struct {
	Name      string
	StackName string
	Output    map[string]string
	Bytes     []byte
}

//NewCloudFormationTemplate takes a cloudformation template filepath and
//reads the file contents, loads it into .Bytes field
func NewCloudFormationTemplate(name string) *CloudFormationTemplate {
	c := CloudFormationTemplate{Name: name}
	c.ReadFile(name)
	return &c
}

//CreateStack calls the cloudformation.CreateStack method. If the stack has already been created,
//then it gets the output from DescribeStack. Also sets .StackName property
func (c *CloudFormationTemplate) CreateStack(suffix string) {
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
			return
		} else {
			fmt.Printf("%s", err)
			panic(err)
		}
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

//ParseOutput takes a cloudformation.DescribeStacksOutput, iterates over the output and sets
//the .Output map field on the instance
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
