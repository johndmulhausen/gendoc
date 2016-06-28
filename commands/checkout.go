package commands

import (
	"fmt"
	"os"

	allprojects "github.com/SvenDowideit/gendoc/allprojects"

	"github.com/codegangsta/cli"
)

var Checkout = cli.Command{
	Name:  "checkout",
	Usage: "checkout versions from "+allprojects.AllProjectsPath+" file",
	Flags: []cli.Flag{
	},
	Action: func(context *cli.Context) error {
        // TODO: checkout the version specified as the arg
		if context.NArg() != 1 {
            return fmt.Errorf("please specify a docs.docker.com branch to checkout")
        }
        publishSetBranch := context.Args()[0]
        err := checkout(allprojects.AllProjectsRepo, publishSetBranch)
        if err != nil {
            return err
        }

		setName, projects, err := allprojects.Load(allprojects.AllProjectsPath)
		if err != nil {
            if os.IsNotExist(err) {
                fmt.Printf("Please run `clone` command first.\n")
            }
			return err
		}
		fmt.Printf("publish-set: %s\n", setName)

        for _, p := range *projects {
            // TODO: don't ignore errors.
            fmt.Printf("-- %s\n", p.RepoName)
            checkout(p.RepoName, p.Ref)
        }
        return nil
	},
}

// TODO: re-write this to use --fetch - defaulting to true
func checkout(repoPath, ref string) error {
    //TODO what if its a tag
    err := allprojects.GitIn(repoPath, "checkout", ref)
    if err != nil {
        // do a fetch, in case it exists in remote
        err = allprojects.GitIn(repoPath, "fetch", "origin", ref+":remotes/origin/"+ref)
        if err != nil {
            // Last resourt, fetch all origin, and undo depth
            err = allprojects.GitIn(repoPath, "fetch", "--all")
            if err != nil {
                return err
            }
            err = allprojects.GitIn(repoPath, "fetch", "--tag")
            if err != nil {
                return err
            }
        }
        err = allprojects.GitIn(repoPath, "checkout", ref)
        if err != nil {
            err = allprojects.GitIn(repoPath, "checkout", "-b", ref, "remotes/origin/"+ref)
            if err != nil {
                return err
            }
        }
    }
    return err
}
