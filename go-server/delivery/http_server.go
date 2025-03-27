package delivery

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/frinfo702/fixer/domain/model"
	"github.com/frinfo702/fixer/usecase"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	searchInteractor *usecase.SearchInteractor
	tmpl             *template.Template
}

type PageData struct {
	Query   string
	Path    string
	Fuzzy   bool
	Results []model.SearchResult
	Error   string
}

func NewHTTPHandler(interactor *usecase.SearchInteractor) *HTTPHandler {
	h := &HTTPHandler{searchInteractor: interactor}
	h.tmpl = template.Must(template.New("page").Parse(pageHTML))
	return h
}

func (h *HTTPHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.handleIndex)
	e.GET("/search", h.handleSearch)
}

func (h *HTTPHandler) handleIndex(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/search")
}

func (h *HTTPHandler) handleSearch(c echo.Context) error {
	queryStr := c.QueryParam("query")
	fuzzy := c.QueryParam("fuzzy") == "on"
	// パスパラメータを取得する（デフォルト値は "."）
	targetDir := c.QueryParam("path")
	if targetDir == "" {
		targetDir = "."
	}

	data := PageData{
		Query: queryStr,
		Path:  targetDir,
		Fuzzy: fuzzy,
	}

	if queryStr != "" {
		searchQuery := model.SearchQuery{
			Path:  targetDir,
			Query: queryStr,
			Fuzzy: fuzzy,
		}
		results, err := h.searchInteractor.Execute(searchQuery)
		if err != nil {
			data.Error = err.Error()
		} else {
			// 結果を最大10個に制限
			if len(results) > 10 {
				data.Results = results[:10]
			} else {
				data.Results = results
			}
		}
	}

	// テンプレートをバッファにレンダリングして、HTMLレスポンスとして返す
	var buf bytes.Buffer
	if err := h.tmpl.Execute(&buf, data); err != nil {
		return c.String(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	return c.HTML(http.StatusOK, buf.String())
}

const pageHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RustySearch - Code Search</title>
    <style>
        :root {
            /* Google-inspired color palette */
            --blue-primary: #4285F4;
            --blue-hover: #1a73e8;
            --red: #EA4335;
            --yellow: #FBBC05;
            --green: #34A853;
            --gray-50: #F8F9FA;
            --gray-100: #F1F3F4;
            --gray-200: #E8EAED;
            --gray-300: #DADCE0;
            --gray-400: #BDC1C6;
            --gray-500: #9AA0A6;
            --gray-600: #80868B;
            --gray-700: #5F6368;
            --gray-800: #3C4043;
            --gray-900: #202124;
            --shadow-sm: 0 1px 2px rgba(60, 64, 67, 0.1);
            --shadow-md: 0 2px 6px rgba(60, 64, 67, 0.15);
            --shadow-lg: 0 4px 12px rgba(60, 64, 67, 0.2);
            --font-family: 'Google Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            --mono-font: 'Google Sans Mono', 'SF Mono', 'Menlo', 'Monaco', 'Courier New', monospace;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: var(--font-family);
            background-color: var(--gray-50);
            color: var(--gray-900);
            line-height: 1.6;
            min-height: 100vh;
        }

        .container {
            max-width: 800px;
            margin: 0 auto;
            padding: 1.5rem;
        }

        header {
            display: flex;
            flex-direction: column;
            align-items: center;
            padding: 2rem 0;
            text-align: center;
        }

        .logo {
            display: flex;
            align-items: center;
            margin-bottom: 0.5rem;
        }

        .logo-text {
            font-size: 2rem;
            font-weight: 500;
            color: var(--gray-900);
            letter-spacing: -0.02em;
        }

        .logo-text span:nth-child(1) { color: var(--blue-primary); }
        .logo-text span:nth-child(2) { color: var(--red); }
        .logo-text span:nth-child(3) { color: var(--yellow); }
        .logo-text span:nth-child(4) { color: var(--blue-primary); }
        .logo-text span:nth-child(5) { color: var(--green); }
        .logo-text span:nth-child(6) { color: var(--red); }
        .logo-text span:nth-child(7) { color: var(--blue-primary); }
        .logo-text span:nth-child(8) { color: var(--green); }
        .logo-text span:nth-child(9) { color: var(--yellow); }
        .logo-text span:nth-child(10) { color: var(--blue-primary); }

        .tagline {
            font-size: 0.95rem;
            color: var(--gray-600);
            margin-bottom: 2rem;
        }

        .search-container {
            width: 100%;
            max-width: 650px;
            margin: 0 auto 2rem;
        }

        .search-box {
            position: relative;
            margin-bottom: 1.5rem;
        }

        .search-input {
            width: 100%;
            padding: 0.85rem 1rem 0.85rem 3rem;
            font-size: 1rem;
            border: 1px solid var(--gray-300);
            border-radius: 24px;
            background-color: white;
            color: var(--gray-900);
            transition: all 0.2s ease;
            box-shadow: var(--shadow-sm);
        }

        .search-input:focus {
            outline: none;
            box-shadow: var(--shadow-md);
            border-color: var(--blue-primary);
        }

        .search-icon {
            position: absolute;
            left: 1rem;
            top: 50%;
            transform: translateY(-50%);
            color: var(--gray-600);
            width: 20px;
            height: 20px;
        }

        .input-group {
            margin-bottom: 1rem;
        }
        
        .input-label {
            display: block;
            font-size: 0.9rem;
            color: var(--gray-700);
            margin-bottom: 0.5rem;
        }
        
        .path-input {
            width: 100%;
            padding: 0.75rem 1rem;
            font-size: 0.95rem;
            border: 1px solid var(--gray-300);
            border-radius: 8px;
            background-color: white;
            color: var(--gray-900);
            transition: all 0.2s ease;
        }
        
        .path-input:focus {
            outline: none;
            border-color: var(--blue-primary);
            box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
        }
        
        .dir-button {
            background-color: var(--gray-100);
            color: var(--gray-800);
            border: 1px solid var(--gray-300);
            padding: 0.5rem 0.75rem;
            font-size: 0.9rem;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.2s ease;
            margin-top: 0.5rem;
        }
        
        .dir-button:hover {
            background-color: var(--gray-200);
        }
        
        .dir-help {
            font-size: 0.8rem;
            color: var(--gray-600);
            margin-top: 0.25rem;
        }

        .search-options {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 1rem;
        }

        .checkbox-container {
            display: flex;
            align-items: center;
            cursor: pointer;
        }

        input[type="checkbox"] {
            appearance: none;
            -webkit-appearance: none;
            width: 18px;
            height: 18px;
            border: 2px solid var(--gray-400);
            border-radius: 2px;
            margin-right: 8px;
            position: relative;
            cursor: pointer;
            background-color: white;
            transition: all 0.2s ease;
        }

        input[type="checkbox"]:checked {
            background-color: var(--blue-primary);
            border-color: var(--blue-primary);
        }

        input[type="checkbox"]:checked::after {
            content: "";
            position: absolute;
            left: 5px;
            top: 1px;
            width: 5px;
            height: 10px;
            border: solid white;
            border-width: 0 2px 2px 0;
            transform: rotate(45deg);
        }

        input[type="checkbox"]:focus {
            outline: none;
            box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.2);
        }

        .checkbox-label {
            font-size: 0.875rem;
            color: var(--gray-700);
        }

        .search-button {
            background-color: var(--blue-primary);
            color: white;
            border: none;
            padding: 0.75rem 1.5rem;
            font-size: 0.95rem;
            font-weight: 500;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.2s ease;
        }

        .search-button:hover {
            background-color: var(--blue-hover);
            box-shadow: var(--shadow-sm);
        }

        .search-button:focus {
            outline: none;
            box-shadow: 0 0 0 2px rgba(66, 133, 244, 0.3);
        }

        .keyboard-hint {
            display: inline-flex;
            align-items: center;
            font-size: 0.75rem;
            color: var(--gray-500);
        }

        kbd {
            background-color: var(--gray-100);
            border: 1px solid var(--gray-300);
            border-radius: 3px;
            box-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
            font-size: 0.7rem;
            font-family: inherit;
            padding: 0.1rem 0.3rem;
            margin: 0 0.1rem;
        }

        .results-container {
            margin-top: 2rem;
        }

        .results-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 1rem;
            padding-bottom: 0.5rem;
            border-bottom: 1px solid var(--gray-200);
        }

        .results-count {
            font-size: 0.95rem;
            color: var(--gray-700);
        }

        .results-count strong {
            color: var(--blue-primary);
        }

        .result {
            background-color: white;
            border-radius: 8px;
            padding: 1rem;
            margin-bottom: 1rem;
            border-left: 3px solid var(--blue-primary);
            transition: all 0.2s ease;
            box-shadow: var(--shadow-sm);
        }

        .result:hover {
            box-shadow: var(--shadow-md);
        }

        .file {
            font-family: var(--mono-font);
            font-weight: 500;
            color: var(--gray-800);
            margin-bottom: 0.5rem;
            font-size: 0.9rem;
            display: flex;
            align-items: center;
        }

        .file-icon {
            color: var(--gray-600);
            margin-right: 0.5rem;
        }

        .line-number {
            color: white;
            background-color: var(--blue-primary);
            font-family: var(--mono-font);
            font-size: 0.75rem;
            padding: 0.125rem 0.375rem;
            border-radius: 12px;
            margin-left: 0.5rem;
        }

        .line-text {
            display: block;
            margin-top: 0.5rem;
            font-family: var(--mono-font);
            font-size: 0.9rem;
            line-height: 1.5;
            padding: 0.75rem;
            background-color: var(--gray-50);
            border-radius: 4px;
            overflow-x: auto;
            color: var(--gray-800);
            border: 1px solid var(--gray-200);
        }

        .error {
            background-color: rgba(234, 67, 53, 0.1);
            color: var(--red);
            padding: 1rem;
            border-radius: 8px;
            margin-bottom: 1.5rem;
            display: flex;
            align-items: center;
        }

        .error-icon {
            margin-right: 0.5rem;
        }

        .no-results {
            text-align: center;
            padding: 3rem 0;
            color: var(--gray-600);
        }

        .no-results-icon {
            font-size: 3rem;
            margin-bottom: 1rem;
            color: var(--gray-400);
        }

        .empty-state {
            text-align: center;
            padding: 4rem 0;
        }

        .empty-state-icon {
            font-size: 3.5rem;
            color: var(--blue-primary);
            margin-bottom: 1.5rem;
        }

        .empty-state-text {
            color: var(--gray-700);
            font-size: 1.25rem;
            margin-bottom: 1.5rem;
        }

        .features-list {
            text-align: left;
            max-width: 450px;
            margin: 0 auto;
            color: var(--gray-600);
            background-color: white;
            padding: 1.5rem;
            border-radius: 8px;
            box-shadow: var(--shadow-sm);
        }

        .features-list li {
            margin-bottom: 0.75rem;
            display: flex;
            align-items: center;
        }

        .features-list li:last-child {
            margin-bottom: 0;
        }

        .feature-icon {
            color: var(--blue-primary);
            margin-right: 0.75rem;
            flex-shrink: 0;
        }

        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }

            .search-options {
                flex-direction: column;
                align-items: flex-start;
                gap: 1rem;
            }

            .keyboard-hint {
                display: none;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div class="logo">
                <h1 class="logo-text">
                    <span>R</span><span>u</span><span>s</span><span>t</span><span>y</span><span>S</span><span>e</span><span>a</span><span>r</span><span>c</span><span>h</span>
                </h1>
            </div>
            <p class="tagline">Fast and fuzzy code search powered by Rust and Go</p>
        </header>

        <div class="search-container">
            <form method="GET" action="/search">
                <div class="search-box">
                    <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="11" cy="11" r="8"></circle>
                        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                    </svg>
                    <input 
                        type="text" 
                        id="query" 
                        name="query" 
                        value="{{.Query}}" 
                        placeholder="Search code..." 
                        class="search-input" 
                        autocomplete="off"
                    />
                </div>
                
                <div class="input-group">
                    <label for="path" class="input-label">検索対象ディレクトリ（絶対パスまたは相対パス）:</label>
                    <input 
                        type="text" 
                        id="path" 
                        name="path" 
                        value="{{.Path}}" 
                        placeholder="/path/to/search" 
                        class="path-input"
                    />
                    <p class="dir-help">例: /Users/username/projects/myapp（絶対パス）または ./src（相対パス）</p>
                </div>
                
                <div class="search-options">
                    <label class="checkbox-container">
                        <input type="checkbox" name="fuzzy" id="fuzzy" {{if .Fuzzy}}checked{{end}} />
                        <span class="checkbox-label">Enable fuzzy search</span>
                    </label>
                    
                    <div class="keyboard-hint">
                        <span>Shortcut: <kbd>Ctrl</kbd> + <kbd>K</kbd> for query</span>
                    </div>
                    
                    <button type="submit" class="search-button">Search</button>
                </div>
            </form>
        </div>

        {{if .Error}}
            <div class="error">
                <svg class="error-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{{.Error}}</span>
            </div>
        {{end}}

        {{if .Results}}
            <div class="results-container">
                <div class="results-header">
                    <div class="results-count">
                        <strong>{{len .Results}}</strong> results found
                    </div>
                </div>
                
                {{range .Results}}
                    <div class="result">
                        <div class="file">
                            <svg class="file-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                                <polyline points="14 2 14 8 20 8"></polyline>
                            </svg>
                            {{.FilePath}}
                            <span class="line-number">{{.LineNumber}}</span>
                        </div>
                        <code class="line-text">{{.LineText}}</code>
                    </div>
                {{end}}
            </div>
        {{else if and (not .Error) (ne .Query "")}}
            <div class="no-results">
                <div class="no-results-icon">🔍</div>
                <p>No matches found for "<strong>{{.Query}}</strong>"</p>
                <p>Try adjusting your search query or enable fuzzy search for more flexible matching.</p>
            </div>
        {{else if and (not .Error) (eq .Query "")}}
            <div class="empty-state">
                <div class="empty-state-icon">🚀</div>
                <p class="empty-state-text">Powerful code search at your fingertips</p>
                <ul class="features-list">
                    <li>
                        <svg class="feature-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 11 12 14 22 4"></polyline>
                            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                        </svg>
                        Search for functions, variables, or specific code patterns
                    </li>
                    <li>
                        <svg class="feature-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 11 12 14 22 4"></polyline>
                            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                        </svg>
                        Enable fuzzy search for more flexible matching
                    </li>
                    <li>
                        <svg class="feature-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 11 12 14 22 4"></polyline>
                            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                        </svg>
                        Results show file location and exact line numbers
                    </li>
                    <li>
                        <svg class="feature-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="9 11 12 14 22 4"></polyline>
                            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
                        </svg>
                        Powered by Rust for lightning-fast performance
                    </li>
                </ul>
            </div>
        {{end}}
    </div>

    <script>
        // Add keyboard shortcuts
        document.addEventListener('keydown', function(e) {
            // Ctrl+K to focus query input
            if (e.key === 'k' && (e.ctrlKey || e.metaKey) && document.activeElement.tagName !== 'INPUT') {
                e.preventDefault();
                document.getElementById('query').focus();
            }
        });
    </script>
</body>
</html>
`
