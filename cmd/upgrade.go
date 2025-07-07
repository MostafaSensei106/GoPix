package cmd

import (
	"github.com/mostafasensei106/gopix/internal/upgrade"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade GoPix to the latest version",
	Long:  `Pull the latest changes from GitHub and update GoPix to the latest available version.`,
	Run: func(cmd *cobra.Command, args []string) {
	 upgrade.UpgradeGoPix()
	},
}
