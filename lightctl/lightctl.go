package lightctl

type LightState int32

const (
	LightStateNone      LightState = 0
	LightStateBlue      LightState = 1
	LightStateRed       LightState = 2
	LightStateGreen     LightState = 3
	LightStateRedBlue   LightState = 4
	LightStateRedGreen  LightState = 5
	LightStateBlueGreen LightState = 6
	LightStateAll       LightState = 7
	LightStateStrobe    LightState = 8
)

func (l LightState) String() string {
	switch l {
	case LightStateBlue:
		return "blue"
	case LightStateRed:
		return "red"
	case LightStateGreen:
		return "green"
	case LightStateRedBlue:
		return "redblue"
	case LightStateRedGreen:
		return "redgreen"
	case LightStateBlueGreen:
		return "bluegreen"
	case LightStateAll:
		return "all"
	case LightStateStrobe:
		return "strobe"
	default:
		return "none"
	}
}

func ParseLightState(color string) LightState {
	switch color {
	case "blue":
		return LightStateBlue
	case "red":
		return LightStateRed
	case "green":
		return LightStateGreen
	case "redblue":
		return LightStateRedBlue
	case "redgreen":
		return LightStateRedGreen
	case "bluegreen":
		return LightStateBlueGreen
	case "all":
		return LightStateAll
	case "strobe":
		return LightStateStrobe
	default:
		return LightStateNone
	}
}