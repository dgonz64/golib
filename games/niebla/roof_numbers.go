package main

import (
	"strconv"

	"golib"
)

const (
	roofNumberSpan   = 0.6
	roofNumberOffset = 0.1
)

var roofDigits = [10][5]uint8{
	{0b111, 0b101, 0b101, 0b101, 0b111},
	{0b010, 0b110, 0b010, 0b010, 0b111},
	{0b111, 0b001, 0b111, 0b100, 0b111},
	{0b111, 0b001, 0b111, 0b001, 0b111},
	{0b101, 0b101, 0b111, 0b001, 0b001},
	{0b111, 0b100, 0b111, 0b001, 0b111},
	{0b111, 0b100, 0b111, 0b101, 0b111},
	{0b111, 0b001, 0b010, 0b010, 0b010},
	{0b111, 0b101, 0b111, 0b101, 0b111},
	{0b111, 0b101, 0b111, 0b001, 0b111},
}

func drawFactoryRoofNumber(
	screen *golib.Screen,
	number int,
	gx, roofY, across float32,
) {
	if number <= 0 {
		return
	}
	digits := strconv.Itoa(number)
	width := float32(len(digits)*4 - 1)
	step := across * roofNumberSpan / max(width, 5)
	roofPoint := func(x, y float32) golib.Vector2 {
		u := (x-width/2)*step + across*roofNumberOffset
		v := (y-2.5)*step - across*roofNumberOffset
		return golib.Vector2{
			X: gx + (u-v)*unitW/2,
			Y: roofY + (u+v)*unitH/2,
		}
	}
	for i, digit := range digits {
		for row, bits := range roofDigits[digit-'0'] {
			for col := 0; col < 3; col++ {
				if bits&(1<<uint(2-col)) == 0 {
					continue
				}
				x, y := float32(i*4+col), float32(row)
				screen.DrawPolygon([]golib.Vector2{
					roofPoint(x, y), roofPoint(x+1, y),
					roofPoint(x+1, y+1), roofPoint(x, y+1),
				}, warFactoryInk)
			}
		}
	}
}
