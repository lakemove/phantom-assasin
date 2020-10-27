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
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
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
	r.GET("/payment-ack", handlePaymentAck)
	r.POST("/payments", handlePayments)
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

func handlePaymentAck(c *gin.Context) {
	c.String(200, "hello gin\n")
}

type LibraTransfer struct {
	SourceSystem   string
	SendReceive    string
	Status         string
	TransferId     string
	TransferStatus string
	AckStatus      string
	Currency       string `xml:"currency"`
	Amount         string `xml:"amount"`
}

func handlePayments(c *gin.Context) {
	var xml LibraTransfer
	if err := c.ShouldBindXML(&xml); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"hello": xml})
}

func createTransfer(t LibraTransfer) (id string, err error) {
	id, err = "123", nil
	return
}

func getTransferStatus(id string) (status string, err error) {
	return "SILO_OP_PENDING_APPROVAL", nil
}

func approveTransfer(id string) (err error) {
	return nil
}

func updateRow(t LibraTransfer) (err error) {
	return nil
}

func sendAck(t LibraTransfer) (err error) {
	return nil
}

func onTransfer(id string) (err error) {
	var row LibraTransfer
	if row.SendReceive == "Send" && row.SourceSystem == "MX_FXCASH" { //transfer out
		if row.Status == "RECEIVED" {
			transferId, err := createTransfer(row)
			if err != nil {
				log.Println("fail to create transfer", id, err)
			}
			row.TransferId = transferId
			row.Status = "DELIVERED"
			err = updateRow(row)
			if err != nil {
				log.Panicln("fail to update db", err)
			}
		}
		if row.Status == "DELIVERED" {
			for {
				if row.TransferStatus == "SILO_OP_COMPLETED" {
					if row.AckStatus != "Settled" {
						err = sendAck(row)
						if err != nil {
							log.Println("fail to ack")
							time.Sleep(1 * time.Minute)
							continue
						}
					}
					break //break for loop, reached final state
				}
				if row.TransferStatus == "SILO_OP_PENDING_APPROVAL" {
					if row.AckStatus == "Received" {
						err = sendAck(row)
						if err != nil {
							log.Println("fail to ack", err)
							time.Sleep(1 * time.Minute)
							continue
						}
					}
					err = approveTransfer(row.TransferId)
					if err != nil {
						log.Println("fail to approve transfer ", id, err)
						time.Sleep(1 * time.Minute)
						continue
					}
				}
				currentTransferStatus, err := getTransferStatus(row.TransferId)
				if err != nil {
					log.Println("fail to pull Transfer status, will wait 10 seconds before next poll", err)
					time.Sleep(10 * time.Second)
					continue
				}
				if row.TransferStatus != currentTransferStatus {
					log.Printf("Transfer status changed from %s to %s", row.TransferStatus, currentTransferStatus)
					row.TransferStatus = currentTransferStatus
					err = updateRow(row)
					if err != nil {
						log.Println("fail to update db", err)
						time.Sleep(1 * time.Minute)
						continue
					}
				}
				time.Sleep(1 * time.Minute)
			}
		}
		return
	}
	if row.SendReceive == "Receive" && row.SourceSystem == "MX_FXCASH" { //transfer in
		if row.Status == "RECEIVED" {
			//todo find a match and update status to MATCHED
		}
		if row.Status == "MATCHED" {
			err = sendAck(row)
			if err != nil {
				//todo
			}
		}
		return
	}
	if row.SourceSystem == "ZODIA" { //transfer in
		return
	}
	log.Printf("unable to process %s\n", id)
	return fmt.Errorf(" unknown type %s %s", row.SourceSystem, row.SendReceive)
}

func main() {
	if err := buildCmd().Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
