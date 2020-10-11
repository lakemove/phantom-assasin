package main

/*
[x] cli subcommand
[ ] sql oracle
[ ] Memorized
[ ] sign/verify EC/RS keys
[ ] xml bind, extract
[ ] json parse/serialize
[ ] http client disable keep-alive
[ ] http client ssl
[ ] ArrayBlockingQueue
*/

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

func httpClient() *http.Client {
	client := &http.Client{}
	return client
}

func buildCmd() *cobra.Command {
	var httpCmd = &cobra.Command{
		Use:   "http",
		Short: "run http server",
		Run:   runHttp,
	}
	httpCmd.Flags().String("addr", ":8000", "bind address, [host]:<port>")
	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "create db tables",
		Run:   runDB,
	}
	var rootCmd = &cobra.Command{
		Use: "leo",
		Run: runRoot,
	}
	rootCmd.PersistentFlags()
	rootCmd.PersistentFlags().StringToString("ds", nil, "datasource options")
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(dbCmd)
	return rootCmd
}

func runHttp(cmd *cobra.Command, args []string) {
	addr, err := cmd.Flags().GetString("addr")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/", handle1)
	r.POST("/xml", handleXml)
	r.Run(addr)
}

func runRoot(cmd *cobra.Command, args []string) {
	log.Println(db)
}

func runDB(cmd *cobra.Command, args []string) {
	ds, err := cmd.Flags().GetStringToString("ds")
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open(ds["driver"], ds["name"])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success", db)
}

func servHttp(addr string) {
}

func handle1(c *gin.Context) {
	c.String(200, "hello gin\n")
}

type LibraTransfer struct {
	Currency string `xml:"currency"`
	Amount   string `xml:"amount"`
}

func handleXml(c *gin.Context) {
	var xml LibraTransfer
	if err := c.ShouldBindXML(&xml); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"hello": xml})
}

func main() {
	if err := buildCmd().Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
