package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/izquiratops/bundler-tools/commands"
)

type CommonFlags struct {
	Entry  string
	Dist   string
	Minify bool
}

func defineCommonFlags(fs *flag.FlagSet) *CommonFlags {
	cf := &CommonFlags{}

	fs.StringVar(&cf.Entry, "entry", "src/main.js", "Path to the entry point")
	fs.StringVar(&cf.Dist, "dist", "dist/", "Path to the distribution directory")
	fs.BoolVar(&cf.Minify, "minify", false, "Enable JavaScript minification")

	return cf
}

func printHelp() {
	fmt.Println("Usage: github.com/izquiratops/bundler-tools <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  build\tBuild the project")
	fmt.Println("  serve\tServe the project")
	fmt.Println("  clean\tClean the distribution directory")
	fmt.Println("  help\tShow this help message")
	fmt.Println("\nOptions:")
	fmt.Println("  -d, --dist string\tPath to the distribution directory (default \"dist/\")")
	fmt.Println("  -e, --entry string\tPath to the entry point (default \"src/main.js\")")
	fmt.Println("  -m, --minify string\tEnable JavaScript minification (default \"false\")")
}

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildFlags := defineCommonFlags(buildCmd)

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveFlags := defineCommonFlags(serveCmd)

	cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	cleanFlags := defineCommonFlags(cleanCmd)

	switch os.Args[1] {
	case "build":
		buildCmd.Parse(os.Args[2:])
		commands.Build(buildFlags.Entry, buildFlags.Dist, buildFlags.Minify)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		commands.Serve(serveFlags.Entry, serveFlags.Dist, serveFlags.Minify)
	case "clean":
		cleanCmd.Parse(os.Args[2:])
		commands.Clean(cleanFlags.Dist)
	case "help":
		printHelp()
	default:
		log.Fatal("expected an action subcommand (build, serve or clean)")
	}
}
