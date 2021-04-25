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
	X               int  `json:"x"`
	Y               int  `json:"y"`
	IsAlive         bool `json:"isAlive"`
	LiveNeighborNum int  `json:"LiveNeighborNum"`
}

type Frame struct {
	FrameNum int `json:"frameNum"`
	Cells    `json:"cells"`
}

type Pattern struct {
	FilePath string `json:"filePath"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type GameData struct {
	BuildConfig `json:"buildConfig"`
	Frames      `json:"frames"`
}

type BuildConfig struct {
	BuildFilePath string  `json:"buildFilePath"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	FrameCount    int     `json:"frameCount"`
	InitPattern   Pattern `json:"initPattern"`
}

type RunConfig struct {
	BuildFilePath string        `json:"buildFilePath"`
	Interval      time.Duration `json:"interval"`
	LiveCellChar  string        `json:"liveCellChar"`
	DeadCellChar  string        `json:"deadCellChar"`
}

type GameConfig struct {
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	FrameCount   int           `json:"frameCount"`
	InitPattern  Pattern       `json:"initPattern"`
	Interval     time.Duration `json:"interval"`
	LiveCellChar string        `json:"liveCellChar"`
	DeadCellChar string        `json:"deadCellChar"`
}
