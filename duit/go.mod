module ccal-duit

go 1.14

require (
	9fans.net/go v0.0.2 // indirect
	//github.com/Aquarian-Age/ts v0.0.0
	github.com/mjl-/duit v0.0.0-20200330125617-580cb0b2843f
	github.com/nongli v0.0.0
	github.com/qxqm v0.0.0
	github.com/sjqm v0.0.0
)

replace (
	9fans.net/go v0.0.0-00010101000000-000000000000 => /home/xuan/src/mjl-duit/go
	//github.com/Aquarian-Age/ts => /home/xuan/src/sr.ht/ts/ts通书/
	github.com/nongli => /home/xuan/src/ccal-cli/
	github.com/qxqm => /home/xuan/src/qxqm/QXQM禽星奇门/
	github.com/sjqm => /home/xuan/src/sjqm/sjqm时家奇门/
)
