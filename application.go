package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"io/ioutil"
	_ "net/http/pprof"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	SERVICES_MAIN string = "iam,kms,waf,directoryservice,ec2,vpc,ebs,elasticbeanstalk,sns,sqs,swf,s3,dynamodb,elasticache,glacier,rds,cloudfront,route53,directconnect,storagegateway,cloudformation,cloudwatch,cloudtrail,config"
)

func main() {
	servicesInput := flag.String("s", "", "Service names you want to download FAQs.")
	servicesMain := flag.Bool("a", false, "You will create the FAQs for all main services.")
	outDir := flag.String("o", "faqs", "The output dir that will contain all the downloaded FAQs")
	flag.Parse()

	if *servicesMain {
		FetchServices(SERVICES_MAIN, *outDir)
		return
	} else if len(*servicesInput) > 0 {
		FetchServices(*servicesInput, *outDir)
		return
	}
	fmt.Println("\n## Usage\n")
	fmt.Println("[-a] Fetch all the main services", SERVICES_MAIN)
	fmt.Println("[-s] Comma separated services: i.e. '-s s3,ec2'")
}

func FetchServices(servicesStr string, outDir string) {
	fmt.Println("Will fetch: ", servicesStr)
	services := strings.Split(servicesStr, ",")
	for _, service := range services {
		ParseFAQ(service, outDir)
	}
}

func ParseFAQ(name string, outDir string) error {
	fmt.Println("Fetching", name)
	doc, err := goquery.NewDocument("https://aws.amazon.com/" + strings.TrimSpace(name) + "/faqs/")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#aws-page-header").Remove()
	doc.Find("#aws-page-footer").Remove()
	doc.Find(".leftnavcontainer").Remove()
	doc.Find("main[role=\"main\"] .content .four.columns").Remove()

	faqHtml, err := doc.Html()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Println("Could not create output directory:", err)
		return err
	}

	fmt.Println("Writing FAQ html has error:", ioutil.WriteFile(outDir+"/faq-"+name+".html", []byte(faqHtml), 0644))
	return nil
}
