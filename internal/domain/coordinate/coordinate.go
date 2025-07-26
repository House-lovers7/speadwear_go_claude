package coordinate

import (
	"github.com/House-lovers7/speadwear-go/internal/domain"
)

type Coordinate struct {
	domain.BaseModel
	UserID           uint    `gorm:"not null;index" json:"user_id"`
	Season           int     `json:"season"`
	TPO              int     `json:"tpo"`
	Picture          string  `gorm:"type:varchar(255)" json:"picture"`
	SiTopLength      int     `json:"si_top_length"`
	SiTopSleeve      int     `json:"si_top_sleeve"`
	SiBottomLength   int     `json:"si_bottom_length"`
	SiBottomType     int     `json:"si_bottom_type"`
	SiDressLength    int     `json:"si_dress_length"`
	SiDressSleeve    int     `json:"si_dress_sleeve"`
	SiOuterLength    int     `json:"si_outer_length"`
	SiOuterSleeve    int     `json:"si_outer_sleeve"`
	SiShoeSize       int     `json:"si_shoe_size"`
	Memo             string  `gorm:"type:text" json:"memo"`
	Rating           float32 `json:"rating"`
	
	// Relations
	User            User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items           []Item           `gorm:"foreignKey:CoordinateID" json:"items,omitempty"`
	Comments        []Comment        `gorm:"foreignKey:CoordinateID" json:"comments,omitempty"`
	LikeCoordinates []LikeCoordinate `gorm:"foreignKey:CoordinateID" json:"like_coordinates,omitempty"`
}

type Comment struct {
	domain.BaseModel
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	CoordinateID uint       `gorm:"not null;index" json:"coordinate_id"`
	Comment      string     `gorm:"type:text;not null" json:"comment"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate   Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

type LikeCoordinate struct {
	domain.BaseModel
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	CoordinateID uint       `gorm:"not null;index" json:"coordinate_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate   Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

// Size定数
const (
	// Top Length
	TopLengthCrop = 1
	TopLengthNormal = 2
	TopLengthLong = 3

	// Top Sleeve
	TopSleeveNone = 1
	TopSleeveCap = 2
	TopSleeveShort = 3
	TopSleeveHalf = 4
	TopSleeveLong = 5

	// Bottom Length
	BottomLengthMini = 1
	BottomLengthShort = 2
	BottomLengthKnee = 3
	BottomLengthMidi = 4
	BottomLengthLong = 5
	BottomLengthMaxi = 6

	// Bottom Type
	BottomTypeSkirt = 1
	BottomTypePants = 2

	// Dress Length
	DressLengthMini = 1
	DressLengthShort = 2
	DressLengthKnee = 3
	DressLengthMidi = 4
	DressLengthLong = 5
	DressLengthMaxi = 6

	// Dress Sleeve
	DressSleeveNone = 1
	DressSleeveCap = 2
	DressSleeveShort = 3
	DressSleeveHalf = 4
	DressSleeveLong = 5

	// Outer Length
	OuterLengthShort = 1
	OuterLengthNormal = 2
	OuterLengthLong = 3

	// Outer Sleeve
	OuterSleeveShort = 1
	OuterSleeveHalf = 2
	OuterSleeveLong = 3
)

// 以下の型は他のパッケージで定義されるが、ここでは参照のために宣言
type User struct{}
type Item struct{}