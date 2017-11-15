package commands

import (
	"fmt"

	"github.com/phase2/rig/util"
	"github.com/urfave/cli"
)

// Remove is the command for deleting a Docker Machine
type Remove struct {
	BaseCommand
}

// Commands returns the operations supported by this command
func (cmd *Remove) Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "remove",
			Usage: "Remove the docker-machine",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force removal. Don't prompt to backup /data",
				},
			},
			Before: cmd.Before,
			Action: cmd.Run,
		},
	}
}

// Run executes the `rig remove` command
func (cmd *Remove) Run(c *cli.Context) error {
	if util.SupportsNativeDocker() {
		return cmd.Success("Remove is not needed on Linux")
	}

	if !cmd.machine.Exists() {
		return cmd.Error(fmt.Sprintf("No machine named '%s' exists.", cmd.machine.Name), "MACHINE-NOT-FOUND", 12)
	}

	cmd.out.Status(fmt.Sprintf("Removing '%s'", cmd.machine.Name))
	force := c.Bool("force")
	if !force {
		cmd.out.Warn("!!!!! This operation is destructive. You may lose important data. !!!!!!!")
		cmd.out.Warn("Run 'rig data-backup' if you want to save your /data volume.")

		if !util.AskYesNo("Are you sure you want to remove '" + cmd.machine.Name + "'") {
			return cmd.Success("Remove was aborted")
		}
	}

	// Run kill first.
	kill := Kill{cmd.BaseCommand}
	if err := kill.Run(c); err != nil {
		return err
	}

	cmd.out.Spin("Removing the docker Virtual Machine")
	if err := cmd.machine.Remove(); err != nil {
		cmd.out.Oops("Failed to remove the docker Virtual Machine")
		return cmd.Error(err.Error(), "MACHINE-REMOVE-FAILED", 12)
	}

	cmd.out.Success("Failed to remove the docker Virtual Machine")
	return cmd.Success(fmt.Sprintf("Machine '%s' removed", cmd.machine.Name))
}
