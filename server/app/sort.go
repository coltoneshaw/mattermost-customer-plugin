package app

// SortField enumerates the available fields we can sort on.
type SortField string

const (
	SortByName       SortField = "name"
	SortByCSM        SortField = "customerSuccessManager"
	SortByAE         SortField = "accountExecutive"
	SortByTAM        SortField = "technicalAccountManager"
	SortByType       SortField = "type"
	SortBySiteURL    SortField = "siteURL"
	SortByLicensedTo SortField = "licensedTo"
)

// SortDirection is the type used to specify the ascending or descending order of returned results.
type SortDirection string

const (
	// DirectionDesc is descending order.
	DirectionDesc SortDirection = "DESC"

	// DirectionAsc is ascending order.
	DirectionAsc SortDirection = "ASC"
)

func IsValidDirection(direction SortDirection) bool {
	return direction == DirectionAsc || direction == DirectionDesc
}
