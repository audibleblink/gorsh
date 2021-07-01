//build +windows

package enum

import _ "embed"

//go:embed embed/sherlock.ps1
var sherlock []byte

//go:embed embed/jaws.ps1
var jaws []byte

//go:embed embed/powerup.ps1
var powerUp []byte

//go:embed embed/winpeas.bat
var winPeas []byte

func Sherlock() *EnumScript {
	return &EnumScript{Data: sherlock}
}

func Jaws() *EnumScript {
	return &EnumScript{Data: jaws}
}

func PowerUp() *EnumScript {
	return &EnumScript{Data: powerUp}
}

func WinPeas() *EnumScript {
	return &EnumScript{Data: winPeas}
}
