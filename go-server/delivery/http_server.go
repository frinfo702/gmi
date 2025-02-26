package delivery

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/frinfo702/rustysearch/domain/model"
	"github.com/frinfo702/rustysearch/usecase"
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
	path := c.QueryParam("path")
	queryStr := c.QueryParam("query")
	fuzzy := c.QueryParam("fuzzy") == "on"

	if path == "" {
		path = "."
	}

	data := PageData{
		Query: queryStr,
		Path:  path,
		Fuzzy: fuzzy,
	}

	if queryStr != "" {
		searchQuery := model.SearchQuery{
			Path:  path,
			Query: queryStr,
			Fuzzy: fuzzy,
		}
		results, err := h.searchInteractor.Execute(searchQuery) // これを実装する
		if err != nil {
			data.Error = err.Error()
		} else {
			data.Results = results
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
    <title>RustySearch</title>
    <style>
        :root {
            /* Color palette */
            --black-primary: #191919;
            --black-secondary: #262625;
            --black-tertiary: #40403E;
            --gray-primary: #666663;
            --gray-secondary: #91918D;
            --gray-tertiary: #BFBFBA;
            --white-primary: #E5E4DF;
            --white-secondary: #F0F0EB;
            --white-tertiary: #FAFAF7;
            --orange-primary: #CC785C;
            --orange-secondary: #D4A27F;
            --orange-tertiary: #EBDBBC;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            background-color: var(--white-tertiary);
            color: var(--black-primary);
            line-height: 1.6;
            padding: 0;
            margin: 0;
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }

        header {
            padding: 2rem 0;
        }

        h1 {
            font-size: 2.5rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
            color: var(--black-primary);
            letter-spacing: -0.03em;
        }

        .tagline {
            font-size: 1.125rem;
            color: var(--gray-primary);
            margin-bottom: 2rem;
        }

        .search-container {
            background-color: var(--white-secondary);
            border-radius: 12px;
            padding: 2rem;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
            margin-bottom: 2.5rem;
        }

        .form-group {
            margin-bottom: 1.5rem;
            position: relative;
        }

        label {
            display: block;
            font-size: 0.875rem;
            font-weight: 500;
            margin-bottom: 0.5rem;
            color: var(--gray-primary);
        }

        input[type="text"] {
            width: 100%;
            padding: 0.75rem 1rem;
            font-size: 1rem;
            border: 1px solid var(--gray-tertiary);
            border-radius: 6px;
            background-color: var(--white-tertiary);
            color: var(--black-primary);
            transition: all 0.2s ease;
        }

        input[type="text"]:focus {
            outline: none;
            border-color: var(--orange-secondary);
            box-shadow: 0 0 0 3px rgba(212, 162, 127, 0.2);
        }

        .checkbox-group {
            display: flex;
            align-items: center;
            margin-bottom: 1.5rem;
        }

        .checkbox-container {
            display: flex;
            align-items: center;
            cursor: pointer;
        }

        input[type="checkbox"] {
            appearance: none;
            -webkit-appearance: none;
            width: 20px;
            height: 20px;
            border: 1px solid var(--gray-tertiary);
            border-radius: 4px;
            margin-right: 8px;
            position: relative;
            cursor: pointer;
            background-color: var(--white-tertiary);
        }

        input[type="checkbox"]:checked {
            background-color: var(--orange-primary);
            border-color: var(--orange-primary);
        }

        input[type="checkbox"]:checked::after {
            content: "";
            position: absolute;
            left: 6px;
            top: 2px;
            width: 5px;
            height: 10px;
            border: solid white;
            border-width: 0 2px 2px 0;
            transform: rotate(45deg);
        }

        input[type="checkbox"]:focus {
            outline: none;
            box-shadow: 0 0 0 3px rgba(212, 162, 127, 0.2);
        }

        .checkbox-label {
            font-size: 0.875rem;
            color: var(--gray-primary);
        }

        button {
            background-color: var(--orange-primary);
            color: var(--white-tertiary);
            border: none;
            padding: 0.75rem 1.5rem;
            font-size: 1rem;
            font-weight: 500;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s ease;
        }

        button:hover {
            background-color: #b86a4f;
        }

        button:focus {
            outline: none;
            box-shadow: 0 0 0 3px rgba(204, 120, 92, 0.3);
        }

        .results-container {
            margin-top: 2rem;
        }

        .results-header {
            font-size: 1.25rem;
            font-weight: 600;
            margin-bottom: 1.5rem;
            color: var(--black-primary);
        }

        .results-count {
            color: var(--orange-primary);
            font-weight: 700;
        }

        .result {
            background-color: var(--white-secondary);
            border-radius: 8px;
            padding: 1.25rem;
            margin-bottom: 1rem;
            border-left: 3px solid var(--orange-secondary);
            transition: all 0.2s ease;
        }

        .result:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(0, 0, 0, 0.06);
        }

        .file {
            font-family: 'SF Mono', 'Menlo', 'Monaco', 'Courier New', monospace;
            font-weight: 600;
            color: var(--black-primary);
            margin-bottom: 0.5rem;
            font-size: 0.9rem;
        }

        .line-number {
            color: var(--orange-primary);
            font-family: 'SF Mono', 'Menlo', 'Monaco', 'Courier New', monospace;
            font-size: 0.8rem;
            padding: 0.125rem 0.375rem;
            background-color: var(--orange-tertiary);
            border-radius: 4px;
            margin-left: 0.5rem;
        }

        .line-text {
            display: block;
            margin-top: 0.5rem;
            font-family: 'SF Mono', 'Menlo', 'Monaco', 'Courier New', monospace;
            font-size: 0.9rem;
            line-height: 1.5;
            padding: 0.75rem;
            background-color: var(--white-tertiary);
            border-radius: 4px;
            overflow-x: auto;
            color: var(--black-secondary);
        }

        .error {
            background-color: rgba(255, 0, 0, 0.1);
            color: #d63031;
            padding: 1rem;
            border-radius: 6px;
            margin-bottom: 1.5rem;
        }

        .no-results {
            text-align: center;
            padding: 3rem 0;
            color: var(--gray-secondary);
        }

        .no-results-icon {
            font-size: 2.5rem;
            margin-bottom: 1rem;
            color: var(--gray-tertiary);
        }

        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }

            .search-container {
                padding: 1.5rem;
            }
        }

        /* Keyboard shortcuts hint */
        .keyboard-shortcuts {
            position: absolute;
            right: 10px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 0.75rem;
            color: var(--gray-secondary);
            background-color: var(--white-primary);
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
        }

        kbd {
            background-color: var(--white-tertiary);
            border: 1px solid var(--gray-tertiary);
            border-radius: 3px;
            box-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
            font-size: 0.7rem;
            font-family: inherit;
            padding: 0.1rem 0.3rem;
            margin: 0 0.1rem;
        }

        /* Empty state */
        .empty-state {
            text-align: center;
            padding: 4rem 0;
        }

        .empty-state-icon {
            font-size: 3rem;
            color: var(--gray-tertiary);
            margin-bottom: 1rem;
        }

        .empty-state-text {
            color: var(--gray-secondary);
            font-size: 1.125rem;
            margin-bottom: 1.5rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>RustySearch</h1>
            <p class="tagline">Fast and fuzzy code search powered by Rust and Go</p>
        </header>

        <div class="search-container">
            <form method="GET" action="/search">
                <div class="form-group">
                    <label for="path">Directory Path</label>
                    <input type="text" id="path" name="path" value="{{.Path}}" placeholder="Enter repository or folder path..." />
                    <div class="keyboard-shortcuts"><kbd>/</kbd> to focus</div>
                </div>
                
                <div class="form-group">
                    <label for="query">Search Query</label>
                    <input type="text" id="query" name="query" value="{{.Query}}" placeholder="Type your search query..." />
                    <div class="keyboard-shortcuts"><kbd>Ctrl</kbd> + <kbd>K</kbd> to focus</div>
                </div>
                
                <div class="checkbox-group">
                    <label class="checkbox-container">
                        <input type="checkbox" name="fuzzy" id="fuzzy" {{if .Fuzzy}}checked{{end}} />
                        <span class="checkbox-label">Enable fuzzy search (more flexible matching)</span>
                    </label>
                </div>
                
                <button type="submit">Search Files</button>
            </form>
        </div>

        {{if .Error}}
            <div class="error">
                <strong>Error:</strong> {{.Error}}
            </div>
        {{end}}

        {{if .Results}}
            <div class="results-container">
                <h3 class="results-header">Found <span class="results-count">{{len .Results}}</span> result(s):</h3>
                
                {{range .Results}}
                    <div class="result">
                        <div class="file">
                            {{.FilePath}}
                            <span class="line-number">Line {{.LineNumber}}</span>
                        </div>
                        <code class="line-text">{{.LineText}}</code>
                    </div>
                {{end}}
            </div>
        {{else if and (not .Error) (ne .Query "")}}
            <div class="no-results">
                <div class="no-results-icon">🔍</div>
                <p>No matches found. Try adjusting your search query or enable fuzzy search.</p>
            </div>
        {{else if and (not .Error) (eq .Query "")}}
            <div class="empty-state">
                <div class="empty-state-icon">⚡</div>
                <p class="empty-state-text">Enter a search query to find code in your projects</p>
                <ul style="text-align: left; max-width: 400px; margin: 0 auto; color: var(--gray-primary);">
                    <li>Search for functions, classes, or specific code patterns</li>
                    <li>Enable fuzzy search for more flexible matching</li>
                    <li>Results show file location and exact line numbers</li>
                </ul>
            </div>
        {{end}}
    </div>

    <script>
        // Add keyboard shortcuts
        document.addEventListener('keydown', function(e) {
            // "/" to focus path input
            if (e.key === '/' && document.activeElement.tagName !== 'INPUT') {
                e.preventDefault();
                document.getElementById('path').focus();
            }
            
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
