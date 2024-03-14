package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"med-chat-bot/pkg/cfg"
	"net/http"
)

const (
	OpTypeIndex = "index"
	FaWorld     = "faworld"
)

type CoreElkClient interface {
	Search(ctx context.Context, indexNameOrAlias string, query elastic.Query, pretty bool) (*elastic.SearchResult, error)
	CheckIndexExist(ctx context.Context, indexName string) (bool, error)
	CreateIndex(ctx context.Context, indexName string, alias *string) error
	CreateIndexIfNotExist(ctx context.Context, indexName string, alias *string) error
	DeleteIndex(ctx context.Context, indexName string) error
	DeleteIndexIfExist(ctx context.Context, indexName string) error
	BulkSaveToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error
	BulkUpdateToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error
	BulkUpsertToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error
	BulkDeleteToElasticsearch(ctx context.Context, indexName string, ids []string) error
	SaveToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error
	UpdateToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error
	UpsertToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error
	DeleteFromElasticsearch(ctx context.Context, indexName string, id string) error
	CompletionSuggesterItem(ctx context.Context, docIndex string, keyword string) ([]string, error)
	SearchWithKeywordAndSuggest(ctx context.Context, indexName string, keyword string, ids []string) (*elastic.SearchResult, error)
}

type coreElkClient struct {
	client *elastic.Client
}

func NewCoreElkClient(cfgReader *viper.Viper) (CoreElkClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	url := cfgReader.GetString(cfg.ElasticsearchUrl)
	userName := cfgReader.GetString(cfg.ElasticsearchUserName)
	password := cfgReader.GetString(cfg.ElasticsearchPassword)
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetHttpClient(httpClient),
		elastic.SetBasicAuth(userName, password),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	if err != nil {
		return nil, err
	}

	return &coreElkClient{
		client: client,
	}, nil
}

func (_this *coreElkClient) Search(ctx context.Context, indexNameOrAlias string, query elastic.Query, pretty bool) (*elastic.SearchResult, error) {
	return _this.client.Search().Index(indexNameOrAlias).Query(query).Pretty(pretty).Do(ctx)
}

func (_this *coreElkClient) CheckIndexExist(ctx context.Context, indexName string) (bool, error) {
	return elastic.NewIndicesExistsService(_this.client).Index([]string{indexName}).Do(ctx)
}

func (_this *coreElkClient) CreateIndex(ctx context.Context, indexName string, alias *string) error {
	_, err := _this.client.CreateIndex(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if alias != nil && *alias != "" {
		_, err := _this.client.Alias().Add(indexName, *alias).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (_this *coreElkClient) CreateIndexIfNotExist(ctx context.Context, indexName string, alias *string) error {
	exists, err := _this.CheckIndexExist(ctx, indexName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return _this.CreateIndex(ctx, indexName, alias)
}

func (_this *coreElkClient) DeleteIndex(ctx context.Context, indexName string) error {
	_, err := _this.client.DeleteIndex(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) DeleteIndexIfExist(ctx context.Context, indexName string) error {
	exists, err := _this.CheckIndexExist(ctx, indexName)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return _this.DeleteIndex(ctx, indexName)
}

func (_this *coreElkClient) BulkSaveToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error {
	bulk := _this.client.Bulk()
	for id, object := range data {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().OpType(OpTypeIndex).Index(indexName).Id(id).Doc(object))
	}
	_, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) BulkUpdateToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error {
	bulk := _this.client.Bulk()
	for id, object := range data {
		bulk = bulk.Add(elastic.NewBulkUpdateRequest().Index(indexName).Id(id).Doc(object))
	}
	_, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) BulkUpsertToElasticsearch(ctx context.Context, indexName string, data map[string]interface{}) error {
	bulk := _this.client.Bulk()
	for id, object := range data {
		bulk = bulk.Add(elastic.NewBulkUpdateRequest().Index(indexName).Id(id).Doc(object).DocAsUpsert(true))
	}
	_, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) BulkDeleteToElasticsearch(ctx context.Context, indexName string, ids []string) error {
	bulk := _this.client.Bulk()
	for _, id := range ids {
		bulk = bulk.Add(elastic.NewBulkDeleteRequest().Index(indexName).Id(id))
	}
	_, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) SaveToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error {
	_, err := _this.client.Index().Index(indexName).Id(id).BodyJson(data).Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) UpdateToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error {
	_, err := _this.client.Update().Index(indexName).Id(id).Doc(data).Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) UpsertToElasticsearch(ctx context.Context, indexName string, id string, data interface{}) error {
	_, err := _this.client.Update().Index(indexName).Id(id).Doc(data).DocAsUpsert(true).Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) DeleteFromElasticsearch(ctx context.Context, indexName string, id string) error {
	_, err := _this.client.Delete().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		return err
	}
	_, err = _this.client.Flush().Index(indexName).Do(ctx)
	return err
}

func (_this *coreElkClient) executeQueryForItem(c context.Context, docIndex string, query elastic.Query, offset, pageSize int64) (*elastic.SearchResult, error) {
	// result search
	ss := elastic.NewSearchSource().Query(query).
		Sort("createdAt", false).
		From(int(offset)).Size(int(pageSize))
	source, _ := ss.Source()
	jsonQuery, _ := json.Marshal(source)
	fmt.Println(string(jsonQuery))
	searchResult, err := _this.client.Search().Index(docIndex).SearchSource(ss).Pretty(false).Do(c)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (_this *coreElkClient) CompletionSuggesterItem(ctx context.Context, docIndex string, keyword string) ([]string, error) {
	var rs []string

	// Specify highlighter
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("description"))
	hl = hl.PreTags("<em>").PostTags("</em>")
	query := elastic.NewMatchQuery("description", keyword)
	ss := elastic.NewSearchSource().Query(query).Highlight(hl).From(int(0)).Size(int(20))

	source, _ := ss.Source()
	jsonQuery, _ := json.Marshal(source)
	fmt.Println(string(jsonQuery))
	searchResult, err := _this.client.Search().Index(docIndex).SearchSource(ss).Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	if searchResult.Hits != nil && searchResult.Hits.Hits != nil && len(searchResult.Hits.Hits) > 0 {
		for _, ops := range searchResult.Hits.Hits {
			if ops.Source != nil {
				var item ItemsAutoDto
				err := json.Unmarshal(ops.Source, &item)
				if err != nil {
					continue
				}
				if item.Status != "INACTIVE" {
					rs = append(rs, item.Description)
				}
			}
		}
	}
	return rs, nil
}

// SearchWithKeywordAndSuggest thực hiện tìm kiếm từ khóa, trả về gợi ý và highlight kết quả.
func (_this *coreElkClient) SearchWithKeywordAndSuggest(ctx context.Context, indexName string, keyword string, ids []string) (*elastic.SearchResult, error) {
	// Tạo một truy vấn match
	matchQuery := elastic.NewMatchQuery("text", keyword)

	// Tạo BoolQuery và thêm matchQuery
	boolQuery := elastic.NewBoolQuery().Must(matchQuery)

	// Chỉ thêm termsQuery nếu ids có độ dài lớn hơn 0
	if len(ids) > 0 {
		// Chuyển đổi mảng ids từ []string sang []interface{}
		var interfaceIds []interface{}
		for _, id := range ids {
			interfaceIds = append(interfaceIds, id)
		}

		// Tạo một truy vấn terms và thêm vào BoolQuery
		termsQuery := elastic.NewTermsQuery("id", interfaceIds...)
		boolQuery = boolQuery.Filter(termsQuery)
	}

	// Tạo một highlighter
	highlight := elastic.NewHighlight().Field("text").PreTags("<b>").PostTags("</b>")

	// Tạo một suggester
	suggester := elastic.NewCompletionSuggester("suggester").
		Text(keyword).
		Field("suggest").
		Size(5) // Số lượng gợi ý trả về

	// Thực hiện truy vấn tìm kiếm với highlight và suggester
	searchResult, err := _this.client.Search().
		Index(indexName).
		Query(boolQuery).
		Highlight(highlight).
		Suggester(suggester).
		Size(50).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	return searchResult, nil
}
