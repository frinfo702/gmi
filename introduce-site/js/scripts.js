// DOM要素の取得
const demoSearch = document.getElementById('demo-search');
const demoResults = document.getElementById('demo-results');
const waitlistForm = document.getElementById('waitlist-form');
const waitlistSuccess = document.getElementById('waitlist-success');

// デモ検索機能の実装
demoSearch.addEventListener('keyup', handleDemoSearch);

function handleDemoSearch(e) {
  if (e.key === 'Enter') {
    const query = e.target.value.trim();
    if (query) {
      displaySearchResults(query);
    }
  }
}

function displaySearchResults(query) {
  const resultsContainer = document.getElementById('demo-results');
  query = query.toLowerCase();
  
  // Clear previous results
  resultsContainer.innerHTML = '';
  
  // Mock search results based on context understanding rather than fuzzy matching
  let results = [];
  
  // Generate relevant results based on the query content
  if (query.includes('authentication') || query.includes('login') || query.includes('auth')) {
    results = [
      {
        filePath: 'src/auth/service.go',
        lineNumber: 42,
        score: 0.92,
        lineText: 'func (s *AuthService) Authenticate(credentials *model.Credentials) (*model.User, error) {',
        context: 'Core authentication logic that validates user credentials against the database'
      },
      {
        filePath: 'src/auth/controller.go',
        lineNumber: 28,
        score: 0.87,
        lineText: 'func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {',
        context: 'REST API endpoint handler for user login requests'
      },
      {
        filePath: 'src/middleware/jwt.go',
        lineNumber: 15,
        score: 0.83,
        lineText: 'func GenerateToken(user *model.User) (string, error) {',
        context: 'JWT token generation for authenticated users'
      }
    ];
  } else if (query.includes('database') || query.includes('db') || query.includes('connection')) {
    results = [
      {
        filePath: 'src/infrastructure/database.go',
        lineNumber: 23,
        score: 0.94,
        lineText: 'func NewDatabaseConnection(config *Config) (*Database, error) {',
        context: 'Establishes database connections with configuration parameters'
      },
      {
        filePath: 'src/config/database_config.go',
        lineNumber: 12,
        score: 0.85,
        lineText: 'type DatabaseConfig struct {',
        context: 'Configuration structure for database connection settings'
      },
      {
        filePath: 'src/models/repository.go',
        lineNumber: 31,
        score: 0.78,
        lineText: 'func (r *Repository) Connect() error {',
        context: 'Generic repository pattern implementation for database access'
      }
    ];
  } else if (query.includes('api') || query.includes('endpoint') || query.includes('rest')) {
    results = [
      {
        filePath: 'src/api/router.go',
        lineNumber: 45,
        score: 0.91,
        lineText: 'func RegisterRoutes(e *echo.Echo, controllers *Controllers) {',
        context: 'Main router configuration that registers all API endpoints'
      },
      {
        filePath: 'src/api/middleware.go',
        lineNumber: 28,
        score: 0.84,
        lineText: 'func AuthMiddleware() echo.MiddlewareFunc {',
        context: 'Authentication middleware for protecting API endpoints'
      },
      {
        filePath: 'src/controllers/user_controller.go',
        lineNumber: 57,
        score: 0.79,
        lineText: 'func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {',
        context: 'API handler for retrieving user profile information'
      }
    ];
  } else {
    // Default results for any other query
    results = [
      {
        filePath: 'src/main.go',
        lineNumber: 28,
        score: 0.75,
        lineText: 'func main() {',
        context: 'Application entry point that initializes all components'
      },
      {
        filePath: 'src/utils/logger.go',
        lineNumber: 42,
        score: 0.65,
        lineText: 'func NewLogger(config *LogConfig) *Logger {',
        context: 'Configurable logging utility for application-wide use'
      },
      {
        filePath: 'src/models/user.go',
        lineNumber: 15,
        score: 0.60,
        lineText: 'type User struct {',
        context: 'Core user data model with all user properties'
      }
    ];
  }
  
  // Sort by score (highest first)
  results.sort((a, b) => b.score - a.score);
  
  // Display results in a nice format
  if (results.length > 0) {
    const resultsList = document.createElement('ul');
    resultsList.className = 'search-result-list';
    
    results.forEach(result => {
      const li = document.createElement('li');
      li.className = 'search-result-item';
      
      // Calculate relevance percentage
      const relevancePercent = Math.round(result.score * 100);
      
      li.innerHTML = `
        <div class="result-header">
          <span class="result-file">${result.filePath}</span>
          <span class="result-line">Line ${result.lineNumber}</span>
          <span class="result-score">${relevancePercent}% match</span>
        </div>
        <div class="result-code">${result.lineText}</div>
        <div class="result-context">${result.context}</div>
      `;
      
      resultsList.appendChild(li);
    });
    
    resultsContainer.appendChild(resultsList);
  } else {
    resultsContainer.innerHTML = '<div class="placeholder-text">No results found</div>';
  }
}

// Mailchimpフォーム送信処理
waitlistForm.addEventListener('submit', handleWaitlistSubmit);

function handleWaitlistSubmit(e) {
  // フォームはMailchimpに直接送信されるため、preventDefault()は不要
  
  const submitButton = waitlistForm.querySelector('button[type="submit"]');
  const originalButtonText = submitButton.textContent;
  submitButton.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Submitting...';
  submitButton.disabled = true;

  // Google Analyticsでのイベント記録
  if (typeof gtag !== 'undefined') {
    gtag('event', 'waitlist_signup', {
      'event_category': 'engagement',
      'event_label': document.getElementById('interest').value
    });
  }

  // LocalStorageに登録済みとして記録
  setTimeout(() => {
    localStorage.setItem('rustysearch_waitlist', 'true');
    
    // フォームを非表示にして成功メッセージを表示
    waitlistForm.style.display = 'none';
    waitlistSuccess.style.display = 'block';
    
    // ボタンを元に戻す
    submitButton.innerHTML = originalButtonText;
    submitButton.disabled = false;
  }, 1000);
}

// ページ読み込み時に既に登録済みかチェック
document.addEventListener('DOMContentLoaded', () => {
  if (localStorage.getItem('rustysearch_waitlist') === 'true') {
    waitlistForm.style.display = 'none';
    waitlistSuccess.style.display = 'block';
  }
  
  // スムーススクロールの実装
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute('href'));
      if (target) {
        window.scrollTo({
          top: target.offsetTop - 100,
          behavior: 'smooth'
        });
      }
    });
  });
}); 
