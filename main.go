package main

import (
	"context"
	"os"
	"time"

	"github.com/goneup/go-pingen-sdk/pingen"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func main() {
	log.Info("Startup")
	err := godotenv.Load(".env.prod")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	org := os.Getenv("PINGEN_ORG")

	var c *pingen.Client
	c, err = pingen.NewClient(clientID, clientSecret, true, org, context.Background())
	if err != nil {
		log.Fatalf("Init failed, error %w", err)
	}
	c.ListLetters()

	//demo run
	bytes, err := os.ReadFile("demo_mustermann.pdf")
	if err != nil {
		log.Fatal(err)
	}

	//test send
	log.Info("Testing sending a letter")
	createData := &pingen.CreateData{}
	createData.Data.Attributes.FileOriginalName = "Testupload.pdf"
	createData.Data.Attributes.AddressPosition = "left"
	createData.Data.Attributes.AutoSend = false

	result, _ := c.CreateLetter(bytes, createData)

	time.Sleep(5 * time.Second) //wait until api catches up
	//c.ListLetters()
	c.GetLetter(result.Data.ID)

	sendData := &pingen.SendData{}
	sendData.Data.Attributes.DeliveryProduct = "cheap"
	sendData.Data.Attributes.PrintMode = "simplex"
	sendData.Data.Attributes.PrintSpectrum = "color"
	c.SendLetter(result.Data.ID, sendData)
	c.GetLetter(result.Data.ID)

	//test delete
	log.Info("Testing create & delete")
	createData = &pingen.CreateData{}
	createData.Data.Attributes.FileOriginalName = "DeleteTest.pdf"
	createData.Data.Attributes.AddressPosition = "left"
	createData.Data.Attributes.AutoSend = false

	result, _ = c.CreateLetter(bytes, createData)

	time.Sleep(5 * time.Second) //wait until api catches up
	c.GetLetter(result.Data.ID)
	c.DeleteLetter(result.Data.ID)
}
