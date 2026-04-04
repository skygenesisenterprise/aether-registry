package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/skygenesisenterprise/aether-bank/server/src/config"
	"github.com/skygenesisenterprise/aether-bank/server/src/models"
	"golang.org/x/crypto/bcrypt"
)

type PrismaService struct {
	DB *sql.DB
}

var prismaInstance *PrismaService

func NewPrismaService(cfg *config.Config) (*PrismaService, error) {
	if prismaInstance != nil {
		return prismaInstance, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	prismaInstance = &PrismaService{
		DB: db,
	}

	if err := prismaInstance.initTables(); err != nil {
		fmt.Printf("\033[1;33m[!] Warning: Failed to initialize tables: %v\033[0m\n", err)
	}

	if err := prismaInstance.initBankingTables(); err != nil {
		fmt.Printf("\033[1;33m[!] Warning: Failed to initialize banking tables: %v\033[0m\n", err)
	}

	return prismaInstance, nil
}

func (p *PrismaService) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		phone VARCHAR(50),
		avatar_url TEXT,
		role VARCHAR(50) DEFAULT 'USER',
		is_active BOOLEAN DEFAULT true,
		email_verified TIMESTAMP,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`
	_, err := p.DB.Exec(schema)
	return err
}

func GetPrismaService() *PrismaService {
	return prismaInstance
}

func (p *PrismaService) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

func (p *PrismaService) ListArticles(status, category, search string, page, pageSize int) ([]models.Article, int, error) {
	ctx := context.Background()

	query := "SELECT id, title, slug, COALESCE(excerpt, ''), content, content_html, status, featured, published_at, scheduled_at, view_count, read_time, COALESCE(image_url, ''), COALESCE(image_alt, ''), COALESCE(seo_title, ''), COALESCE(seo_description, ''), COALESCE(seo_keywords, ''), COALESCE(locale, 'fr'), author_id, COALESCE(category_id, ''), created_at, updated_at FROM articles WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM articles WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if category != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR excerpt ILIKE $%d)", argIndex, argIndex)
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR excerpt ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		var a models.Article
		err := rows.Scan(
			&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Content, &a.ContentHtml, &a.Status,
			&a.Featured, &a.PublishedAt, &a.ScheduledAt, &a.ViewCount, &a.ReadTime,
			&a.ImageUrl, &a.ImageAlt, &a.SeoTitle, &a.SeoDescription, &a.SeoKeywords,
			&a.Locale, &a.AuthorID, &a.CategoryID, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, a)
	}

	return articles, total, nil
}

func (p *PrismaService) GetArticle(id string) (*models.Article, error) {
	ctx := context.Background()

	query := "SELECT id, title, slug, COALESCE(excerpt, ''), content, content_html, status, featured, published_at, scheduled_at, view_count, read_time, COALESCE(image_url, ''), COALESCE(image_alt, ''), COALESCE(seo_title, ''), COALESCE(seo_description, ''), COALESCE(seo_keywords, ''), COALESCE(locale, 'fr'), author_id, COALESCE(category_id, ''), created_at, updated_at FROM articles WHERE id = $1"

	var a models.Article
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Content, &a.ContentHtml, &a.Status,
		&a.Featured, &a.PublishedAt, &a.ScheduledAt, &a.ViewCount, &a.ReadTime,
		&a.ImageUrl, &a.ImageAlt, &a.SeoTitle, &a.SeoDescription, &a.SeoKeywords,
		&a.Locale, &a.AuthorID, &a.CategoryID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return &a, nil
}

func (p *PrismaService) CreateArticle(req *models.CreateArticleRequest, authorID string) (*models.Article, error) {
	ctx := context.Background()

	id := fmt.Sprintf("art_%d", time.Now().UnixNano())
	slug := generateSlug(req.Title)
	now := time.Now()
	status := models.ArticleStatusDraft

	query := `INSERT INTO articles (id, title, slug, excerpt, content, status, featured, view_count, read_time, image_url, image_alt, seo_title, seo_description, seo_keywords, locale, author_id, category_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
			  RETURNING id, title, slug, COALESCE(excerpt, ''), content, content_html, status, featured, published_at, scheduled_at, view_count, read_time, COALESCE(image_url, ''), COALESCE(image_alt, ''), COALESCE(seo_title, ''), COALESCE(seo_description, ''), COALESCE(seo_keywords, ''), COALESCE(locale, 'fr'), author_id, COALESCE(category_id, ''), created_at, updated_at`

	var a models.Article
	err := p.DB.QueryRowContext(ctx, query,
		id, req.Title, slug, req.Excerpt, req.Content, status, false, 0, 5,
		req.ImageUrl, req.ImageAlt, req.SeoTitle, req.SeoDescription, req.SeoKeywords,
		"fr", authorID, req.CategoryID, now, now,
	).Scan(
		&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Content, &a.ContentHtml, &a.Status,
		&a.Featured, &a.PublishedAt, &a.ScheduledAt, &a.ViewCount, &a.ReadTime,
		&a.ImageUrl, &a.ImageAlt, &a.SeoTitle, &a.SeoDescription, &a.SeoKeywords,
		&a.Locale, &a.AuthorID, &a.CategoryID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return &a, nil
}

func (p *PrismaService) UpdateArticle(id string, req *models.UpdateArticleRequest) (*models.Article, error) {
	ctx := context.Background()

	query := `UPDATE articles SET title = COALESCE(NULLIF($2, ''), title), excerpt = COALESCE(NULLIF($3, ''), excerpt), 
			  content = COALESCE(NULLIF($4, ''), content), category_id = COALESCE(NULLIF($5, ''), category_id),
			  image_url = COALESCE(NULLIF($6, ''), image_url), updated_at = $7
			  WHERE id = $1
			  RETURNING id, title, slug, COALESCE(excerpt, ''), content, content_html, status, featured, published_at, scheduled_at, view_count, read_time, COALESCE(image_url, ''), COALESCE(image_alt, ''), COALESCE(seo_title, ''), COALESCE(seo_description, ''), COALESCE(seo_keywords, ''), COALESCE(locale, 'fr'), author_id, COALESCE(category_id, ''), created_at, updated_at`

	var a models.Article
	err := p.DB.QueryRowContext(ctx, query,
		id, req.Title, req.Excerpt, req.Content, req.CategoryID, req.ImageUrl, time.Now(),
	).Scan(
		&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Content, &a.ContentHtml, &a.Status,
		&a.Featured, &a.PublishedAt, &a.ScheduledAt, &a.ViewCount, &a.ReadTime,
		&a.ImageUrl, &a.ImageAlt, &a.SeoTitle, &a.SeoDescription, &a.SeoKeywords,
		&a.Locale, &a.AuthorID, &a.CategoryID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return &a, nil
}

func (p *PrismaService) DeleteArticle(id string) error {
	ctx := context.Background()

	_, err := p.DB.ExecContext(ctx, "DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	return nil
}

func (p *PrismaService) PublishArticle(id string) error {
	ctx := context.Background()

	now := time.Now()
	_, err := p.DB.ExecContext(ctx, "UPDATE articles SET status = $1, published_at = $2, updated_at = $2 WHERE id = $3",
		models.ArticleStatusPublished, now, id)
	if err != nil {
		return fmt.Errorf("failed to publish article: %w", err)
	}
	return nil
}

func (p *PrismaService) ArchiveArticle(id string) error {
	ctx := context.Background()

	_, err := p.DB.ExecContext(ctx, "UPDATE articles SET status = $1, updated_at = $2 WHERE id = $3",
		models.ArticleStatusArchived, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to archive article: %w", err)
	}
	return nil
}

func generateSlug(title string) string {
	slug := ""
	for i := 0; i < len(title); i++ {
		c := title[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			slug += string(c)
		} else if c == ' ' || c == '-' {
			slug += "-"
		}
	}
	if slug == "" {
		slug = "article"
	}
	return slug
}

func (p *PrismaService) ListUsers(search string, page, pageSize int) ([]models.EtheriaUser, int, error) {
	ctx := context.Background()

	query := "SELECT id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(avatar_url, ''), role, is_active, email_verified, created_at, updated_at FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argIndex, argIndex, argIndex)
		countQuery += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []models.EtheriaUser{}
	for rows.Next() {
		var u models.EtheriaUser
		err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarUrl, &u.Role,
			&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (p *PrismaService) GetUser(id string) (*models.EtheriaUser, error) {
	ctx := context.Background()

	query := "SELECT id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(avatar_url, ''), COALESCE(password_hash, ''), role, is_active, email_verified, created_at, updated_at FROM users WHERE id = $1"

	var u models.EtheriaUser
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarUrl, &u.Password,
		&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

func (p *PrismaService) GetUserByEmail(email string) (*models.EtheriaUser, error) {
	ctx := context.Background()

	query := "SELECT id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(avatar_url, ''), COALESCE(password_hash, ''), role, is_active, email_verified, created_at, updated_at FROM users WHERE email = $1"

	var u models.EtheriaUser
	err := p.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarUrl, &u.Password,
		&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}

func (p *PrismaService) CreateUser(email, firstName, lastName, role, password string) (*models.EtheriaUser, error) {
	ctx := context.Background()

	id := fmt.Sprintf("user_%d", time.Now().UnixNano())
	now := time.Now()
	defaultRole := models.RoleUser
	switch role {
	case "ADMIN":
		defaultRole = models.RoleAdmin
	case "EDITOR":
		defaultRole = models.RoleEditor
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO users (id, email, first_name, last_name, role, password_hash, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(avatar_url, ''), role, is_active, email_verified, created_at, updated_at`

	var u models.EtheriaUser
	err = p.DB.QueryRowContext(ctx, query,
		id, email, firstName, lastName, defaultRole, string(hashedPassword), true, now, now,
	).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarUrl, &u.Role,
		&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &u, nil
}

func (p *PrismaService) UpdateUser(id, firstName, lastName, role string, isActive bool) (*models.EtheriaUser, error) {
	ctx := context.Background()

	query := `UPDATE users SET first_name = COALESCE(NULLIF($2, ''), first_name), 
			  last_name = COALESCE(NULLIF($3, ''), last_name), role = COALESCE(NULLIF($4, ''), role),
			  is_active = $5, updated_at = $6
			  WHERE id = $1
			  RETURNING id, email, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(avatar_url, ''), role, is_active, email_verified, created_at, updated_at`

	var u models.EtheriaUser
	err := p.DB.QueryRowContext(ctx, query, id, firstName, lastName, role, isActive, time.Now()).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarUrl, &u.Role,
		&u.IsActive, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &u, nil
}

func (p *PrismaService) DeleteUser(id string) error {
	ctx := context.Background()

	_, err := p.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (p *PrismaService) AdminExists() (bool, error) {
	ctx := context.Background()

	var count int
	err := p.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE role = 'ADMIN' AND is_active = true").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check admin existence: %w", err)
	}
	return count > 0, nil
}

// ==================== BANKING METHODS ====================

func (p *PrismaService) initBankingTables() error {
	schema := `
	-- Bank Accounts Table
	CREATE TABLE IF NOT EXISTS bank_accounts (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255),
		account_type VARCHAR(50) DEFAULT 'CURRENT',
		status VARCHAR(50) DEFAULT 'PENDING_KYC',
		currency VARCHAR(10) DEFAULT 'EUR',
		iban VARCHAR(50) UNIQUE NOT NULL,
		bic VARCHAR(20) NOT NULL,
		holder_name VARCHAR(255) NOT NULL,
		holder_type VARCHAR(50) DEFAULT 'individual',
		balance BIGINT DEFAULT 0,
		available BIGINT DEFAULT 0,
		pending BIGINT DEFAULT 0,
		overdraft BIGINT DEFAULT 0,
		provider_ref VARCHAR(255),
		metadata JSONB DEFAULT '{}',
		client_id VARCHAR(255),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Banking Cards Table
	CREATE TABLE IF NOT EXISTS banking_cards (
		id VARCHAR(255) PRIMARY KEY,
		account_id VARCHAR(255) NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
		card_type VARCHAR(50) DEFAULT 'VIRTUAL',
		status VARCHAR(50) DEFAULT 'PENDING',
		holder_name VARCHAR(255) NOT NULL,
		pan VARCHAR(255),
		last4 VARCHAR(4) NOT NULL,
		expiry_month INT NOT NULL,
		expiry_year INT NOT NULL,
		cvv VARCHAR(10),
		brand VARCHAR(50) DEFAULT 'VISA',
		daily_limit BIGINT DEFAULT 100000,
		monthly_limit BIGINT DEFAULT 300000,
		provider_ref VARCHAR(255),
		frozen_at TIMESTAMP,
		frozen_reason VARCHAR(500),
		estimated_delivery TIMESTAMP,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Transactions Table
	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(255) PRIMARY KEY,
		account_id VARCHAR(255) NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
		type VARCHAR(50) NOT NULL,
		status VARCHAR(50) DEFAULT 'PENDING',
		amount BIGINT NOT NULL,
		currency VARCHAR(10) DEFAULT 'EUR',
		description TEXT,
		reference VARCHAR(255),
		balance_after BIGINT,
		counterparty_name VARCHAR(255),
		counterparty_iban VARCHAR(50),
		merchant_name VARCHAR(255),
		merchant_category VARCHAR(100),
		reconciled BOOLEAN DEFAULT FALSE,
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMP DEFAULT NOW(),
		completed_at TIMESTAMP
	);

	-- Transfers Table
	CREATE TABLE IF NOT EXISTS transfers (
		id VARCHAR(255) PRIMARY KEY,
		account_id VARCHAR(255) NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
		type VARCHAR(50) DEFAULT 'SEPA',
		direction VARCHAR(50) NOT NULL,
		amount BIGINT NOT NULL,
		currency VARCHAR(10) DEFAULT 'EUR',
		status VARCHAR(50) DEFAULT 'PENDING',
		description TEXT,
		reference VARCHAR(255),
		iban VARCHAR(50),
		bic VARCHAR(20),
		counterparty_name VARCHAR(255),
		counterparty_iban VARCHAR(50),
		external_ref VARCHAR(255),
		provider_ref VARCHAR(255),
		fees BIGINT DEFAULT 0,
		idempotency_key VARCHAR(255),
		metadata JSONB DEFAULT '{}',
		executed_at TIMESTAMP,
		completed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Clients Table
	CREATE TABLE IF NOT EXISTS clients (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		company VARCHAR(255),
		status VARCHAR(50) DEFAULT 'PROSPECT',
		kyc_verified BOOLEAN DEFAULT FALSE,
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_bank_accounts_user_id ON bank_accounts(user_id);
	CREATE INDEX IF NOT EXISTS idx_bank_accounts_status ON bank_accounts(status);
	CREATE INDEX IF NOT EXISTS idx_bank_accounts_iban ON bank_accounts(iban);
	CREATE INDEX IF NOT EXISTS idx_banking_cards_account_id ON banking_cards(account_id);
	CREATE INDEX IF NOT EXISTS idx_banking_cards_status ON banking_cards(status);
	CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
	CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
	CREATE INDEX IF NOT EXISTS idx_transfers_account_id ON transfers(account_id);
	CREATE INDEX IF NOT EXISTS idx_transfers_status ON transfers(status);
	`
	_, err := p.DB.Exec(schema)
	return err
}

// Bank Account Methods

func (p *PrismaService) CreateBankAccount(acc *models.BankingAccount) error {
	ctx := context.Background()

	query := `INSERT INTO bank_accounts (id, user_id, account_type, status, currency, iban, bic, holder_name, holder_type, balance, available, pending, overdraft, metadata, client_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	metadata := "{}"
	if acc.Metadata != nil {
		// Convert map to JSON string
	}

	_, err := p.DB.ExecContext(ctx, query,
		acc.ID, "", acc.Type, acc.Status, acc.Currency, acc.IBAN, acc.BIC,
		acc.Holder.Name, acc.Holder.Type, acc.Balance.Available, acc.Balance.Available,
		acc.Balance.Pending, acc.Balance.Overdraft, metadata, "",
		acc.CreatedAt, acc.UpdatedAt,
	)
	return err
}

func (p *PrismaService) GetBankAccount(id string) (*models.BankingAccount, error) {
	ctx := context.Background()

	query := `SELECT id, COALESCE(user_id, ''), account_type, status, currency, iban, bic, holder_name, holder_type, balance, available, pending, overdraft, COALESCE(metadata, '{}'), COALESCE(client_id, ''), created_at, updated_at FROM bank_accounts WHERE id = $1`

	var acc models.BankingAccount
	var holderName, holderType string
	var balance, available, pending, overdraft int64

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&acc.ID, &acc.UserID, &acc.Type, &acc.Status, &acc.Currency,
		&acc.IBAN, &acc.BIC, &holderName, &holderType,
		&balance, &available, &pending, &overdraft,
		&acc.Metadata, &acc.ClientID, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	acc.Holder = &models.AccountHolder{Name: holderName, Type: models.HolderType(holderType)}
	acc.Balance = &models.AccountBalance{
		Current:   balance,
		Available: available,
		Pending:   pending,
		Overdraft: overdraft,
	}

	return &acc, nil
}

func (p *PrismaService) GetBankAccountByIBAN(iban string) (*models.BankingAccount, error) {
	ctx := context.Background()

	query := `SELECT id, COALESCE(user_id, ''), account_type, status, currency, iban, bic, holder_name, holder_type, balance, available, pending, overdraft, COALESCE(metadata, '{}'), COALESCE(client_id, ''), created_at, updated_at FROM bank_accounts WHERE iban = $1`

	var acc models.BankingAccount
	var holderName, holderType string
	var balance, available, pending, overdraft int64

	err := p.DB.QueryRowContext(ctx, query, iban).Scan(
		&acc.ID, &acc.UserID, &acc.Type, &acc.Status, &acc.Currency,
		&acc.IBAN, &acc.BIC, &holderName, &holderType,
		&balance, &available, &pending, &overdraft,
		&acc.Metadata, &acc.ClientID, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	acc.Holder = &models.AccountHolder{Name: holderName, Type: models.HolderType(holderType)}
	acc.Balance = &models.AccountBalance{
		Current:   balance,
		Available: available,
		Pending:   pending,
		Overdraft: overdraft,
	}

	return &acc, nil
}

func (p *PrismaService) ListBankAccounts(status, accountType string, limit, offset int) ([]models.BankingAccount, int, error) {
	ctx := context.Background()

	query := `SELECT id, COALESCE(user_id, ''), account_type, status, currency, iban, bic, holder_name, holder_type, balance, available, pending, overdraft, COALESCE(metadata, '{}'), COALESCE(client_id, ''), created_at, updated_at FROM bank_accounts WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM bank_accounts WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, strings.ToUpper(status))
		argIndex++
	}

	if accountType != "" {
		query += fmt.Sprintf(" AND account_type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND account_type = $%d", argIndex)
		args = append(args, strings.ToUpper(accountType))
		argIndex++
	}

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	accounts := []models.BankingAccount{}
	for rows.Next() {
		var acc models.BankingAccount
		var holderName, holderType string
		var balance, available, pending, overdraft int64

		err := rows.Scan(
			&acc.ID, &acc.UserID, &acc.Type, &acc.Status, &acc.Currency,
			&acc.IBAN, &acc.BIC, &holderName, &holderType,
			&balance, &available, &pending, &overdraft,
			&acc.Metadata, &acc.ClientID, &acc.CreatedAt, &acc.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		acc.Holder = &models.AccountHolder{Name: holderName, Type: models.HolderType(holderType)}
		acc.Balance = &models.AccountBalance{
			Current:   balance,
			Available: available,
			Pending:   pending,
			Overdraft: overdraft,
		}
		accounts = append(accounts, acc)
	}

	return accounts, total, nil
}

func (p *PrismaService) UpdateBankAccountBalance(id string, balance, available, pending int64) error {
	ctx := context.Background()

	query := `UPDATE bank_accounts SET balance = $2, available = $3, pending = $4, updated_at = NOW() WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id, balance, available, pending)
	return err
}

func (p *PrismaService) UpdateBankAccountStatus(id, status string) error {
	ctx := context.Background()

	query := `UPDATE bank_accounts SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id, strings.ToUpper(status))
	return err
}

// Card Methods

func (p *PrismaService) CreateBankingCard(card *models.BankingCard) error {
	ctx := context.Background()

	query := `INSERT INTO banking_cards (id, account_id, card_type, status, holder_name, pan, last4, expiry_month, expiry_year, cvv, brand, daily_limit, monthly_limit, provider_ref, frozen_at, frozen_reason, estimated_delivery, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err := p.DB.ExecContext(ctx, query,
		card.ID, card.AccountID, card.Type, card.Status, card.HolderName,
		card.PAN, card.Last4, card.ExpiryMonth, card.ExpiryYear, card.CVV,
		card.Brand, 100000, 300000, "",
		card.FrozenAt, card.FrozenReason, card.EstimatedDelivery,
		card.CreatedAt, card.CreatedAt,
	)
	return err
}

func (p *PrismaService) GetBankingCard(id string) (*models.BankingCard, error) {
	ctx := context.Background()

	query := `SELECT id, account_id, card_type, status, holder_name, COALESCE(pan, ''), last4, expiry_month, expiry_year, COALESCE(cvv, ''), brand, daily_limit, monthly_limit, COALESCE(provider_ref, ''), frozen_at, COALESCE(frozen_reason, ''), estimated_delivery, created_at, updated_at FROM banking_cards WHERE id = $1`

	var card models.BankingCard
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&card.ID, &card.AccountID, &card.Type, &card.Status, &card.HolderName,
		&card.PAN, &card.Last4, &card.ExpiryMonth, &card.ExpiryYear, &card.CVV,
		&card.Brand, &card.SpendingLimits.Daily, &card.SpendingLimits.Monthly,
		&card.ProviderRef, &card.FrozenAt, &card.FrozenReason, &card.EstimatedDelivery,
		&card.CreatedAt, &card.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (p *PrismaService) ListBankingCards(accountID, status, cardType string, limit, offset int) ([]models.BankingCard, int, error) {
	ctx := context.Background()

	query := `SELECT id, account_id, card_type, status, holder_name, COALESCE(pan, ''), last4, expiry_month, expiry_year, COALESCE(cvv, ''), brand, daily_limit, monthly_limit, COALESCE(provider_ref, ''), frozen_at, COALESCE(frozen_reason, ''), estimated_delivery, created_at, updated_at FROM banking_cards WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM banking_cards WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if accountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND account_id = $%d", argIndex)
		args = append(args, accountID)
		argIndex++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, strings.ToUpper(status))
		argIndex++
	}

	if cardType != "" {
		query += fmt.Sprintf(" AND card_type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND card_type = $%d", argIndex)
		args = append(args, strings.ToUpper(cardType))
		argIndex++
	}

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	cards := []models.BankingCard{}
	for rows.Next() {
		var card models.BankingCard
		err := rows.Scan(
			&card.ID, &card.AccountID, &card.Type, &card.Status, &card.HolderName,
			&card.PAN, &card.Last4, &card.ExpiryMonth, &card.ExpiryYear, &card.CVV,
			&card.Brand, &card.SpendingLimits.Daily, &card.SpendingLimits.Monthly,
			&card.ProviderRef, &card.FrozenAt, &card.FrozenReason, &card.EstimatedDelivery,
			&card.CreatedAt, &card.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		cards = append(cards, card)
	}

	return cards, total, nil
}

func (p *PrismaService) UpdateBankingCardStatus(id, status string) error {
	ctx := context.Background()

	query := `UPDATE banking_cards SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id, strings.ToUpper(status))
	return err
}

func (p *PrismaService) FreezeBankingCard(id, reason string) error {
	ctx := context.Background()

	query := `UPDATE banking_cards SET status = 'FROZEN', frozen_at = NOW(), frozen_reason = $2, updated_at = NOW() WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id, reason)
	return err
}

func (p *PrismaService) UnfreezeBankingCard(id string) error {
	ctx := context.Background()

	query := `UPDATE banking_cards SET status = 'ACTIVE', frozen_at = NULL, frozen_reason = '', updated_at = NOW() WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}

// Transaction Methods

func (p *PrismaService) CreateTransaction(txn *models.Transaction) error {
	ctx := context.Background()

	query := `INSERT INTO transactions (id, account_id, type, status, amount, currency, description, reference, balance_after, counterparty_name, counterparty_iban, merchant_name, merchant_category, reconciled, metadata, created_at, completed_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	metadata := "{}"
	_, err := p.DB.ExecContext(ctx, query,
		txn.ID, txn.AccountID, txn.Type, txn.Status, txn.Amount, txn.Currency,
		txn.Description, txn.Reference, txn.BalanceAfter, txn.CounterpartyName,
		txn.CounterpartyIban, txn.MerchantName, txn.MerchantCategory,
		txn.Reconciled, metadata, txn.CreatedAt, txn.CompletedAt,
	)
	return err
}

func (p *PrismaService) ListTransactions(accountID, status, txnType string, limit, offset int) ([]models.Transaction, int, error) {
	ctx := context.Background()

	query := `SELECT id, account_id, type, status, amount, currency, COALESCE(description, ''), COALESCE(reference, ''), balance_after, COALESCE(counterparty_name, ''), COALESCE(counterparty_iban, ''), COALESCE(merchant_name, ''), COALESCE(merchant_category, ''), reconciled, COALESCE(metadata, '{}'), created_at, completed_at FROM transactions WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM transactions WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if accountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND account_id = $%d", argIndex)
		args = append(args, accountID)
		argIndex++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, strings.ToUpper(status))
		argIndex++
	}

	if txnType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, strings.ToUpper(txnType))
		argIndex++
	}

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	for rows.Next() {
		var txn models.Transaction
		err := rows.Scan(
			&txn.ID, &txn.AccountID, &txn.Type, &txn.Status, &txn.Amount, &txn.Currency,
			&txn.Description, &txn.Reference, &txn.BalanceAfter, &txn.CounterpartyName,
			&txn.CounterpartyIban, &txn.MerchantName, &txn.MerchantCategory,
			&txn.Reconciled, &txn.Metadata, &txn.CreatedAt, &txn.CompletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, txn)
	}

	return transactions, total, nil
}

// Client Methods

func (p *PrismaService) CreateClient(client *models.Client) error {
	ctx := context.Background()

	query := `INSERT INTO clients (id, email, name, company, status, kyc_verified, metadata, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	metadata := "{}"
	_, err := p.DB.ExecContext(ctx, query,
		client.ID, client.Email, client.Name, client.Company, client.Status,
		client.KYCVerified, metadata, client.CreatedAt, client.UpdatedAt,
	)
	return err
}

func (p *PrismaService) ListClients(limit, offset int) ([]models.Client, int, error) {
	ctx := context.Background()

	query := `SELECT id, email, name, COALESCE(company, ''), status, kyc_verified, COALESCE(metadata, '{}'), created_at, updated_at FROM clients`
	countQuery := `SELECT COUNT(*) FROM clients`

	var total int
	err := p.DB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	clients := []models.Client{}
	for rows.Next() {
		var client models.Client
		err := rows.Scan(
			&client.ID, &client.Email, &client.Name, &client.Company,
			&client.Status, &client.KYCVerified, &client.Metadata,
			&client.CreatedAt, &client.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, client)
	}

	return clients, total, nil
}
