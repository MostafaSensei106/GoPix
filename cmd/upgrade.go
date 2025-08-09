package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MostafaSensei106/GoPix/internal/upgrade"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade GoPix to the latest version",
	Long:  `Pull the latest changes from GitHub and update GoPix to the latest available version.`,
	Run: func(cmd *cobra.Command, args []string) {
	 upgrade.UpgradeGoPix()
	},
}
