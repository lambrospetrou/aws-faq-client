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
	SERVICES_MAIN string = "ec2,vpc,ebs,s3,dynamodb,glacier,rds,cloudfront,route53,directconnect,storagegateway,sns,sqs,cloudwatch,cloudtrail,config"
)

func main() {
	servicesInput := flag.String("s", "", "Please give the service name you want to read FAQ.")
	servicesMain := flag.Bool("a", false, "You will create the FAQs for all services")
	flag.Parse()

	if *servicesMain {
		FetchServices(SERVICES_MAIN)
		return
	} else if len(*servicesInput) > 0 {
		FetchServices(*servicesInput)
		return
	}
	fmt.Println("\n## Usage\n")
	fmt.Println("[-a] Fetch all the main services", SERVICES_MAIN)
	fmt.Println("[-s] Comma separated services: i.e. '-s s3,ec2'")
}

func FetchServices(servicesStr string) {
	fmt.Println("Will fetch: ", servicesStr)
	services := strings.Split(servicesStr, ",")
	for _, service := range services {
		ParseFAQ(service)
	}
}

func ParseFAQ(name string) error {
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

	if err := os.MkdirAll("faqs", 0755); err != nil {
		fmt.Println("Could not create 'faqs' directory!")
		return err
	}

	fmt.Println("Writing FAQ html has error:", ioutil.WriteFile("faqs/faq-"+name+".html", []byte(faqHtml), 0644))
	return nil
}
