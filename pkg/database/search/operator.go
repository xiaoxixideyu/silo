package search

// Operator use mongo
// https://www.mongodb.com/docs/manual/reference/operator/query/
const (
	OPEqual            = "$eq"
	OPEqualIgnoreCase  = "$eqIgnoreCase"
	OPNotEqual         = "$ne"
	OPGT               = "$gt"
	OPGTE              = "$gte"
	OPLT               = "$lt"
	OPLTE              = "$lte"
	OPIn               = "$in"
	OPNotIn            = "$notIn"
	OPOr               = "$or"
	OPBetween          = "$between"
	OPLike             = "$like"
	OPSort             = "$sort"
	OPContains         = "$contains"
	OPLeftJoin         = "$left"
	OPJoin             = "$join"
	OPGroup            = "$group"
	OPGroupOr          = "$groupor"
	OPSet              = "$set"
	OPUnSet            = "$unset"
	OPSetOnInsert      = "$setOnInsert"
	OPInc              = "$inc"
	OPIS               = "$is"
	OPPattern          = "$pattern"
	OPArrayContains    = "$arrayContains" // 包含数组的所有元素
	OPExpr             = "$expr"
	OPIsNotNil         = "$isNotNil"
	OPIsNil            = "$isNil"
	OPJson             = "$jsonValue"
	OPJsonLike         = "$jsonLike"
	OPJsonToStringLike = "$jsonToStringLike" // Convert field to JSON string for LIKE operations
)

// Item Order
const (
	SortASC           = "ASC"
	SortDESC          = "DESC"
	SortDESCNullsLast = "DESC NULLS LAST"
)

// Custom Operators
const (
	OPCustomSort = "$customSort"
)
