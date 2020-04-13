package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// hati2 bagian mapping klo err ident aja kelar udah semua
const MappingExample = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"ticket":{
			"properties":{
				"name":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`

// Example for struct example elasitic search
type Example struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}

func TestElastic(t *testing.T) {
	g := NewServiceGenerator()
	times := time.Now()
	ctx := context.Background()
	es := NewEsCore("http://fec-ticketing-stag-es.statefulset.svc.cluster.local:9200", "ticketing", MappingExample, "ticket")

	errCreateIndex := es.CreateIndex(ctx)
	assert.NoError(t, errCreateIndex)

	ID1 := g.Do(10)
	ID2 := g.Do(10)
	ID3 := "3"

	docs1 := &Example{Name: "Maulana", Message: "Amin paling serius", Created: times}
	errAdd := es.AddDocument(ctx, ID1, docs1)
	assert.NoError(t, errAdd)

	docs2 := &Example{Name: "Izzudin", Message: "gak serius", Created: times}
	errAdd2 := es.AddDocument(ctx, ID2, docs2)
	assert.NoError(t, errAdd2)

	docs3 := &Example{Name: "Burhanudin", Message: "jalan2", Created: times}
	errAdd3 := es.AddDocument(ctx, ID3, docs3)
	assert.NoError(t, errAdd3)

	// query with no params
	result, total, errQueryAll := es.Query(ctx, nil, nil, nil, 0, 10, "name", true)
	assert.NoError(t, errQueryAll)
	assert.Equal(t, 3, total)
	assert.NotEmpty(t, result)

	for _, hit := range result {
		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
		var t Example
		err := json.Unmarshal(*hit.Source, &t)
		if err != nil {
			// Deserialization failed
		}
		fmt.Printf("Tweet by %s: %s\n", t.Name, t.Message)
	}
	// query with params
	queryParam := es.GenerateQueryTerms()

	// queryParam = append(queryParam, es.GetQueryTerm("message", "serius")) nah ini full text search karena di mapping nya sbg text
	queryParam = append(queryParam, es.GetQueryTerm("name", "Maulana"))   // ini harus exact sama soalnya dia tipe mappingnya keyword
	queryParam = append(queryParam, es.GetQueryTerm("message", "serius")) // ini harus exact sama soalnya dia tipe mappingnya keyword

	// case di atas jadi 2 contoh condisi name maulana or message serius jadi yang nampil 2 rec
	// contoh yang di bawah untuk mengambil satu kondisi doank
	resultParam, totalWParam, errQueryWParam := es.Query(ctx, queryParam[0], nil, nil, 0, 10, "name", true)
	assert.NoError(t, errQueryWParam)
	assert.Equal(t, 1, totalWParam)
	assert.NotEmpty(t, resultParam)

	// ini multiple condition
	// queryParam2 := es.GenerateQueryTerms()
	// queryParam2 = append(queryParam2, es.GetQueryTerm("name", "Burhanudin"))
	queryParamFilter := es.GetBoolQuery(es.GetQueryTerm("message", "jalan2"), es.GetNewMatchQuery("name", "Burhanudin"))

	resultParam2, totalWParam2, errQueryWParam2 := es.Query(ctx, nil, queryParamFilter, nil, 0, 10, "name", true)
	assert.NoError(t, errQueryWParam2)
	assert.Equal(t, 1, totalWParam2)
	assert.NotEmpty(t, resultParam2)

	for _, hit := range resultParam2 {
		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
		var t Example
		err := json.Unmarshal(*hit.Source, &t)
		if err != nil {
			// Deserialization failed
		}
		fmt.Printf("With param by %s: %s\n", t.Name, t.Message)
	}

	// for update
	multiscript := es.GenerateScriptLines()
	// script := es.GetScriptLine("name = params.name", "name", "udin aja")
	script2 := es.GetScriptLine("message = params.message", "message", "ke sorean kayaknya bro")
	multiscript = append(multiscript, script2)

	ver, errUpdate := es.UpdateDocument(ctx, ID3, multiscript, map[string]interface{}{"name": "", "message": ""})

	assert.NoError(t, errUpdate)
	log.Println(ver)
	result, total, errQueryAll = es.Query(ctx, es.GetQueryTerm("message", "sorean"), nil, nil, 0, 10, "name", true)

	for _, hit := range result {
		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
		var t Example
		err := json.Unmarshal(*hit.Source, &t)
		if err != nil {
			// Deserialization failed
		}
		fmt.Printf("updates %s: %s\n", t.Name, t.Message)
	}

	assert.Equal(t, 1, total)

	// delete data
	errDelete := es.DeleteDocumentByID(ctx, ID3)
	assert.NoError(t, errDelete)

	result, total, errQueryAll = es.Query(ctx, nil, nil, nil, 0, 10, "name", true)
	assert.NoError(t, errQueryAll)
	assert.Equal(t, 2, total)
	assert.NotEmpty(t, result)

	err := es.DeleteIndex(ctx)
	assert.NoError(t, err)

}
