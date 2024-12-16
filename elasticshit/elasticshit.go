package elasticshit

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"log"

	"github.com/olivere/elastic/v7"
)

type SimpleRetrier struct {
	Max      int
	Interval time.Duration
}

var (
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // nolint: gosec
		},
	}

	esClient *elastic.Client
	ctx      = context.Background()
)

func (s *SimpleRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	if retry <= s.Max {
		log.Printf(`es SimpleRetrier, retry: %d, err: %+v`, retry, err)
		return s.Interval, true, nil
	}
	return 0, false, nil
}

func init() {

	fn := []elastic.ClientOptionFunc{
		elastic.SetURL("http://elastic:3be6!aXWzwTkuX@es-cn-pe333cojq000fn9ee.elasticsearch.aliyuncs.com:9200"),
		elastic.SetSniff(false),
		elastic.SetGzip(false),
		elastic.SetHttpClient(defaultHttpClient),
		elastic.SetRetrier(&SimpleRetrier{Max: 1, Interval: time.Second}),
		// elastic.SetBasicAuth("elastic", "3be6!aXWzwTkuX"),
	}

	client, err := elastic.NewClient(fn...)

	if err != nil {
		log.Fatalf("init es client error: %+v", err)
	}

	esClient = client
}

func demo1(status bool) {
	mainID := 31907

	now := time.Now()
	// script := elastic.NewScript("ctx._source['status'] = params.status").Param("status", status)
	script := elastic.NewScript("ctx._source.status = params.status; ctx._source.updated_at = params.updated_at;").Params(
		map[string]any{"status": status, "updated_at": now.Format("2006-01-02T15:04:05.999999Z")})

	_, err := esClient.UpdateByQuery().Index("knowledge-alias").Query(
		elastic.NewBoolQuery().Should(
			elastic.NewTermQuery("main_id", mainID),
			elastic.NewTermQuery("question_id", mainID),
		)).Script(script).Do(context.Background())

	log.Println("elapsed", time.Since(now))

	if err != nil {
		log.Fatalf("update by query error: %+v", err)
	}
}

func demo2(status bool) {
	ids := []int64{32005, 31907, 31908, 31909}
	now := time.Now()
	data := map[string]any{"status": status, "updated_at": now.Format("2006-01-02T15:04:05.999999Z")}
	bulkRequest := esClient.Bulk().Index("knowledge-alias")
	for _, id := range ids {
		req := elastic.NewBulkUpdateRequest().Id(fmt.Sprintf("%d", id)).Doc(data)
		bulkRequest = bulkRequest.Add(req)
	}
	_, err := bulkRequest.Refresh("true").Do(ctx)
	log.Println("elapsed", time.Since(now))

	if err != nil {
		log.Fatalf("bulk update error: %+v", err)
	}
}

func Main() {
	demo1(false)
}
