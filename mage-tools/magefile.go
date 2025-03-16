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

// update the dependency
func Update() error {
	return sh.Run("go", "mod", "download")
}

// build Creates the binary in the current directory.
func Build() error {
	mg.Deps(Clean)
	mg.Deps(BuidTCPClient)
	mg.Deps(BuidTCPServer)
	mg.Deps(BuidUDPClient)
	err := sh.Run("go", "build", "-o", "./bin/udp-server", "./cmd/udp-server/main.go")
	if err != nil {
		return err
	}
	return nil
}

// BuidTCPServer - 編譯 TCP server
func BuidTCPServer() error {
	err := sh.Run("go", "build", "-o", "./bin/tcp-server", "./cmd/tcp-server/main.go")
	if err != nil {
		return err
	}
	return nil
}

// BuildTCPClient - 編譯 TCP client
func BuidTCPClient() error {
	err := sh.Run("go", "build", "-o", "./bin/tcp-client", "./cmd/tcp-client/main.go")
	if err != nil {
		return err
	}
	return nil
}

// BuildUDPClient - 編譯 UDP client
func BuidUDPClient() error {
	err := sh.Run("go", "build", "-o", "./bin/udp-client", "./cmd/udp-client/main.go")
	if err != nil {
		return err
	}
	return nil
}

// LaunchTCPServer - start tcp server
func LaunchTCPServer() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/tcp-server")
	if err != nil {
		return err
	}
	return nil
}

// LaunchUDPServer - start udp server
func LaunchUDPServer() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/udp-server")
	if err != nil {
		return err
	}
	return nil
}

// LaunchTCPClient - start tcp client
func LaunchTCPClient() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/tcp-client")
	if err != nil {
		return err
	}
	return nil
}

// LaunchUDPClient - start udp client
func LaunchUDPClient() error {
	mg.Deps(Build)
	err := sh.RunV("./bin/udp-client")
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
