package provision

import (
	"fmt"
	"vmctl/src/config"

	"github.com/urfave/cli/v2"
)

// The struct with all possible bootstrap options
type ProvisionOptions struct {
	SkipAll            bool   // Autoapprove all prompts
	ThinPool           string // Name of thinpool (default: flintlock[-dev])
	Disk               string // Disk to use for DirectLVM thinpool (ignored for dev)
	gRPCEndpoint       string // Address for gRPC server (default: 0.0.0.0:9090)
	ParentInterface    string // Interface of the default route of the host
	Development        bool   // Set up development environment (uses Loop thinpools)
	FirecrackerVersion string // Firecracker version to use
	FlintlockVersion   string // Flintlock Version to use
	ContainerdVersion  string // ContainerD Version to use
}

// Setup all flags for the bootstrap commands
func ProvisionFlags() []cli.Flag {
	var bootFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "skip-all",
			Value: true,
			Usage: "Non interactive bootstrap",
		},
		&cli.BoolFlag{
			Name:    "development",
			Aliases: []string{"d", "dev"},
			Value:   true,
			Usage:   "Set up development environment (uses Loop thinpools)",
		},
		&cli.StringFlag{
			Name:  "thinpool",
			Value: "flintlock",
			Usage: "Name of thinpool (default: flintlock[-dev])",
		},
		&cli.StringFlag{
			Name:  "disk",
			Usage: "Disk to use for DirectLVM thinpool (ignored for dev)",
		},
		&cli.StringFlag{
			Name:    "grpc-endpoint",
			Aliases: []string{"b", "grpc"},
			Value:   "0.0.0.0:9090",
			Usage:   "Disk to use for DirectLVM thinpool (ignored for dev)",
		},
		&cli.StringFlag{
			Name:  "fc-version",
			Value: "latest",
			Usage: "Firecracker Version",
		},
		&cli.StringFlag{
			Name:  "fl-version",
			Value: "latest",
			Usage: "Flintlock Version",
		},
		&cli.StringFlag{
			Name:  "cd-version",
			Value: "latest",
			Usage: "Containerd Version",
		},
	}
	return bootFlags
}

// Takes the CLI Context and builds the Options object
// This then can be used by every `init` function for the
// different methods that will be ran on the bootstrap
func buildProvisionOptions(cCtx *cli.Context) ProvisionOptions {
	var skipAll = cCtx.Bool("skip-all")
	var develop = cCtx.Bool("development")
	var thinpool = cCtx.String("thinpool")
	var disk = cCtx.String("disk")
	var gRPCEndpoint = cCtx.String("grpc-endpoint")
	var FirecrackerVersion = cCtx.String("fc-version")
	var FlintlockVersion = cCtx.String("fl-version")
	var ContainerdVersion = cCtx.String("cd-version")

	var newOptions = ProvisionOptions{
		SkipAll:            skipAll,
		Development:        develop,
		ThinPool:           thinpool,
		Disk:               disk,
		gRPCEndpoint:       gRPCEndpoint,
		FirecrackerVersion: FirecrackerVersion,
		FlintlockVersion:   FlintlockVersion,
		ContainerdVersion:  ContainerdVersion,
	}

	return newOptions
}

// Function that actually bootstraps a Flintlock-enabled server
// It will run all checks & steps necessary to get a Flintlock server running
// This includes Containerd, Firecracker, etc.
func ProvisionAction(cCtx *cli.Context) error {
	var bootConfig = buildProvisionOptions(cCtx)

	fmt.Println("Provisioning host...")

	// Checking dependencies
	fmt.Println("Checking dependencies...")
	initDependencies(bootConfig, &config.Cfg)

	// Setting up Networking
	fmt.Println("Setting up network...")
	initNetwork(bootConfig, &config.Cfg)

	// Setting Up Containerd
	fmt.Println("Creating Containerd directories...")
	initContainerd(bootConfig, &config.Cfg)

	// // Setup Thin Pool
	// fmt.Println("Setup Disks...")
	// initDisks(bootConfig, &config.Cfg)

	// // Setup Firecracker
	// // install_firecracker "$fc_version"
	// fmt.Println("Installing Firecracker...")
	// initFirecracker(bootConfig, &config.Cfg)

	// // Installing Flintlock
	// fmt.Println("Installing Flintlock...")
	// initFlintlock(bootConfig, &config.Cfg)

	return nil
}
