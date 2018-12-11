package mysql

type Columns struct {
	TableCatalog *string  `json:"TABLE_CATALOG"`
	TableSchema *string  `json:"TABLE_SCHEMA"`
	TableName *string  `json:"TABLE_NAME"`
	ColumnName string  `json:"COLUMN_NAME"`
	OrdinalPosition *int64  `json:"ORDINAL_POSITION"`
	ColumnDefault *string  `json:"COLUMN_DEFAULT"`
	IsNullAble string  `json:"IS_NULLABLE"`
	DataType *string  `json:"DATA_TYPE"`
	CharacterMaximumLength *int64  `json:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength *int64  `json:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision *int64  `json:"NUMERIC_PRECISION"`
	NumericScale *int64  `json:"NUMERIC_SCALE"`
	DatetimePrecision *int64  `json:"DATETIME_PRECISION"`
	CharacterSetName *string  `json:"CHARACTER_SET_NAME"`
	CollationName *string  `json:"COLLATION_NAME"`
	ColumnType string  `json:"COLUMN_TYPE"`
	ColumnKey *string  `json:"COLUMN_KEY"`
	Extra *string  `json:"EXTRA"`
	Privileges *string  `json:"PRIVILEGES"`
	ColumnComment *string  `json:"COLUMN_COMMENT"`
	GenerationExpression *interface{}  `json:"GENERATION_EXPRESSION"`
}