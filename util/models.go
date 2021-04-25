/*
Copyright © 2021 Brian Longmore brianl.ext@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import "time"

type (
	Cells  []Cell
	Frames []Frame
)

type Cell struct {
	X               int
	Y               int
	IsAlive         bool
	LiveNeighborNum int
}

type Frame struct {
	FrameNum int
	Cells
}

type Pattern struct {
	FilePath string
	X        int
	Y        int
}

type GameData struct {
	BuildConfig
	Frames
}

type BuildConfig struct {
	BuildFilePath string
	Width         int
	Height        int
	FrameCount    int
	InitPattern   Pattern
}

type RunConfig struct {
	BuildFilePath string
	Interval      time.Duration
	LiveCellChar  string
	DeadCellChar  string
}

type GameConfig struct {
	Width        int
	Height       int
	FrameCount   int
	InitPattern  Pattern
	Interval     time.Duration
	LiveCellChar string
	DeadCellChar string
}
