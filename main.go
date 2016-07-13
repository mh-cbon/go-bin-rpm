package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mh-cbon/go-bin-rpm/rpm"
	"github.com/mh-cbon/verbose"
	"github.com/urfave/cli"
)

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
	}

	app.Run(os.Args)
}

func generateSpec(c *cli.Context) error {
	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")

	rpmJson := rpm.Package{}

	if err := rpmJson.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := rpmJson.Normalize(arch, version); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if spec, err := rpmJson.GenerateSpecFile(""); err != nil {
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

	rpmJson := rpm.Package{}

	if err := rpmJson.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

  if buildArea, err = filepath.Abs(buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
  }

	if err := rpmJson.Normalize(arch, version); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

  rpmJson.InitializeBuildArea(buildArea)

	if err = rpmJson.WriteSpecFile("", buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
  }

	if err = rpmJson.RunBuild(buildArea, output); err != nil {
		return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("\n\nAll done!")

	return nil
}
