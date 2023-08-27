package database

type Database interface {
  DoInit() error;
  AllRows() ([]TokenEntity, error);
  SearchRow(UUID string) (TokenEntity, error);
  InsertRow(t TokenEntity) error;
  UpdateRow(UUID string, t TokenEntity) error;
  DeleteRow(UUID string) error;
}
