package server

import (
	"github.com/naiba/cloudssh/cmd/client/dao"
	"github.com/spf13/cobra"
)

// CreateCmd ..
var CreateCmd *cobra.Command

func init() {
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create server",
	}
	CreateCmd.Run = create
}

func create(cmd *cobra.Command, args []string) {
	dao.API.CreateServer(0)
}
