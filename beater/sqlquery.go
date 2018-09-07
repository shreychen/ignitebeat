package beater

// "github.com/elastic/beats/libbeat/logp"

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/shreychen/ignitebeat/config"
)

const (
	SQL_QUERY = "/ignite?cmd=qryfldexe&"

	javaLong   = "java.lang.Long"
	javaString = "java.lang.String"
	javaShort  = "java.lang.Short"
	javaByte   = "java.lang.Byte"
	javaInt    = "java.lang.Integer"
)

type SQuery struct {
	Query  *config.Query
	Server string
}

type FieldMetadata struct {
	FieldName     string `json:"fieldName"`
	FieldTypeName string `json:"fieldTypeName"`
	SchemaName    string `json:"schemaName"`
	TypeName      string `json:"typeName"`
}

type Response struct {
	FieldsMetadata []FieldMetadata `json:"fieldsMetadata,flow"`
	Items          [][]interface{} `json:"items,flow"`
}

type Body struct {
	Error         string   `json:"error"`
	SessionToken  string   `json:"sessionToken"`
	Response      Response `json:"response,flow"`
	SuccessStatus int      `json:"successStatus"`
}

func (q *SQuery) GenURL() string {
	sql := url.QueryEscape(q.Query.Sql)
	// UnEscape '*' ,e.g. select * from aTable;
	sql = strings.Replace(sql, "%2A", "*", -1)
	url := fmt.Sprintf("%s%spageSize=%s&cacheName=%s&qry=%s", q.Server, SQL_QUERY, q.Query.Size, q.Query.CacheName, sql)
	logp.Info("query url: %s", url)
	return url
}

func (q *SQuery) GenEvents() (events []*beat.Event, err error) {

	logp.Info("cache name: %s", q.Query.CacheName)

	body, err := OpenURL(q.GenURL())

	if err != nil {
		logp.Info(err.Error())
		return events, err
	}

	logp.Debug(selectorDetail, "got response body: %s", string(body))

	b, err := BodyParser(body)
	if err != nil {
		logp.Err(err.Error())
		return events, err
	}

	logp.Debug("json", "parsed body %v: ", b)

	return b.MakeEvents()

	// logp.Debug("json", "parsed body: %s", json.Marshal())
	// return events, nil
}

func BodyParser(body []byte) (b Body, err error) {
	err = json.Unmarshal(body, &b)
	return b, err
}

func (bd *Body) MakeEvents() (events []*beat.Event, err error) {
	logp.Info("making events")
	for _, item := range bd.Response.Items {
		event := beat.Event{
			Timestamp: time.Now(),
			Fields:    common.MapStr{},
		}
		for idx, v := range item {
			fieldName := bd.Response.FieldsMetadata[idx].FieldName
			fieldName = strings.Replace(fieldName, "-", "_", -1)
			fieldName = strings.ToLower(fieldName)

			// typeName := bd.Response.FieldsMetadata[idx].TypeName
			// fieldTypeName := bd.Response.FieldsMetadata[idx].FieldTypeName
			// switch fieldTypeName {
			// case javaByte:
			// 	v, _ = v.(int8)
			// case javaLong:
			// 	v, _ = v.(int64)
			// case javaString:
			// 	v, _ = v.(string)
			// case javaShort:
			// 	v, _ = v.(int64)
			// case javaInt:
			// 	v, _ = v.(int)
			// default:
			// 	v = v
			// }
			// event.Fields["type"] = typeName
			event.Fields[fieldName] = v
		}
		events = append(events, &event)
	}

	logp.Info("%d events made by sql query", len(events))
	return events, nil
}
