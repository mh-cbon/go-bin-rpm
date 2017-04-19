// Create binary rpm package with ease
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mh-cbon/go-bin-rpm/rpm"
	"github.com/mh-cbon/verbose"
	"github.com/urfave/cli"
)

// VERSION is the build version number.
var VERSION = "0.0.0"
var logger = verbose.Auto()

func main() {
	app := cli.NewApp()
	app.Name = "go-bin-rpm"
	app.Version = VERSION
	app.Usage = "Generate a binary rpm package"
	app.UsageText = "go-bin-rpm <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "generate-spec",
			Usage:  "Generate the SPEC file",
			Action: generateSpec,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
				cli.StringFlag{
					Name:  "a, arch",
					Value: "",
					Usage: "Target architecture of the build",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Target version of the build",
				},
			},
		},
		{
			Name:   "generate",
			Usage:  "Generate the package",
			Action: generatePkg,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
				cli.StringFlag{
					Name:  "b, build-area",
					Value: "pkg-build",
					Usage: "Path to the build area",
				},
				cli.StringFlag{
					Name:  "a, arch",
					Value: "",
					Usage: "Target architecture of the build",
				},
				cli.StringFlag{
					Name:  "o, output",
					Value: "",
					Usage: "Output package to this path",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Target version of the build",
				},
			},
		},
		{
			Name:   "test",
			Usage:  "Test the package json file",
			Action: testPkg,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
			},
		},
	}

	app.Run(os.Args)
}

func generateSpec(c *cli.Context) error {
	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")

	rpmJSON := rpm.Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := rpmJSON.Normalize(arch, version); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if spec, err := rpmJSON.GenerateSpecFile(""); err != nil {
		return cli.NewExitError(err.Error(), 1)
	} else {
		fmt.Printf("%s", spec)
	}

	return nil
}

func generatePkg(c *cli.Context) error {
	var err error

	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")
	buildArea := c.String("build-area")
	output := c.String("output")

	rpmJSON := rpm.Package{}

	if err3 := rpmJSON.Load(file); err3 != nil {
		return cli.NewExitError(err3.Error(), 1)
	}

	if buildArea, err = filepath.Abs(buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err2 := rpmJSON.Normalize(arch, version); err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	rpmJSON.InitializeBuildArea(buildArea)

	if err = rpmJSON.WriteSpecFile("", buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = rpmJSON.RunBuild(buildArea, output); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("\n\nAll done!")

	return nil
}

func testPkg(c *cli.Context) error {
	file := c.String("file")

	rpmJSON := rpm.Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("File is correct")

	return nil
}
