//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// clean the build binary
func Clean() error {
	return sh.Rm("bin")
}

// Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	err := sh.Run("go", "build", "-o", "./bin/basic-sample", "./cmd/basic-sample/main.go")
	if err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "./bin/hash-sample", "./cmd/hash-sample/main.go")
}

// start the basic-sample
func LaunchBasicSample() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/basic-sample")
	if err != nil {
		return err
	}
	return nil
}

// start the hash-sample
func LaunchHashSample() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/hash-sample")
	if err != nil {
		return err
	}
	return nil
}

// run the test
func Test() error {
	err := sh.RunV("go", "test", "-v", "./...")
	if err != nil {
		return err
	}
	return nil
}
