package database

import (
	"fmt"
	"log"

	"internal-transfer-system/internal/config"
	"internal-transfer-system/internal/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Enable foreign key constraints
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings for better concurrency
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(25)
	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(5)

	log.Printf("Successfully connected to database: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&repository.Account{}, &repository.Transaction{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Add additional constraints and indexes
	err = addConstraints(db)
	if err != nil {
		return fmt.Errorf("failed to add constraints: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func addConstraints(db *gorm.DB) error {
	// Add composite index for transaction queries
	err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_accounts ON transactions(source_account_id, destination_account_id)").Error
	if err != nil {
		return err
	}

	// Add index for transaction amount queries
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_amount ON transactions(amount)").Error
	if err != nil {
		return err
	}

	// Add index for account_id queries
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_accounts_account_id ON accounts(account_id)").Error
	if err != nil {
		return err
	}

	// Add check constraint to ensure no self-transfers
	err = addCheckConstraint(db)
	if err != nil {
		return err
	}

	log.Println("Database constraints added successfully")
	return nil
}

func addCheckConstraint(db *gorm.DB) error {
	// Check if constraint already exists
	var count int64
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.table_constraints 
		WHERE constraint_name = 'check_no_self_transfer' 
		AND table_name = 'transactions'
	`).Count(&count).Error
	
	if err != nil {
		return err
	}
	
	// Only add constraint if it doesn't exist
	if count == 0 {
		err = db.Exec("ALTER TABLE transactions ADD CONSTRAINT check_no_self_transfer CHECK (source_account_id != destination_account_id)").Error
		if err != nil {
			return err
		}
	}
	
	return nil
}

func Setup(cfg *config.Config) (*gorm.DB, error) {
	db, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}