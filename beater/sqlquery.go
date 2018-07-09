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
	fieldName, fieldTypeName, schemaName, typeName string
}

type Response struct {
	fieldsMetadata []FieldMetadata
	items          [][]interface{}
}

type Body struct {
	Error         string `json:"error"`
	sessionToken  string
	response      Response
	successStatus int
}

func (q *SQuery) GenURL() string {
	sql := url.QueryEscape(q.Query.Sql)
	// UnEscape '*' ,e.g. select * from aTable;
	sql = strings.Replace(sql, "%2A", "*", -1)
	url := fmt.Sprintf("%s%spageSize=%s&cacheName=%s&qry=%s", q.Server, SQL_QUERY, string(q.Query.Size), q.Query.CacheName, sql)
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

	b, err := BodyParser(body)
	if err != nil {
		return events, err
	}

	return b.MakeEvents()

	// logp.Debug("json", "parsed body: %s", json.Marshal())
	// return events, nil
}

func BodyParser(body []byte) (b Body, err error) {
	logp.Debug(selectorDetail, "body[%s]", string(body))
	err = json.Unmarshal(body, &b)
	return b, err
}

func (bd *Body) MakeEvents() (events []*beat.Event, err error) {
	for _, item := range bd.response.items {
		for idx, v := range item {
			fieldName := bd.response.fieldsMetadata[idx].fieldName
			fieldTypeName := bd.response.fieldsMetadata[idx].fieldTypeName
			typeName := bd.response.fieldsMetadata[idx].typeName
			switch fieldTypeName {
			case javaByte:
				v, _ = v.(int8)
			case javaLong:
				v, _ = v.(int64)
			case javaString:
				v, _ = v.(string)
			case javaShort:
				v, _ = v.(int64)
			case javaInt:
				v, _ = v.(int)
			default:
				v = v
			}
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": typeName,
				},
			}
			event.Fields[fieldName] = v
			events = append(events, &event)
		}
	}
	return events, nil
}
