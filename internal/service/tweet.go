package service

import (
	"bytes"
	"context"
	"daily/internal/db"
	"daily/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// ParseSiteInfoFromURL extracts content from a provided URL and creates a SiteInfo entity
func ParseSiteInfoFromURL(urlStr string) (*entity.SiteInfo, error) {
	// Validate URL
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}

	// Generate a random ID for the site info
	rand.Seed(time.Now().UnixNano())
	randomID := fmt.Sprintf("%d", rand.Int63n(9000000000000000000)+1000000000000000000)

	// Fetch the page content
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, errors.New("failed to fetch page content")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read page content")
	}
	
	htmlContent := string(body)
	
	// Extract metadata from the page
	title := extractMetaTag(htmlContent, "og:title")
	if title == "" {
		title = extractTitle(htmlContent)
	}
	
	description := extractMetaTag(htmlContent, "og:description")
	if description == "" {
		description = extractMetaTag(htmlContent, "description")
	}
	
	keywords := extractMetaTag(htmlContent, "keywords")
	
	imageURL := extractMetaTag(htmlContent, "og:image")
	if imageURL == "" {
		// Try to find first image in the content
		imageURL = extractFirstImage(htmlContent)
	}
	
	// Create initial SiteInfo entity
	initialSiteInfo := &entity.SiteInfo{
		Id:          randomID,
		Title:       title,
		Keywords:    keywords,
		Description: description,
		OriginUrl:   urlStr,
		ImageUrl:    imageURL,
		CreatedAt:   time.Now(),
	}
	
	// Enhance content using AI model
	enhancedSiteInfo, err := enhanceContentWithAI(initialSiteInfo, htmlContent)
	if err != nil {
		// If enhancement fails, return the original site info
		fmt.Printf("Warning: AI enhancement failed: %v. Using original content.\n", err)
		return initialSiteInfo, nil
	}
	
	return enhancedSiteInfo, nil
}

// SaveSiteInfo saves a SiteInfo entity to the database
func SaveSiteInfo(siteInfo *entity.SiteInfo) error {
	engine := db.NewDB()
	defer engine.Close()
	
	_, err := engine.Insert(siteInfo)
	return err
}

// Helper function to extract meta tag content from HTML
func extractMetaTag(html, property string) string {
	// Try with property attribute first (Open Graph)
	re := regexp.MustCompile(`<meta\s+property=["']` + property + `["']\s+content=["']([^"']+)["']`)
	matches := re.FindStringSubmatch(html)
	
	if len(matches) >= 2 {
		return cleanHtmlEntities(matches[1])
	}
	
	// Try with name attribute
	re = regexp.MustCompile(`<meta\s+name=["']` + property + `["']\s+content=["']([^"']+)["']`)
	matches = re.FindStringSubmatch(html)
	
	if len(matches) >= 2 {
		return cleanHtmlEntities(matches[1])
	}
	
	return ""
}

// Helper function to extract title from HTML
func extractTitle(html string) string {
	re := regexp.MustCompile(`<title>([^<]+)</title>`)
	matches := re.FindStringSubmatch(html)
	
	if len(matches) >= 2 {
		return cleanHtmlEntities(matches[1])
	}
	
	return "Untitled"
}

// Helper function to extract first image from HTML
func extractFirstImage(html string) string {
	re := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches := re.FindStringSubmatch(html)
	
	if len(matches) >= 2 {
		return matches[1]
	}
	
	return ""
}

// Helper function to clean HTML entities
func cleanHtmlEntities(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	return text
}

// enhanceContentWithAI uses AI model to enhance the content
func enhanceContentWithAI(siteInfo *entity.SiteInfo, htmlContent string) (*entity.SiteInfo, error) {
	// Get AI provider from environment variables
	aiProvider := os.Getenv("AI_PROVIDER")
	
	// Default to local if not specified
	if aiProvider == "" || (aiProvider != "openai" && aiProvider != "gemini" && aiProvider != "local") {
		aiProvider = "local"
	}
	
	// Extract text content from HTML for better context
	textContent := extractTextContent(htmlContent)
	
	// Prepare prompt content
	systemPrompt := "你是一个内容优化助手，专门用于优化网页内容以生成高质量的每日精选卡片。你的任务是：\n" +
		"1. 优化描述信息，使其更加简洁、准确、吸引人\n" +
		"2. 提取或生成适合的关键词（最多5个）\n" +
		"3. 如果提供的图片URL不合适或为空，建议一个更好的图片URL\n" +
		"4. 如果标题不够吸引人，优化标题"
	
	userPrompt := fmt.Sprintf("请优化以下网页内容，生成高质量的每日精选卡片：\n\n" +
		"标题：%s\n" +
		"描述：%s\n" +
		"关键词：%s\n" +
		"图片URL：%s\n" +
		"原始URL：%s\n\n" +
		"网页文本内容摘要：\n%s\n\n" +
		"请使用以下格式返回结果：\n" +
		"标题：优化后的标题\n" +
		"描述：优化后的描述\n" +
		"关键词：关键词1,关键词2,关键词3\n" +
		"图片URL：如果原始URL不合适，请提供更好的建议",
		siteInfo.Title,
		siteInfo.Description,
		siteInfo.Keywords,
		siteInfo.ImageUrl,
		siteInfo.OriginUrl,
		textContent)
	
	var aiContent string
	var err error
	
	// Call the appropriate AI API based on provider
	switch aiProvider {
	case "openai":
		aiContent, err = callOpenAIAPI(systemPrompt, userPrompt)
	case "gemini":
		aiContent, err = callGeminiAPI(systemPrompt, userPrompt)
		// 如果Gemini API调用失败，回退到本地增强
		if err != nil {
			fmt.Printf("Warning: Gemini API failed: %v. Falling back to local enhancement.\n", err)
			aiContent, err = enhanceContentLocally(siteInfo, textContent)
		}
	default: // local
		aiContent, err = enhanceContentLocally(siteInfo, textContent)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Extract enhanced content using regex
	enhancedTitle := extractAIContent(aiContent, "标题[:：]\\s*(.+)")
	enhancedDescription := extractAIContent(aiContent, "描述[:：]\\s*(.+)")
	enhancedKeywords := extractAIContent(aiContent, "关键词[:：]\\s*(.+)")
	enhancedImageURL := extractAIContent(aiContent, "图片URL[:：]\\s*(.+)")
	
	// Create enhanced SiteInfo
	enhancedSiteInfo := &entity.SiteInfo{
		Id:          siteInfo.Id,
		Title:       siteInfo.Title,        // Default to original
		Keywords:    siteInfo.Keywords,     // Default to original
		Description: siteInfo.Description,  // Default to original
		OriginUrl:   siteInfo.OriginUrl,    // Keep original URL
		ImageUrl:    siteInfo.ImageUrl,     // Default to original
		CreatedAt:   siteInfo.CreatedAt,    // Keep original timestamp
	}
	
	// Update with enhanced content if available
	if enhancedTitle != "" {
		enhancedSiteInfo.Title = enhancedTitle
	}
	
	if enhancedDescription != "" {
		enhancedSiteInfo.Description = enhancedDescription
	}
	
	if enhancedKeywords != "" {
		enhancedSiteInfo.Keywords = enhancedKeywords
	}
	
	if enhancedImageURL != "" && enhancedImageURL != "无" && enhancedImageURL != "N/A" {
		enhancedSiteInfo.ImageUrl = enhancedImageURL
	}
	
	return enhancedSiteInfo, nil
}

// callOpenAIAPI calls the OpenAI API to enhance content
func callOpenAIAPI(systemPrompt, userPrompt string) (string, error) {
	// Get API key and URL from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiURL := os.Getenv("OPENAI_API_URL")
	model := os.Getenv("OPENAI_MODEL")
	
	// Check if API key and URL are set
	if apiKey == "" || apiURL == "" {
		return "", errors.New("OpenAI API key or URL not set in environment variables")
	}
	
	// If model is not set, use a default model
	if model == "" {
		model = "gpt-4-turbo"
	}
	
	// Prepare request payload
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": userPrompt,
			},
		},
		"temperature": 0.7,
		"max_tokens": 1000,
	}
	
	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API returned error: %s", respBody)
	}
	
	// Parse response
	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	err = json.Unmarshal(respBody, &aiResponse)
	if err != nil {
		return "", err
	}
	
	// Check if we have a response
	if len(aiResponse.Choices) == 0 || aiResponse.Choices[0].Message.Content == "" {
		return "", errors.New("OpenAI API returned empty response")
	}
	
	return aiResponse.Choices[0].Message.Content, nil
}

// callGeminiAPI calls the Gemini API to enhance content using the official generative-ai-go library
func callGeminiAPI(systemPrompt, userPrompt string) (string, error) {
	// Get API key from environment variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	modelName := os.Getenv("GEMINI_MODEL")
	
	// Check if API key is set
	if apiKey == "" {
		return "", errors.New("Gemini API key not set in environment variables")
	}
	
	// If model is not set, use a default model
	if modelName == "" {
		modelName = "gemini-pro"
	}
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Initialize the Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()
	
	// Get the model
	model := client.GenerativeModel(modelName)
	
	// Set generation parameters
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(1000)
	
	// Combine system prompt and user prompt for Gemini
	combinedPrompt := systemPrompt + "\n\n" + userPrompt
	
	// Create a prompt
	prompt := genai.Text(combinedPrompt)
	
	// Generate content
	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}
	
	// Check if we have a response
	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("Gemini API returned empty response")
	}
	
	// Extract the text from the response
	result, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", errors.New("unexpected response format from Gemini API")
	}
	
	return string(result), nil
}

// extractAIContent extracts content from AI response using regex
func extractAIContent(aiContent, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(aiContent)
	
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	
	return ""
}

// extractTextContent extracts readable text content from HTML
func extractTextContent(html string) string {
	// Remove script and style tags and their content
	html = regexp.MustCompile(`(?s)<script.*?</script>`).ReplaceAllString(html, "")
	html = regexp.MustCompile(`(?s)<style.*?</style>`).ReplaceAllString(html, "")
	
	// Remove HTML tags
	html = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(html, " ")
	
	// Replace multiple spaces with a single space
	html = regexp.MustCompile(`\s+`).ReplaceAllString(html, " ")
	
	// Trim and limit length
	html = strings.TrimSpace(html)
	if len(html) > 3000 {
		html = html[:3000] + "..."
	}
	
	return html
}

// enhanceContentLocally provides a local implementation for content enhancement
// when external AI APIs are not available
func enhanceContentLocally(siteInfo *entity.SiteInfo, textContent string) (string, error) {
	// 提取或生成关键词
	keywords := extractKeywords(siteInfo.Title, textContent)
	
	// 优化描述
	description := optimizeDescription(siteInfo.Description, textContent)
	
	// 优化标题
	title := optimizeTitle(siteInfo.Title)
	
	// 构建格式化的响应
	response := fmt.Sprintf("标题：%s\n描述：%s\n关键词：%s\n图片URL：%s",
		title,
		description,
		keywords,
		siteInfo.ImageUrl)
	
	return response, nil
}

// extractKeywords extracts keywords from title and content
func extractKeywords(title, content string) string {
	// 简单实现：从标题中提取关键词
	words := strings.Fields(title)
	
	// 过滤掉常见的停用词和短词
	filteredWords := []string{}
	for _, word := range words {
		if len(word) > 2 && !isStopWord(word) {
			filteredWords = append(filteredWords, word)
		}
	}
	
	// 如果提取的关键词不足，从内容中提取一些
	if len(filteredWords) < 3 {
		contentWords := strings.Fields(content)
		for _, word := range contentWords {
			if len(filteredWords) >= 5 {
				break
			}
			if len(word) > 3 && !isStopWord(word) && !contains(filteredWords, word) {
				filteredWords = append(filteredWords, word)
			}
		}
	}
	
	// 限制关键词数量
	if len(filteredWords) > 5 {
		filteredWords = filteredWords[:5]
	}
	
	return strings.Join(filteredWords, ",")
}

// isStopWord checks if a word is a common stop word
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "a": true, "to": true, "of": true, "in": true,
		"is": true, "that": true, "for": true, "on": true, "with": true,
		"的": true, "了": true, "和": true, "是": true, "在": true,
		"不": true, "有": true, "一": true, "个": true, "上": true,
	}
	return stopWords[strings.ToLower(word)]
}

// contains checks if a string slice contains a specific string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// optimizeDescription optimizes the description
func optimizeDescription(description, content string) string {
	// 如果描述为空，从内容中提取
	if description == "" {
		if len(content) > 150 {
			description = content[:150] + "..."
		} else {
			description = content
		}
	}
	
	// 如果描述太长，截断它
	if len(description) > 200 {
		description = description[:197] + "..."
	}
	
	return description
}

// optimizeTitle optimizes the title
func optimizeTitle(title string) string {
	// 如果标题太长，截断它
	if len(title) > 100 {
		title = title[:97] + "..."
	}
	
	return title
}
