package cohere

import (
	"done-hub/common"
	"done-hub/common/config"
	"done-hub/types"
	"net/http"
)

func (p *CohereProvider) CreateRerank(request *types.RerankRequest) (*types.RerankResponse, *types.OpenAIErrorWithStatusCode) {
	documents, err := request.GetDocumentsList()
	if err != nil {
		return nil, common.ErrorWrapper(err, "invalid_documents", http.StatusBadRequest)
	}

	url, errWithCode := p.GetSupportedAPIUri(config.RelayModeRerank)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// 获取请求地址
	fullRequestURL := p.GetFullRequestURL(url)
	if fullRequestURL == "" {
		return nil, common.ErrorWrapper(nil, "invalid_jina_config", http.StatusInternalServerError)
	}

	// 获取请求头
	headers := p.GetRequestHeaders()

	rerankReq := getRerankRequest(request, documents)

	// 创建请求
	req, err := p.Requester.NewRequest(http.MethodPost, fullRequestURL, p.Requester.WithBody(rerankReq), p.Requester.WithHeader(headers))
	if err != nil {
		return nil, common.ErrorWrapper(err, "new_request_failed", http.StatusInternalServerError)
	}
	defer req.Body.Close()

	cResponse := &RerankResponse{}

	// 发送请求
	_, errWithCode = p.Requester.SendRequest(req, cResponse, false)
	if errWithCode != nil {
		return nil, errWithCode
	}

	return p.ConvertToRerank(cResponse, request)
}

func getRerankRequest(request *types.RerankRequest, documents []string) *RerankRequest {
	return &RerankRequest{
		Model:           request.Model,
		Query:           request.Query,
		TopN:            request.TopN,
		ReturnDocuments: true,
		Documents:       documents,
	}
}

func (p *CohereProvider) ConvertToRerank(response *RerankResponse, request *types.RerankRequest) (*types.RerankResponse, *types.OpenAIErrorWithStatusCode) {
	rerank := &types.RerankResponse{
		Model:   request.Model,
		Results: make([]types.RerankResult, 0),
		Usage: &types.Usage{
			PromptTokens: response.Meta.BilledUnits.SearchUnits,
			TotalTokens:  response.Meta.BilledUnits.SearchUnits,
		},
	}

	for _, result := range response.Results {
		rerankResult := types.RerankResult{
			Index:          result.Index,
			RelevanceScore: result.RelevanceScore,
			Document: types.RerankResultDocument{
				Text: result.Document.Text,
			},
		}
		rerank.Results = append(rerank.Results, rerankResult)
	}

	*p.Usage = *rerank.Usage

	return rerank, nil
}
