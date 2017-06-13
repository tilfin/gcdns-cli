package main

import (
	"fmt"
	"github.com/apcera/termtables"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app         = kingpin.New("gcdns", "Google Cloud DNS CLI")
	keyFile     = app.Flag("keyfile", "JSON key file").String()
	project     = app.Flag("project", "Project name of Google Cloud").Required().String()
	managedZone = app.Flag("mz", "Target managed zone of project").Required().String()

	listCmd = app.Command("list", "Show record sets")
	setCmd  = app.Command("set", "Set A record for target host")
	host    = setCmd.Arg("host", "FQDN of host").Required().String()
	ip      = setCmd.Arg("ip", "FQDN of host").Required().IP()
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	if *keyFile != "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *keyFile)
	}
	switch cmd {
	case listCmd.FullCommand():
		err := printRecordSets(*project, *managedZone)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	case setCmd.FullCommand():
		err := setRecord(*host, ip.String(), *project, *managedZone)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}

func setRecord(host string, ip string, project string, managedZone string) error {
	c, err := createDnsClient()
	if err != nil {
		return err
	}

	hostName := host + "."
	req := c.ResourceRecordSets.List(project, managedZone)
	req.Name(hostName)
	rrsetList, err := req.Do()
	if err != nil {
		return err
	}

	chgRrs := &dns.ResourceRecordSet{
		Kind:    "dns#resourceRecordSet",
		Name:    hostName,
		Type:    "A",
		Rrdatas: []string{ip},
	}

	change := &dns.Change{
		Kind:      "dns#change",
		Additions: []*dns.ResourceRecordSet{chgRrs},
	}

	if len(rrsetList.Rrsets) == 1 {
		change.Deletions = rrsetList.Rrsets
	}

	_, err = c.Changes.Create(project, managedZone, change).Do()
	if err != nil {
		return err
	}

	fmt.Printf("A record of %s has changed (new IP address: %s).\n", host, ip)
	return nil
}

func printRecordSets(project string, managedZone string) error {
	c, err := createDnsClient()
	if err != nil {
		return err
	}

	rrsetList, err := c.ResourceRecordSets.List(project, managedZone).Do()
	if err != nil {
		return err
	}

	table := termtables.CreateTable()
	table.AddHeaders("Name", "Type", "Value")

	for i, rrset := range rrsetList.Rrsets {
		if i > 0 {
			table.AddSeparator()
		}

		table.AddRow(rrset.Name, rrset.Type, rrset.Rrdatas[0])
		for _, rrdata := range rrset.Rrdatas[1:] {
			table.AddRow("", "", rrdata)
		}
	}

	fmt.Println(table.Render())
	return nil
}

func createDnsClient() (*dns.Service, error) {
	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.TODO()
	client, err := google.DefaultClient(ctx, dns.NdevClouddnsReadwriteScope)
	if err != nil {
		return nil, err
	}

	return dns.New(client)
}
