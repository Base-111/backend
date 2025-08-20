package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"github.com/Base-111/backend/internal/entities/admin/domain"
	appErrors "github.com/Base-111/backend/pkg/errors"
	"github.com/Base-111/backend/pkg/errors/sql"
	"github.com/pkg/errors"
)

const promptTable = "admin.prompts"

const promptNameField = "name"
const promptIDField = "id"
const promptTextField = "text"
const promptSystemField = "is_system"

type PromptRepository struct {
	pool *pgxpool.Pool
	qb   sq.StatementBuilderType
}

func NewPromptRepository(pool *pgxpool.Pool, qb sq.StatementBuilderType) *PromptRepository {
	return &PromptRepository{
		pool: pool,
		qb:   qb,
	}
}

func (p *PromptRepository) Insert(ctx context.Context, prompt domain.Prompt) error {
	rawQuery := p.qb.Insert(promptTable).
		Columns(promptNameField, promptTextField, promptSystemField).
		Values(prompt.Name, prompt.Text, prompt.System)

	query, args, err := rawQuery.ToSql()
	if err != nil {
		return sql.WrapSQLError("build sql query", query, args, err)
	}

	_, err = p.pool.Exec(ctx, query, args...)

	if err != nil {
		return sql.WrapSQLError("insert prompt", query, args, err)
	}

	return nil
}

func (p *PromptRepository) GetById(ctx context.Context, id int) (domain.Prompt, error) {
	query, args, err := p.qb.Select(
		promptIDField,
		promptNameField,
		promptTextField,
		promptSystemField,
	).
		From(promptTable).
		Where(sq.Eq{promptIDField: id}).
		ToSql()

	if err != nil {
		return domain.Prompt{}, sql.WrapSQLError("build sql query", query, args, err)
	}

	row, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return domain.Prompt{}, sql.WrapSQLError("get prompt", query, args, err)
	}

	product, err := pgx.CollectOneRow(row, func(row pgx.CollectableRow) (domain.Prompt, error) {
		var product domain.Prompt
		err = row.Scan(
			&product.Id,
			&product.Name,
			&product.Text,
			&product.System,
		)
		if err != nil {
			return domain.Prompt{}, err
		}

		return product, nil
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return domain.Prompt{}, appErrors.NewNotFoundError(errors.New("prompt not found"))
		}

		return domain.Prompt{}, sql.WrapSQLError("get prompt", query, args, err)
	}

	return product, nil
}

func (p *PromptRepository) GetAll(ctx context.Context, params domain.PromptFilterParams) ([]domain.Prompt, error) {
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)
	query, args, err := p.qb.Select(
		promptIDField,
		promptNameField,
		promptTextField,
		promptSystemField,
	).
		From(promptTable).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return []domain.Prompt{}, sql.WrapSQLError("build sql query", query, args, err)
	}

	row, err := p.pool.Query(ctx, query, args...)

	if err != nil {
		return []domain.Prompt{}, sql.WrapSQLError("collect prompts", query, args, err)
	}

	products, err := pgx.CollectRows(row, func(row pgx.CollectableRow) (domain.Prompt, error) {
		var product domain.Prompt
		err = row.Scan(
			&product.Id,
			&product.Name,
			&product.Text,
			&product.System,
		)
		if err != nil {
			return domain.Prompt{}, err
		}

		return product, nil
	})

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Prompt{}, appErrors.NewNotFoundError(errors.New("prompt not found"))
		}

		return nil, sql.WrapSQLError("collect prompt", query, args, err)
	}

	return products, nil
}

func (p *PromptRepository) Update(ctx context.Context, id int, product domain.Prompt) error {
	updateQuery := p.qb.Update(promptTable).
		SetMap(map[string]any{
			promptNameField:   product.Name,
			promptTextField:   product.Text,
			promptSystemField: product.System,
		}).
		Where(sq.Eq{"id": id})

	sqlQuery, arguments, err := updateQuery.ToSql()
	if err != nil {
		return sql.WrapSQLError("build sql query", sqlQuery, arguments, err)
	}

	_, err = p.pool.Exec(ctx, sqlQuery, arguments...)
	if err != nil {
		return sql.WrapSQLError("update prompt", sqlQuery, arguments, err)
	}

	return nil
}

func (p *PromptRepository) Delete(ctx context.Context, id int) error {
	rawQuery := p.qb.Delete(promptTable).
		Where(sq.Eq{promptIDField: id})

	query, args, err := rawQuery.ToSql()
	if err != nil {
		return sql.WrapSQLError("build sql query", query, args, err)
	}

	_, err = p.pool.Exec(ctx, query, args...)

	if err != nil {
		return sql.WrapSQLError("delete prompt", query, args, err)
	}

	return nil
}
