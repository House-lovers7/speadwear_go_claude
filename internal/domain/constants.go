package domain

// Season constants
const (
	SeasonSpring     = 1
	SeasonSummer     = 2
	SeasonAutumn     = 3
	SeasonWinter     = 4
	SeasonAllSeason  = 5
)

// TPO (Time, Place, Occasion) constants
const (
	TPOWork   = 1
	TPOCasual = 2
	TPOFormal = 3
	TPOSports = 4
	TPOHome   = 5
)

// Color constants
const (
	ColorBlack  = 1
	ColorWhite  = 2
	ColorGray   = 3
	ColorBrown  = 4
	ColorBeige  = 5
	ColorGreen  = 6
	ColorBlue   = 7
	ColorPurple = 8
	ColorYellow = 9
	ColorPink   = 10
	ColorRed    = 11
	ColorOrange = 12
	ColorSilver = 13
	ColorGold   = 14
	ColorOther  = 15
)

// Size constants for coordinates
const (
	// Top Length
	TopLengthCrop   = 1
	TopLengthNormal = 2
	TopLengthLong   = 3

	// Top Sleeve
	TopSleeveNone  = 1
	TopSleeveCap   = 2
	TopSleeveShort = 3
	TopSleeveHalf  = 4
	TopSleeveLong  = 5

	// Bottom Length
	BottomLengthMini  = 1
	BottomLengthShort = 2
	BottomLengthKnee  = 3
	BottomLengthMidi  = 4
	BottomLengthLong  = 5
	BottomLengthMaxi  = 6

	// Bottom Type
	BottomTypeSkirt = 1
	BottomTypePants = 2

	// Dress Length
	DressLengthMini  = 1
	DressLengthShort = 2
	DressLengthKnee  = 3
	DressLengthMidi  = 4
	DressLengthLong  = 5
	DressLengthMaxi  = 6

	// Dress Sleeve
	DressSleeveNone  = 1
	DressSleeveCap   = 2
	DressSleeveShort = 3
	DressSleeveHalf  = 4
	DressSleeveLong  = 5

	// Outer Length
	OuterLengthShort  = 1
	OuterLengthNormal = 2
	OuterLengthLong   = 3

	// Outer Sleeve
	OuterSleeveShort = 1
	OuterSleeveHalf  = 2
	OuterSleeveLong  = 3
)

// Notification actions
const (
	NotificationActionFollow      = "follow"
	NotificationActionLike        = "like"
	NotificationActionComment     = "comment"
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

// SeasonNames maps season constants to their names
var SeasonNames = map[int]string{
	SeasonSpring:    "春",
	SeasonSummer:    "夏",
	SeasonAutumn:    "秋",
	SeasonWinter:    "冬",
	SeasonAllSeason: "オールシーズン",
}

// TPONames maps TPO constants to their names
var TPONames = map[int]string{
	TPOWork:   "仕事",
	TPOCasual: "カジュアル",
	TPOFormal: "フォーマル",
	TPOSports: "スポーツ",
	TPOHome:   "ホーム",
}

// ColorNames maps color constants to their names
var ColorNames = map[int]string{
	ColorBlack:  "ブラック",
	ColorWhite:  "ホワイト",
	ColorGray:   "グレー",
	ColorBrown:  "ブラウン",
	ColorBeige:  "ベージュ",
	ColorGreen:  "グリーン",
	ColorBlue:   "ブルー",
	ColorPurple: "パープル",
	ColorYellow: "イエロー",
	ColorPink:   "ピンク",
	ColorRed:    "レッド",
	ColorOrange: "オレンジ",
	ColorSilver: "シルバー",
	ColorGold:   "ゴールド",
	ColorOther:  "その他",
}