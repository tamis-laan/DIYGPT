package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"services/pkg/schema"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tmaxmax/go-sse"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Read configuration
func init() {

	// Define command-line flags
	pflag.String("kafka.bootstrap", "localhost:9092", "Kafka bootstrap servers")
	pflag.String("kafka.clientid", "wiki producer", "The kafka client id")
	pflag.String("url", "https://stream.wikimedia.org/v2/stream/test", "The stream to scrape from")
	pflag.String("topic", "wiki.test", "The topic to publish to")
	pflag.String("since", "", "From when to start the stream")
	pflag.String("filter.database", "^(enwiki|nlwiki)$", "Regex filter for database property")
	pflag.String("filter.title", "^[^:]*$", "Regex filter for title property")

	// Parse command-line flags
	pflag.Parse()

	// Bind command-line flags to Viper
	viper.BindPFlags(pflag.CommandLine)

	// Replace . with _ in environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Make configurable from environment variables
	viper.AutomaticEnv()

}

func unpack(jsonData string) (schema.WikiPageCreateLocal, error) {

	// Create struct
	var origin schema.WikiPageCreateOrigin

	// Unpack json
	err := json.Unmarshal([]byte(jsonData), &origin)

	// Handle error
	if err != nil {
		return schema.WikiPageCreateLocal{}, err
	}

	// Create the new structure
	novel := schema.WikiPageCreateLocal{
		Uri:                origin.Meta.URI,
		Domain:             origin.Meta.Domain,
		Database:           origin.Database,
		PageID:             origin.PageID,
		PageTitle:          origin.PageTitle,
		PageNamespace:      origin.PageNamespace,
		UserIsBot:          origin.Performer.UserIsBot,
		UserID:             origin.Performer.UserID,
		UserRegistrationDt: origin.Performer.UserRegistrationDt,
		UserEditCount:      origin.Performer.UserEditCount,
		Comment:            origin.Comment,
		ParsedComment:      origin.ParsedComment,
		Dt:                 origin.Dt,
	}

	// Return
	return novel, nil

}

// Scrape url
func scrape(client *http.Client, url string) (string, error) {

	// Create request
	req, err := http.NewRequest("GET", url, nil)

	// Handle errors
	if err != nil {
		return "", err
	}

	// Send the request
	res, err := client.Do(req)

	// Handle errors
	if err != nil {
		return "", err
	}

	// Close body
	defer res.Body.Close()

	// Handle incorrect response
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Non 200 response %d", res.StatusCode))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)

	// Handle error
	if err != nil {
		return "", err
	}

	// Remove all <table> elements
	doc.Find("table").Remove()

	// Remove all script script and meta stuff
	doc.Find("script, style, noscript, .metadata").Remove()

	// Remove footer
	doc.Find("div.printfooter").Remove()

	// Remove references
	doc.Find("sup.reference").Remove()

	// Remove reference list
	doc.Find("div.reflist, h2 span#References").Remove()

	// Remove edit buttons
	doc.Find("span.mw-editsection").Remove()

	// Extract content
	content := doc.Find("div#mw-content-text.mw-body-content").Text()

	// Strip tabs
	content = strings.ReplaceAll(content, `\t`, "")

	// Regex to find multiple new lines in a row
	newlineRegex, err := regexp.Compile(`(\n\s*)+`)

	// Handle error
	if err != nil {
		log.Fatal(err)
	}

	// Replace multiple new lines with a single new line
	content = newlineRegex.ReplaceAllString(content, "\n")

	// Reject empty pages
	if content == "" {
		return "", errors.New("Empty page content!")
	}

	// Return extracted content
	return content, nil

}

// Start producer
func main() {

	// Create producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.bootstrap"),
		"client.id":         viper.GetString("kafka.clientid"),
		"acks":              "all",
	})

	// Error check
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		log.Fatal(err)
		os.Exit(1)
	}

	// Close producer
	defer p.Close()

	// Log to console
	log.Printf("%+v\n", p)

	// Define the topic
	topic := viper.GetString("topic")

	// Create a channel that tells us if something is delivered or not
	deliveryChan := make(chan kafka.Event, 10000)

	// Create http client
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, viper.GetString("url"), http.NoBody)

	// Handle error
	if err != nil {
		log.Fatalln(err)
	}

	// Start a connection
	conn := sse.DefaultClient.NewConnection(req)

	// Define a regexp pattern for matching correct databases
	filter_database, err := regexp.Compile(viper.GetString("filter.database"))

	// Handle error
	if err != nil {
		log.Fatal(err)
	}

	// Define a regexp pattern for matching correct databases
	filter_title, err := regexp.Compile(viper.GetString("filter.title"))

	// Handle error
	if err != nil {
		log.Fatal(err)
	}

	// Scraper http client
	cclient := http.Client{}

	// Subscribe to message event
	conn.SubscribeEvent("message", func(event sse.Event) {

		// Extract out relevant values into struct
		content, err := unpack(event.Data)

		// Only focus on nlwiki
		if content.UserIsBot == true ||
			filter_database.MatchString(content.Database) == false ||
			filter_title.MatchString(content.PageTitle) == false {
			return
		}

		// Handle error
		if err != nil {
			log.Fatal(err)
		}

		// Scrape page
		page, err := scrape(&cclient, content.Uri)

		// Handle error
		if err != nil {
			log.Println(err)
			return
		}

		// Attach page
		content.Page = page

		// Turn page into json
		contentJson, err := json.MarshalIndent(content, "", "  ")

		// Handle error
		if err != nil {
			log.Fatal(err)
		}

		// Construct message
		kmessage := kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte("mykey"),
			Value:          []byte(contentJson),
		}

		// Send a message
		err = p.Produce(&kmessage, deliveryChan)

		// Error check
		if err != nil {
			log.Fatal(err)
		}

		// Retrieve message
		e := <-deliveryChan
		m := e.(*kafka.Message)

		// Print message
		log.Println(m, string(m.Value))
	})

	// Start connection
	conn.Connect()

}
