package server

import (
	"github.com/naiba/cloudssh/cmd/client/dao"
	"github.com/spf13/cobra"
)

// DeleteCmd ..
var DeleteCmd *cobra.Command

func init() {
	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete server(s)",
	}
	DeleteCmd.Flags().UintSlice("id", []uint{}, "sever id list --id=\"1,3,4\"")
	DeleteCmd.Run = delete
}

func delete(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetUintSlice("id")
	dao.API.BatchDeleteServer(id, 0)
}
