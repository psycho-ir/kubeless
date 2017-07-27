/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/kubeless/kubeless/pkg/spec"
	"github.com/kubeless/kubeless/pkg/utils"
	"github.com/spf13/cobra"
)

var ingressCreateCmd = &cobra.Command{
	Use:   "create <name> FLAG",
	Short: "create a route to function",
	Long:  `create a route to function`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("Need exactly one argument - ingress name")
		}
		ingressName := args[0]
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logrus.Fatal(err)
		}
		domain, err := cmd.Flags().GetString("domain")
		if err != nil {
			logrus.Fatal(err)
		}

		function, err := cmd.Flags().GetString("function")
		if err != nil {
			logrus.Fatal(err)
		}

		err = validateInput(function, ns)
		if err != nil {
			logrus.Fatal(err)
		}

		client := utils.GetClientOutOfCluster()

		err = utils.CreateIngress(client, ingressName, function, domain, ns)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func validateInput(function, ns string) error {
	tprClient, err := utils.GetTPRClientOutOfCluster()
	if err != nil {
		return err
	}
	funcList := spec.FunctionList{}
	err = tprClient.Get().
		Resource("functions").
		Namespace(ns).
		Do().
		Into(&funcList)
	if err != nil {
		return err
	}

	for _, f := range funcList.Items {
		if f.Metadata.Name == function {
			return nil
		}
	}

	return errors.New("function doesn't exist")
}

func init() {
	ingressCreateCmd.Flags().StringP("domain", "", "", "Specify runtime")
	ingressCreateCmd.Flags().StringP("function", "", "", "Specify handler")
}
