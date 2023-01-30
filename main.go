package main

import (
	"context"
	"go-pingen-sdk/pingen"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func main() {
	log.Info("Startup")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	org := os.Getenv("PINGEN_ORG")

	c := pingen.NewClient(clientID, clientSecret, true, org, context.Background())

	//demo run
	bytes, err := os.ReadFile("demo_mustermann.pdf")
	if err != nil {
		log.Fatal(err)
	}

	createData := pingen.CreateData{}
	createData.Data.Attributes.FileOriginalName = "Testupload.pdf"
	createData.Data.Attributes.AddressPosition = "left"
	createData.Data.Attributes.AutoSend = false

	result, _ := c.CreateLetter(bytes, &createData)

	time.Sleep(5 * time.Second) //wait until api catches up
	c.ListLetters()
	c.GetLetter(result.Data.ID)
	c.DeleteLetter(result.Data.ID)
}
