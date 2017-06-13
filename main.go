package main

import (
	"fmt"
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
	return err
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

	for _, rrset := range rrsetList.Rrsets {
		fmt.Printf("%s\t%s\n", rrset.Name, rrset.Type)
		for _, rrdata := range rrset.Rrdatas {
			fmt.Println(rrdata)
		}
		fmt.Println("-------------")
	}

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
