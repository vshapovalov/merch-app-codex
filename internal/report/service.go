package report

import (
	"context"
	"database/sql"

	"merch-app-codex/internal/storage/mysql"
)

// Service exposes reporting capabilities on top of the storage layer.
type Service struct {
	repo *mysql.Repository
}

// NewService constructs a new reporting service.
func NewService(repo *mysql.Repository) *Service {
	return &Service{repo: repo}
}

// CompanyVisitSummary aggregates visits and sold items for a company.
type CompanyVisitSummary struct {
	CompanyID   string  `json:"company_id"`
	TotalVisits int64   `json:"total_visits"`
	TotalItems  int64   `json:"total_items"`
	TotalAmount float64 `json:"total_amount"`
}

// CompanyVisitSummary gathers visit statistics for the provided company ID.
func (s *Service) CompanyVisitSummary(ctx context.Context, companyID string) (CompanyVisitSummary, error) {
	summary := CompanyVisitSummary{CompanyID: companyID}

	if err := s.repo.DB().WithContext(ctx).
		Model(&mysql.Visit{}).
		Joins("JOIN retail_points ON retail_points.id = visits.retail_point_id").
		Where("retail_points.company_id = ?", companyID).
		Count(&summary.TotalVisits).Error; err != nil {
		return summary, err
	}

	var aggregate struct {
		TotalItems  sql.NullInt64
		TotalAmount sql.NullFloat64
	}

	if err := s.repo.DB().WithContext(ctx).
		Model(&mysql.VisitItem{}).
		Select("COALESCE(SUM(visit_items.present_quantity), 0) AS total_items, COALESCE(SUM(visit_items.present_quantity * visit_items.price), 0) AS total_amount").
		Joins("JOIN visits ON visits.id = visit_items.visit_id").
		Joins("JOIN retail_points ON retail_points.id = visits.retail_point_id").
		Where("retail_points.company_id = ?", companyID).
		Scan(&aggregate).Error; err != nil {
		return summary, err
	}

	summary.TotalItems = aggregate.TotalItems.Int64
	summary.TotalAmount = aggregate.TotalAmount.Float64

	return summary, nil
}
