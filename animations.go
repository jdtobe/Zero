package main

import (
	"context"

	"github.com/jdtobe/Zero/animation"
	"github.com/jdtobe/Zero/color"
)

const xmasTreeShowName = "xmas-tree"

const originalShowName = "original"

type animationFunc func(context.Context, []color.HSV) error

func animationByName(name string) animationFunc {
	switch name {
	case xmasTreeShowName:
		return animation.XmasTree
	case originalShowName:
		fallthrough
	default:
		return animation.Original
	}
}
