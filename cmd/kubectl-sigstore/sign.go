//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sigstore/k8s-manifest-sigstore/pkg/k8smanifest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdSign() *cobra.Command {

	var imageRef string
	var inputDir string
	var keyPath string
	var output string
	var updateAnnotation bool
	var imageAnnotations []string
	cmd := &cobra.Command{
		Use:   "sign -f <YAMLFILE> [-i <IMAGE>]",
		Short: "A command to sign Kubernetes YAML manifests",
		RunE: func(cmd *cobra.Command, args []string) error {

			err := sign(inputDir, imageRef, keyPath, output, updateAnnotation, imageAnnotations)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&inputDir, "filename", "f", "", "file name which will be signed (if dir, all YAMLs inside it will be signed)")
	cmd.PersistentFlags().StringVarP(&imageRef, "image", "i", "", "signed image name which bundles yaml files")
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file name (if empty, use `<input>.signed`)")
	cmd.PersistentFlags().StringVarP(&keyPath, "key", "k", "", "path to your signing key (if empty, do key-less signing)")
	cmd.PersistentFlags().BoolVarP(&updateAnnotation, "annotation-metadata", "", true, "whether to update annotation and generate signed yaml file")
	cmd.PersistentFlags().StringArrayVarP(&imageAnnotations, "annotation", "a", []string{}, "extra key=value pairs to sign")

	return cmd
}

func sign(inputDir, imageRef, keyPath, output string, updateAnnotation bool, annotations []string) error {
	if output == "" && updateAnnotation {
		output = inputDir + ".signed"
	}

	anntns, err := parseAnnotations(annotations)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	so := &k8smanifest.SignOption{
		ImageRef:         imageRef,
		KeyPath:          keyPath,
		Output:           output,
		UpdateAnnotation: updateAnnotation,
		ImageAnnotations: anntns,
	}

	_, err = k8smanifest.Sign(inputDir, so)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	log.Info("signed manifest generated at ", output)
	return nil
}

func parseAnnotations(annotations []string) (map[string]interface{}, error) {
	annotationsMap := map[string]interface{}{}

	for _, annotation := range annotations {
		kvp := strings.SplitN(annotation, "=", 2)
		if len(kvp) != 2 {
			return nil, fmt.Errorf("invalid flag: %s, expected key=value", annotation)
		}

		annotationsMap[kvp[0]] = kvp[1]
	}
	return annotationsMap, nil
}
