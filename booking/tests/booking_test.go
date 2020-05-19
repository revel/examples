package tests_test

// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"testing"

	"github.com/revel/cmd/model"
	revelParser "github.com/revel/cmd/parser"
	"github.com/revel/revel"
	"strings"
)

func getRevelContainer() *model.RevelContainer{

	paths, _ := model.NewRevelPaths("prod","github.com/revel/examples/booking", model.NewWrappedRevelCallback(nil, nil))

	return paths
}

// A test for processing the booking application source
func TestProcessBookingSource(t *testing.T) {
	revel.Init("prod", "github.com/revel/examples/booking", "")
	sourceInfo, err := revelParser.ProcessSource(getRevelContainer())
	if err != nil {
		t.Fatal("Failed to process booking source with error:", err)
	}

	controllerPackage := "github.com/revel/examples/booking/app/controllers"
	expectedControllerSpecs := []*model.TypeInfo{
		{"Application", controllerPackage, "controllers", nil, nil},
		{"Hotels", controllerPackage, "controllers", nil, nil},
	}
	specList := []*model.TypeInfo{}
	for _,x := range sourceInfo.ControllerSpecs() {
		if strings.HasPrefix(x.ImportPath,controllerPackage) {
			specList = append(specList,x)
		}
	}
	if len(specList) != len(expectedControllerSpecs) {
		t.Errorf("Unexpected number of controllers found.  Expected %d, Found %d",
			len(expectedControllerSpecs), len(sourceInfo.ControllerSpecs()))
	}

NEXT_TEST:
	for _, expected := range expectedControllerSpecs {
		for _, actual := range sourceInfo.ControllerSpecs() {
			if actual.StructName == expected.StructName {
				if actual.ImportPath != expected.ImportPath {
					t.Errorf("%s expected to have import path %s, actual %s",
						actual.StructName, expected.ImportPath, actual.ImportPath)
				}
				if actual.PackageName != expected.PackageName {
					t.Errorf("%s expected to have package name %s, actual %s",
						actual.StructName, expected.PackageName, actual.PackageName)
				}
				continue NEXT_TEST
			}
		}
		t.Errorf("Expected to find controller %s, but did not.  Actuals: %s",
			expected.StructName, sourceInfo.ControllerSpecs())
	}
}

// Performance test for booking application
// this tests the speed of the command line utility to process the source of the booking application
func BenchmarkProcessBookingSource(b *testing.B) {
	revel.Init("", "github.com/revel/examples/booking", "")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := revelParser.ProcessSource(getRevelContainer())
		if err != nil {
			b.Error("Unexpected error:", err)
		}
	}
}
