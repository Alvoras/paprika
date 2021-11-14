package main

import (
	"fmt"
	"github.com/alvoras/paprika/internal/cli"
	"os"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//res, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	//if err != nil {
	//	panic(err)
	//}

	//for _, resas := range res.Reservations{
	//	for _, inst := range resas.Instances{
	//
	//		if inst.State.Name == "running"{
	//			fmt.Println(*inst.InstanceId)
	//		}
	//	}
	//}
}
