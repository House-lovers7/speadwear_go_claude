package item

import (
	"github.com/House-lovers7/speadwear-go/internal/domain"
)

type Item struct {
	domain.BaseModel
	UserID       uint    `gorm:"not null;index" json:"user_id"`
	CoordinateID *uint   `json:"coordinate_id,omitempty"`
	SuperItem    string  `gorm:"type:varchar(100)" json:"super_item"`
	Season       int     `json:"season"`
	TPO          int     `json:"tpo"`
	Color        int     `json:"color"`
	Content      string  `gorm:"type:text" json:"content"`
	Memo         string  `gorm:"type:text" json:"memo"`
	Picture      string  `gorm:"type:varchar(255)" json:"picture"`
	Rating       float32 `json:"rating"`
	
	// Relations
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate  *Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

// 定数定義
const (
	// Season
	SeasonSpring = 1
	SeasonSummer = 2
	SeasonAutumn = 3
	SeasonWinter = 4
	SeasonAllSeason = 5

	// TPO (Time, Place, Occasion)
	TPOWork = 1
	TPOCasual = 2
	TPOFormal = 3
	TPOSports = 4
	TPOHome = 5

	// Color
	ColorBlack = 1
	ColorWhite = 2
	ColorGray = 3
	ColorBrown = 4
	ColorBeige = 5
	ColorGreen = 6
	ColorBlue = 7
	ColorPurple = 8
	ColorYellow = 9
	ColorPink = 10
	ColorRed = 11
	ColorOrange = 12
	ColorSilver = 13
	ColorGold = 14
	ColorOther = 15
)

// SuperItem categories
var SuperItemCategories = []string{
	"アウター",
	"トップス",
	"ボトムス",
	"ワンピース",
	"シューズ",
	"バッグ",
	"アクセサリー",
	"帽子",
	"その他",
}

// 以下の型は他のパッケージで定義されるが、ここでは参照のために宣言
type User struct{}
type Coordinate struct{}