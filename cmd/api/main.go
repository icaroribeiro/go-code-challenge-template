package main

import (
	"fmt"
	"os"

	cmd "github.com/icaroribeiro/new-go-code-challenge-template/tools/api/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "api",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// @title New Go Code Challenge Template API
// @version 1.0
// @Description A REST API developed using Golang, Json Web Token and PostgreSQL database.
// @tag.name health check
// @tag.description It refers to the operation related to health check.
// @tag.name authentication
// @tag.description It refers to the operations related to authentication.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email icaroribeiro@hotmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /
// @schemes http
func main() {
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
