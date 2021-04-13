/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

type FrameCells []Cell

type CellNeighbors []Cell

type Pattern struct {
	Path string
	X    int
	Y    int
}

type FrameConfig struct {
	Width          int
	Height         int
	FrameCount     int
	Interval       time.Duration
	DeadCellChar   string
	LivingCellChar string
}

type Cell struct {
	X               int
	Y               int
	IsAlive         bool
	LivingNeighbors CellNeighbors
}
