package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"github.com/alexflint/go-arg"
	"google.golang.org/api/option"
)

var args struct {
	DocumentID                 string `arg:"required,env:DOCUMENT_ID"`
	PaperlessToken             string `arg:"required,env:PAPERLESS_TOKEN"`
	PaperlessEndpoint          string `arg:"required,env:PAPERLESS_ENDPOINT"`
	InputFile                  string `arg:"required,env:DOCUMENT_SOURCE_PATH"`
	DocumentAIProjectID        string `arg:"required,env:DOCUMENTAI_PROJECT_ID"`
	DocumentAILocation         string `arg:"required,env:DOCUMENTAI_LOCATION"`
	DocumentAIProcessorID      string `arg:"required,env:DOCUMENTAI_PROCESSOR_ID"`
	DocumentAIProcessorVersion string `arg:"required,env:DOCUMENTAI_PROCESSOR_VERSION"`
	GoogleCredentialsFile      string `arg:"required,env:GOOGLE_APPLICATION_CREDENTIALS"`
	TestString                 string `arg:"env:TEST_STRING"`
}

func main() {
	arg.MustParse(&args)

	inputBytes, err := ioutil.ReadFile(args.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	var ocrText string
	if args.TestString == "" {
		ctx := context.Background()
		c, err := documentai.NewDocumentProcessorClient(ctx,
			option.WithCredentialsFile(args.GoogleCredentialsFile),
			option.WithEndpoint(args.DocumentAILocation+"-documentai.googleapis.com:443"))

		if err != nil {
			log.Fatalf("Cannot create client: %v", err)
		}
		defer c.Close()
		processorName := fmt.Sprintf("projects/%s/locations/%s/processors/%s/processorVersions/%s", args.DocumentAIProjectID, args.DocumentAILocation, args.DocumentAIProcessorID, args.DocumentAIProcessorVersion)
		log.Printf("Processing file: %s, ID: %s", args.InputFile, args.DocumentID)
		req := &documentaipb.ProcessRequest{
			Name: processorName,
			Source: &documentaipb.ProcessRequest_RawDocument{
				RawDocument: &documentaipb.RawDocument{
					Content:  inputBytes,
					MimeType: "application/pdf",
				},
			},
		}
		resp, err := c.ProcessDocument(ctx, req)
		if err != nil {
			log.Fatalf("Cannot process document: %v", err)
		}
		ocrText = resp.Document.GetText()
	} else {
		log.Printf("Skipping OCR, pasting %s", args.TestString)
		ocrText = args.TestString
	}
	log.Printf("Fixing OCR text...")
	httpClient := &http.Client{}
	body, errEnc := json.Marshal(map[string]string{"content": ocrText})
	if errEnc != nil {
		log.Fatalf("Cannot encode response: %v", errEnc)
	}

	patchReq, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/documents/%s/", args.PaperlessEndpoint, args.DocumentID), bytes.NewReader(body))
	patchReq.Header.Add("Content-Type", "application/json")
	patchReq.Header.Add("Accept", "application/json; version=2")
	patchReq.Header.Add("Authorization", fmt.Sprintf("Token %s", args.PaperlessToken))

	if resp, err := httpClient.Do(patchReq); err != nil {
		log.Fatalf("Cannot patch document: %v", err)
	} else {
		bb, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Patched document: %s, %s", resp.Status, string(bb))
	}
	log.Println("Sucessfully patched!")
}
