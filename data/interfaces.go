package data

import (
	"database/sql"

	"github.com/aykleo/ion/data/sqlite"
)

type IData interface {
	GetUser() User
	GetData() IData
	SetDB(db *sql.DB)

	GetInitialData(path string) (IData, bool)

	SetUsername(name, path string)

	AddSecret(args []string, path string) error
	UpdateSecretValue(args []string, path string) error
	UpdateSecretName(args []string, path string) error
	UpdateSecretTags(args []string, path string) error
	ListSecrets(args []string, path string) ([]Secret, bool, error)
	SearchSecret(args []string) ([]Secret, error)
	RemoveSecret(args []string, path string) error
	CopySecretToClipboard(args []string, path string) error

	AddAlias(args []string, path string) error
}

func NewData() IData {
	data := InitData()
	return &data
}

func (s *Data) GetData() IData {
	return s
}

func (s *Data) SetDB(db *sql.DB) {
	s.db = db
}

func (s *Data) GetInitialData(path string) (IData, bool) {

	if s.db == nil {
		return s, false
	}

	if err := s.loadDataFromDB(); err != nil {
		return s, false
	}

	return s, true
}

func (s *Data) loadDataFromDB() error {
	sqliteUser, err := sqlite.GetUser(s.db)
	if err != nil {
		return err
	}
	s.User = User{Username: sqliteUser.Username}

	sqliteSecrets, err := sqlite.GetSecrets(s.db)
	if err != nil {
		return err
	}
	s.Secrets = make([]Secret, len(sqliteSecrets))
	for i, sqliteSecret := range sqliteSecrets {
		s.Secrets[i] = Secret{
			ID:        sqliteSecret.ID,
			Name:      sqliteSecret.Name,
			Salt:      sqliteSecret.Salt,
			Value:     sqliteSecret.Value,
			Tags:      sqliteSecret.Tags,
			CreatedAt: sqliteSecret.CreatedAt,
			UpdatedAt: sqliteSecret.UpdatedAt,
		}
	}
	s.buildSecretIndex()

	sqliteAliases, err := sqlite.GetAliases(s.db)
	if err != nil {
		return err
	}
	s.Aliases = make([]Alias, len(sqliteAliases))
	for i, sqliteAlias := range sqliteAliases {
		s.Aliases[i] = Alias{
			ID:        sqliteAlias.ID,
			Name:      sqliteAlias.Name,
			Command:   sqliteAlias.Command,
			CreatedAt: sqliteAlias.CreatedAt,
			UpdatedAt: sqliteAlias.UpdatedAt,
		}
	}

	return nil
}
