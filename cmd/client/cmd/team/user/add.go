package user

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/liamylian/x-rsa/golang/xrsa"
	"github.com/naiba/cloudssh/cmd/client/dao"
	"github.com/naiba/cloudssh/internal/apiio"
	"github.com/naiba/cloudssh/pkg/xcrypto"
	"github.com/spf13/cobra"
)

// AddCmd ..
var AddCmd *cobra.Command

func init() {
	AddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user to team",
	}
	AddCmd.Run = add
}

func add(cmd *cobra.Command, args []string) {
	teamID, _ := cmd.Parent().Parent().PersistentFlags().GetUint64("oid")
	if teamID == 0 {
		log.Println("must set team ID")
		return
	}
	var req apiio.AddTeamUserRequest

	fmt.Print("User Email: ")
	fmt.Scanf("%s", &req.Email)

	fmt.Print("Permission (1 readOnly 2 readWrite 3 owner): ")
	fmt.Scanf("%d", &req.Permission)

	u, err := dao.API.GetUser(req.Email)
	if err != nil {
		log.Println("GetUser", err)
		return
	}
	cs, err := xcrypto.NewCipherString(dao.Conf.User.Privatekey)
	if err != nil {
		log.Println("NewCipherString", err)
		return
	}
	userPrivateKeyByte, err := cs.Decrypt(dao.Conf.MasterKey)
	if err != nil {
		log.Println("Decrypt user privateKey", err)
		return
	}
	userXr, err := xrsa.NewXRsa([]byte(u.Data.Pubkey), userPrivateKeyByte)
	if err != nil {
		log.Println("NewXRsa", err)
		return
	}
	info, err := dao.API.GetTeamByID(teamID)
	if err != nil {
		log.Println("GetTeamByID", err)
		return
	}
	xr, err := dao.Conf.GerUserXRsa()
	if err != nil {
		log.Println("GerUserXRsa", err)
		return
	}
	teamPkStr, err := xr.PrivateDecrypt(info.TeamUser.PrivateKey)
	if err != nil {
		log.Println("PrivateDecrypt", err)
		return
	}
	req.Prikey, err = userXr.PublicEncrypt(teamPkStr)
	if err != nil {
		log.Println("PublicEncrypt", err)
		return
	}

	body, err := dao.API.Do(fmt.Sprintf("/team/%d/user", teamID), "POST", req)
	if err != nil {
		log.Println("API Request", err)
		return
	}
	var resp apiio.Response
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Println("API Request", string(body), err)
		return
	}
	if !resp.Success {
		log.Println("API Request", resp.Message)
		return
	}
	log.Println(resp.Message)
}
