package ms

/**
Copyright (c) 2020 DXC Technology - Dan Hushon. All rights reserved

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc., DXC Technology nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// [START bigquery_hw_imports]
import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/olivere/ndjson"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// [END bigquery_hw_imports]

const (
	dataSetName = "boxscores"
)

//NBJson supports basic operators to create NBSJON
type NBJson interface {
	tableName() string
	marshalNBJSON(*bytes.Buffer) error
}

// sliceContains reports whether the provided string is present in the given slice of strings.
func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

// importJSONTruncate demonstrates loading data from newline-delimeted JSON data in Cloud Storage
// and overwriting/truncating data in the existing table.  Need to have ~200 rows of data to
// improve accuracy of schema determination -- via https://cloud.google.com/bigquery/docs/loading-data-cloud-storage-json
func importJSONTruncate(projectID *string, datasetID *string, tableID *string, gscReference string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, *projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	//rRef := bigquery.NewReaderSource(reader)
	rRef := bigquery.NewGCSReference(gscReference)
	rRef.SourceFormat = bigquery.JSON
	rRef.AutoDetect = true
	loader := client.Dataset(*datasetID).Table(*tableID).LoaderFrom(rRef)
	loader.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("job completed with error: %v", status.Err())
	}
	return nil
}

func createDataset(projectID string, datasetID string) (bool, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {

		return false, err
	}
	return true, nil
}

// listDatasets demonstrates iterating through the collection of datasets in a project.
func existsDataset(projectID string, testDataSetID string) (bool, error) {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("bigquery.NewClient: %v", err)
	}

	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		fmt.Printf("iterator: %s, found\n", dataset.DatasetID)
		if dataset.DatasetID == testDataSetID {
			return true, nil
		}
	}
	return false, nil
}

func (bs *BoxScore) tableName() string {
	var tn = string(bs.League)
	return ("boxscores" + tn)
}

func (sb *ScoreBoard) tableName() string {
	if len(sb.BoxScores) > 0 {
		return string("boxscores" + sb.BoxScores[0].League)
	}
	return ""
}

func (bs *BoxScore) marshalNBJSON(b *bytes.Buffer) error {
	r := ndjson.NewWriter(b)
	if err := r.Encode(bs); err != nil {
		return err
	}
	return nil
}

func (sb *ScoreBoard) marshalNBJSON(b *bytes.Buffer) error {
	r := ndjson.NewWriter(b)
	for i := 0; i < len(sb.BoxScores); i++ {
		if err := r.Encode(sb.BoxScores[i]); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string, b *bytes.Buffer) {
	//OPEN FILE TO APPEND CERT INFORMATION INTO
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l, _ := b.WriteTo(f)
	fmt.Printf("%s bytes written %d\n", filename, l)
}

//InsertRow 1 row into named project and dataset.  note that BigQuery supports
//Newline Delimited JSON (ndjson) so we need to determine if we have a singleton or an array
func InsertRow(projectID string, datasetID string, s *ScoreBoard) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	defer client.Close()
	if err != nil {
		//log.Panicf("NewClient failed: %v", err)
		return err
	}
	var b bytes.Buffer
	tableName := s.tableName()
	if err := s.marshalNBJSON(&b); err != nil {
		//log.Panicf("row Preparation failure: %v", err)
		return err
	}
	// dump buffer to file, we can then use the file to load BigQuery?
	Write(&b, &projectID, "monumental-boxes-nba", "synthetic", nil)
	// now load the written file to the bigquery tablespace
	importJSONTruncate(&projectID, &datasetID, &tableName, "gs://monumental-boxes-nba/synthetic")
	//writeFile("testout.json",&b)
	//insert data vs. reload
	/*inserter := client.Dataset(datasetID).Table(tableName).Inserter()
	inserter.IgnoreUnknownValues = true
	inserter.SkipInvalidRows = true
	if err := inserter.Put(ctx, &b); err != nil {
		return err
	}*/
	return nil
}

// DeleteDataset ... demonstrates the deletion of an empty dataset.
func deleteDataset(projectID, datasetID string) error {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}

	// To recursively delete a dataset and contents, use DeleteWithContents.
	if err := client.Dataset(datasetID).DeleteWithContents(ctx); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}
	return nil
}

//CreateTable ...
func CreateTable(projectID, datasetID string, tableID string, metadata *bigquery.TableMetadata) error { // createTablePartitioned demonstrates creating a table and specifying a time partitioning configuration.
	/*ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		//Handle Error
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	*/
	/*
		if err := t.Create(ctx,
			&bigquery.TableMetadata{
				Name:           "My New Table",
				Schema:         schema,
				ExpirationTime: time.Now().Add(24 * time.Hour),
			}); err != nil {
			// TODO: Handle error.
		}*/
	/*metadata = &bigquery.TableMetadata{
		Name:"boxscoresNBA",
		Schema: schema,
	}*/
	//if err := tableRef.Create(ctx, &metadata); err != nil {
	//	//Handle Error
	//	return err
	//}
	return nil
}

func main() {
	var bsc = ScoreBoard{
		[]BoxScore{
			BoxScore{EntityID{"2019-12-28.WSH.DET",nil,""}, "2019-12-28.WSH.DET", "NBA", Season{2019, 1},
				&Competitor{EntityID{"DET-NBA-2019",nil,""}, "Detroit Pistons", "DET", Record{0, 1, []Item{}},0,  &[]Score{}, "Detroit", "0x0000", "0xffff", true, false, nil},
				&Competitor{EntityID{"WAS-NBA-2019",nil,""}, "Washington Wizards", "WAS", Record{1, 0, []Item{}},0, &[]Score{}, "Washington", "0E3764", "e31837", true, false, nil},
				&Venue{EntityID{}, "", "Little Caesars Arena", &Address{}, 10000, true},
				&GameStatus{0.0,0,"Final","Thu, December 28th at 7:00 PM EST"},
				&[]Link{
					Link{"http://www.espn.com/nba/team/roster/_/name/det/detroit-pistons",
						[]string{"roster"}, "roster",nil, false},
				},
				&GameDetail{},
			},
			BoxScore{EntityID{"2017-02-03.TOR.BOS",nil,""}, "2017-02-03.TOR.BOS", "NBA", Season{2017, 1},
				&Competitor{EntityID{"TOR-NBA-2017",nil,""}, "Toronto Raptors", "TOR", Record{1, 0, []Item{}}, 109, &[]Score{}, "Toronto", "0x0000", "0xffff", true, false, nil},
				&Competitor{EntityID{"BOS-NBA-2017",nil,""}, "Boston Celtics", "BOS", Record{0, 1, []Item{}}, 104, &[]Score{}, "Boston", "0x0000", "0xffff", true, false, nil},
				&Venue{EntityID{}, "", "TD Garden", &Address{}, 10000, true},
				&GameStatus{0.0,0,"Final","Thu, February 3rd at 7:00 PM EST"},
				&[]Link{},
				&GameDetail{},
			},
		},
	}
	project := flag.String("project", "", "The Google Cloud Platform project ID. Required.")
	flag.Parse()

	for _, f := range []string{"project"} {
		if flag.Lookup(f).Value.String() == "" {
			log.Fatalf("The %s flag is required.", f)
		}
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, *project, option.WithCredentialsFile("../ms-testbed.json"))
	//client, err := bigquery.NewClient(ctx, *project)
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()
	exists, err := existsDataset(*project, dataSetName)
	if err != nil {
		log.Panicf("error on CreateDataset: %s\n", err.Error())
	}
	if !exists {
		_, err := createDataset(*project, dataSetName)
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok {
				if gerr.Code == 409 { // already exists
					fmt.Printf("Dataset: %s Already Exists A.OK\n", dataSetName)
				} else {
					log.Panicf("error on CreateDataset: %s\n", gerr.Error())
				}
			} else {
				log.Panicf("error on CreateDataset: %s\n", gerr.Error())
			}
		}
	}

	if err = InsertRow(*project, dataSetName, &bsc); err != nil {
		fmt.Printf("error on InsertRow: %s\n", err.Error())
		//log.Panicf("error on InsertRow: %s\n", err.Error())
		//if 404 error could do a create table and then retry?
	}

	fmt.Println("exiting")
}
