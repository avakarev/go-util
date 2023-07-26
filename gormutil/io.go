package gormutil

// TableFilter defines table filtering options
type TableFilter struct {
	IncludeTables []string `json:"includeTables"`
	ExcludeTables []string `json:"excludeTables"`
}

// contains checks whenther given element is a member of given slice
func contains[T string | int | int64 | float64](slice []T, elem T) bool {
	for _, e := range slice {
		if elem == e {
			return true
		}
	}
	return false
}

// FilterTables returns tables with respect of include/exclude filters
func FilterTables(tables []string, f *TableFilter) []string {
	if f == nil || (len(f.IncludeTables) == 0 && len(f.ExcludeTables) == 0) {
		return tables
	}
	filteredTables := make([]string, 0)
	for _, t := range tables {
		if !contains(f.ExcludeTables, t) && (len(f.IncludeTables) == 0 || contains(f.IncludeTables, t)) {
			filteredTables = append(filteredTables, t)
		}
	}
	return filteredTables
}

// Tables retuns list of existing db tables respecing the given filter
func (db *DB) Tables(filter *TableFilter) ([]string, error) {
	tables, err := db.Conn().Migrator().GetTables()
	if err != nil {
		return nil, err
	}
	return FilterTables(tables, filter), nil
}

// ExportTable returns all rows of given table as slice of maps
func (db *DB) ExportTable(table string) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	if err := db.Conn().Table(table).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// Export returns map of all tables with their rows
func (db *DB) Export(filter *TableFilter) (map[string]interface{}, error) {
	tables, err := db.Tables(filter)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for _, t := range tables {
		data, err := db.ExportTable(t)
		if err != nil {
			return nil, err
		}
		m[t] = data
	}

	return m, nil
}

// ImportTable inserts given rows into given table
func (db *DB) ImportTable(table string, rows []interface{}) error {
	for _, row := range rows {
		if err := db.Conn().Table(table).Create(row).Error; err != nil {
			return err
		}
	}
	return nil
}

// Import inserts given data into db
func (db *DB) Import(data map[string]interface{}, filter *TableFilter) error {
	tables, err := db.Tables(filter)
	if err != nil {
		return err
	}

	for _, t := range tables {
		value, ok := data[t]
		if !ok {
			continue
		}
		rows, ok := value.([]interface{})
		if !ok {
			continue
		}
		if err := db.ImportTable(t, rows); err != nil {
			return err
		}
	}

	return nil
}
