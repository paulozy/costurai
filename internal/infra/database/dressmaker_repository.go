package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/paulozy/costurai/internal/entity"
)

type DressmakerRepository struct {
	DB *sql.DB
}

func NewDressmakerRepository(db *sql.DB) *DressmakerRepository {
	return &DressmakerRepository{
		DB: db,
	}
}

func (r *DressmakerRepository) Create(dressmaker *entity.Dressmaker) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO dressmakers (id, email, password, name, contact, location, services, grade, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, ST_MakePoint($6, $7)::geography, $8, $9, $10, $11)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		dressmaker.ID,
		dressmaker.Email,
		dressmaker.Password,
		dressmaker.Name,
		dressmaker.Contact,
		dressmaker.Location.Latitude,
		dressmaker.Location.Longitude,
		strings.Join(dressmaker.Services, ","),
		dressmaker.Grade,
		dressmaker.CreatedAt,
		dressmaker.UpdatedAt,
	)

	return err
}

func (r *DressmakerRepository) FindByEmail(email string) (*entity.Dressmaker, error) {
	var dressmaker entity.Dressmaker

	row := r.DB.QueryRow(`
		SELECT 
			id, 
			email,
			password, 
			name, 
			contact, 
			ST_Y(location::geometry) as latitude, 
			ST_X(location::geometry) as longitude, 
			services, 
			grade, 
			created_at, 
			updated_at
		FROM dressmakers
		WHERE email = $1
	`, email)

	var services string
	err := row.Scan(
		&dressmaker.ID,
		&dressmaker.Email,
		&dressmaker.Password,
		&dressmaker.Name,
		&dressmaker.Contact,
		&dressmaker.Location.Latitude,
		&dressmaker.Location.Longitude,
		&services,
		&dressmaker.Grade,
		&dressmaker.CreatedAt,
		&dressmaker.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	dressmaker.Services = strings.Split(services, ", ")

	return &dressmaker, nil
}

func (r *DressmakerRepository) Exists(email string) (bool, error) {
	var exists bool

	row := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM dressmakers
			WHERE email = $1
		)
	`, email)

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *DressmakerRepository) Find(searchParams GetDressmakersParams) ([]entity.Dressmaker, error) {
	var dressmakers []entity.Dressmaker

	var stmt *sql.Stmt
	var err error
	var rows *sql.Rows

	searchByServices := searchParams.Services != "" && searchParams.Latitude == 0 && searchParams.Longitude == 0 && searchParams.Distance == 0
	searchyByProximity := searchParams.Latitude != 0 && searchParams.Longitude != 0 && searchParams.Distance != 0 && searchParams.Services == ""

	if searchParams.Default {
		stmt, err = r.DB.Prepare(`
		SELECT 
				id, 
				email,
				password, 
				name, 
				contact, 
				ST_Y(location::geometry) as latitude, 
				ST_X(location::geometry) as longitude, 
				services, 
				grade, 
				created_at, 
				updated_at
			FROM dressmakers
		`)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}

		rows, err = stmt.Query()
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}
	} else if searchByServices {
		keywords := strings.Split(searchParams.Services, ", ")
		query := generateILikeQuery(keywords)

		stmt, err = r.DB.Prepare(query)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}

		args := make([]interface{}, len(keywords))
		for i, keyword := range keywords {
			args[i] = "%" + keyword + "%"
		}

		rows, err = stmt.Query(args...)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}
	} else if searchyByProximity {
		stmt, err = r.DB.Prepare(`
			SELECT 
				id, 
				email,
				password, 
				name, 
				contact, 
				ST_Y(location::geometry) as latitude, 
				ST_X(location::geometry) as longitude, 
				services, 
				grade, 
				created_at, 
				updated_at
			FROM dressmakers
			WHERE 
				ST_DWithin(
						location,
						ST_MakePoint($1, $2)::geography,
						$3
				);
		`)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}

		rows, err = stmt.Query(searchParams.Latitude, searchParams.Longitude, searchParams.Distance)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}
	} else {
		keywords := strings.Split(searchParams.Services, ", ")
		query := generateFullyILikeQuery(keywords)

		stmt, err = r.DB.Prepare(query)
		if err != nil {
			fmt.Printf("error on exec query: %v\n", err)
			return nil, err
		}

		args := make([]interface{}, len(keywords)+3)
		args[0] = searchParams.Latitude
		args[1] = searchParams.Longitude
		args[2] = searchParams.Distance

		for i, keyword := range keywords {
			args[i+3] = "%" + keyword + "%"
		}

		rows, err = stmt.Query(args...)
	}

	for rows.Next() {
		var dressmaker entity.Dressmaker
		var services string

		err := rows.Scan(
			&dressmaker.ID,
			&dressmaker.Email,
			&dressmaker.Password,
			&dressmaker.Name,
			&dressmaker.Contact,
			&dressmaker.Location.Latitude,
			&dressmaker.Location.Longitude,
			&services,
			&dressmaker.Grade,
			&dressmaker.CreatedAt,
			&dressmaker.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("error on scan rows: %v\n", err)
			return nil, err
		}

		dressmaker.Services = strings.Split(services, ", ")
		dressmakers = append(dressmakers, dressmaker)
	}

	return dressmakers, nil
}

func (r *DressmakerRepository) FindByID(id string) (*entity.Dressmaker, error) {
	var dressmaker entity.Dressmaker

	row := r.DB.QueryRow(`
		SELECT 
			id, 
			email,
			password, 
			name, 
			contact, 
			ST_Y(location::geometry) as latitude, 
			ST_X(location::geometry) as longitude, 
			services, 
			grade, 
			created_at, 
			updated_at
		FROM dressmakers
		WHERE id = $1
	`, id)

	var services string
	var reviews []entity.Review

	rows, err := r.DB.Query(`
		SELECT *
		FROM dressmakers_reviews
		WHERE dressmaker_id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var review entity.Review

		err := rows.Scan(
			&review.ID,
			&review.DressmakerID,
			&review.UserID,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Grade,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	err = row.Scan(
		&dressmaker.ID,
		&dressmaker.Email,
		&dressmaker.Password,
		&dressmaker.Name,
		&dressmaker.Contact,
		&dressmaker.Location.Latitude,
		&dressmaker.Location.Longitude,
		&services,
		&dressmaker.Grade,
		&dressmaker.CreatedAt,
		&dressmaker.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	dressmaker.Services = strings.Split(services, ", ")
	dressmaker.Reviews = reviews

	return &dressmaker, nil
}

func (r *DressmakerRepository) FindByProximity(latitude, longitude float64, maxDistance int) ([]entity.Dressmaker, error) {
	var dressmakers []entity.Dressmaker

	stmt, err := r.DB.Prepare(`
		SELECT 
			id, 
			email,
			password, 
			name, 
			contact, 
			ST_Y(location::geometry) as latitude, 
			ST_X(location::geometry) as longitude, 
			services, 
			grade, 
			created_at, 
			updated_at
		FROM dressmakers
		WHERE ST_DWithin(
				location,
				ST_MakePoint($1, $2)::geography,
				$3
		);
	`)
	if err != nil {
		fmt.Printf("error on prepare query: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(latitude, longitude, maxDistance)
	if err != nil {
		fmt.Printf("error on exec query: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var dressmaker entity.Dressmaker
		var services string

		err := rows.Scan(
			&dressmaker.ID,
			&dressmaker.Email,
			&dressmaker.Password,
			&dressmaker.Name,
			&dressmaker.Contact,
			&dressmaker.Location.Latitude,
			&dressmaker.Location.Longitude,
			&services,
			&dressmaker.Grade,
			&dressmaker.CreatedAt,
			&dressmaker.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("error on scan rows: %v\n", err)
			return nil, err
		}

		dressmaker.Services = strings.Split(services, ", ")
		dressmakers = append(dressmakers, dressmaker)
	}

	return dressmakers, nil
}

func (r *DressmakerRepository) FindByServices(services string) ([]entity.Dressmaker, error) {
	var dressmakers []entity.Dressmaker

	keywords := strings.Split(services, ", ")
	query := generateILikeQuery(keywords)

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		fmt.Printf("error on prepare query: %v\n", err)
		return nil, err
	}

	args := make([]interface{}, len(keywords))
	for i, keyword := range keywords {
		args[i] = "%" + keyword + "%"
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Printf("error on exec query: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var dressmaker entity.Dressmaker
		var services string

		err := rows.Scan(
			&dressmaker.ID,
			&dressmaker.Email,
			&dressmaker.Password,
			&dressmaker.Name,
			&dressmaker.Contact,
			&dressmaker.Location.Latitude,
			&dressmaker.Location.Longitude,
			&services,
			&dressmaker.Grade,
			&dressmaker.CreatedAt,
			&dressmaker.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("error on scan rows: %v\n", err)
			return nil, err
		}

		dressmaker.Services = strings.Split(services, ", ")
		dressmakers = append(dressmakers, dressmaker)
	}

	return dressmakers, nil
}

func (r *DressmakerRepository) Update(dressmaker *entity.Dressmaker) error {
	stmt, err := r.DB.Prepare(`
		UPDATE dressmakers
		SET name = $1, contact = $2, location = ST_MakePoint($3, $4)::geography, services = $5, grade = $6, updated_at = $7
		WHERE id = $8
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		dressmaker.Name,
		dressmaker.Contact,
		dressmaker.Location.Latitude,
		dressmaker.Location.Longitude,
		strings.Join(dressmaker.Services, ","),
		dressmaker.Grade,
		dressmaker.UpdatedAt,
		dressmaker.ID,
	)

	return err
}

func (r *DressmakerRepository) GetServices() ([]string, error) {
	uniqueServicesMap := make(map[string]bool)
	var servicesUniqueAttr []string

	rows, err := r.DB.Query(`
		SELECT 
			services
		FROM dressmakers
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var services string

		err := rows.Scan(
			&services,
		)
		if err != nil {
			return nil, err
		}

		splitedServices := strings.Split(services, ",")
		for _, service := range splitedServices {
			if _, found := uniqueServicesMap[service]; !found {
				uniqueServicesMap[service] = true
				servicesUniqueAttr = append(servicesUniqueAttr, service)
			}
		}

	}

	return servicesUniqueAttr, nil
}

func generateILikeQuery(keywords []string) string {
	baseQuery := `
        SELECT
            id,
            email,
            password,
            name,
            contact,
            ST_Y(location::geometry) as latitude,
            ST_X(location::geometry) as longitude,
            services,
            grade,
            created_at,
            updated_at
        FROM dressmakers
        WHERE `

	var conditions []string
	for i := range keywords {
		conditions = append(conditions, fmt.Sprintf("services ILIKE $%d", i+1))
	}

	fullQuery := baseQuery + strings.Join(conditions, " OR ")

	return fullQuery
}

func generateFullyILikeQuery(keywords []string) string {
	baseQuery := `
		SELECT
			id,
			email,
			password,
			name,
			contact,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			services,
			grade,
			created_at,
			updated_at
		FROM dressmakers
		WHERE 
			ST_DWithin(
				location,
				ST_MakePoint($1, $2)::geography,
				$3
			)
	`

	var conditions []string
	for i := range keywords {
		conditions = append(conditions, fmt.Sprintf("services ILIKE $%d", i+1))
	}

	fullQuery := baseQuery + strings.Join(conditions, " AND ")

	return fullQuery
}
