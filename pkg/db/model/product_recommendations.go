package model

// TableNameProductRecommendation is the table name of <product_recommendations>
const TableNameProductRecommendation = "product_recommendations"

// ProductRecommendation mapped from table <product_recommendations>
type ProductRecommendation struct {
	ID   int    `gorm:"column:id;type:int;not null;primaryKey;autoIncrement"`
	Name string `gorm:"column:name;type:varchar(256);not null"`
}

// TableName ProductRecommendation's table name
func (*ProductRecommendation) TableName() string {
	return TableNameProductRecommendation
}
