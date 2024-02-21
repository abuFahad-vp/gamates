package main

import (
	//"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
  SW = 600
  SH = 600
)

type Grid struct {
  color color.RGBA
  greyFunc func(float64, float64) float64
}

type Pos struct {
  X1,Y1,X2,Y2 int
}

func main() {
	rl.InitWindow(SW, SH, "Gamates")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
  nX := 4
  nY := 4
  if len(os.Args) == 3 {
    nX,_= strconv.Atoi(os.Args[1])
    nY,_ = strconv.Atoi(os.Args[2])
  }

  grid,ans:= generateGrid(nX,nY)

  // ------ STATES -----------
  nowX := 999
  nowY := 999

  prevX := 999
  prevY := 999

  founds := make(map[Pos]bool)

  score := 0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
    //----------------------------
    DrawGrid(float64(nX),float64(nY),grid,ans,&nowX,&nowY,&prevX,&prevY,&founds,&score)
    //fmt.Println(founds)
    //----------------------------
		rl.EndDrawing()
	}
}

func DrawGrid(nX,nY float64, grid [][]Grid,ans map[Pos]int,nowX,nowY,prevX,prevY *int,founds *map[Pos]bool, score *int) {

  paddingX := SW * 0.01
  paddingY := SH * 0.01

  gridX := (SW - paddingX) / nX
  gridY := (SH - paddingY) / nY

  width := gridX - paddingX
  height := gridY - paddingY

  //fmt.Println(gridX,gridY,paddingX,paddingY,width,height)

  if *score >= int(nX * nY / 2) {
    msg := "YOU WIN!!!!"
    fontsize := 20
    msgWidth := rl.MeasureText(msg,int32(fontsize))
    rl.DrawText(msg,(SW-msgWidth)/2,SH/2-20,int32(fontsize),rl.Black)
    return;
  }

  for i,x := paddingX,0;i < SW;i,x = i + paddingX + width, x + 1 {
    for k,y := paddingY,0; k < SH; k,y = k + paddingY + height, y + 1 {
      //fmt.Println(x,y)
      if _,exits := (*founds)[Pos{x,y,0,0}];exits {
        rl.DrawRectangleV(rl.Vector2{X:float32(i),Y:float32(k)},rl.Vector2{X: float32(width),Y: float32(height)},rl.RayWhite)
      } else if *nowX == x && *nowY == y || *prevX == x && *prevY ==  y {
        if *nowX != 999 && *prevX != 999 && (ans[Pos{*nowX,*nowY,0,0}] == ans[Pos{*prevX,*prevY,0,0}]){
            (*founds)[Pos{*nowX,*nowY,0,0}] = true
            (*founds)[Pos{*prevX,*prevY,0,0}] = true
            (*score)++
          }
        DrawGridOfPatterns(i,k,i+width,k+height,grid[x][y].color,grid[x][y].greyFunc)
      } else {
        rl.DrawRectangleV(rl.Vector2{X:float32(i),Y:float32(k)},rl.Vector2{X: float32(width),Y: float32(height)},rl.Black)
      }

      mousePos := rl.GetMousePosition()
      if rl.CheckCollisionPointRec(mousePos,
      rl.Rectangle{X: float32(i),Y: float32(k),Width: float32(width),Height: float32(height)}) {
        if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
          //fmt.Printf("(%d,%d), (%d,%d), (%d,%d)\n",*prevX,*prevY,*nowX,*nowY,x,y)
          *prevX = 999
          *prevY = 999
          if *nowX != x || *nowY != y {
            *prevX = *nowX
            *prevY = *nowY
            //fmt.Println(*click)
          }
          *nowX = x
          *nowY = y
        }
      }
    }
  }
}

func DrawGridOfPatterns(x1,y1,x2,y2 float64, color color.RGBA, f func (float64,float64) float64) {
  for i,x := 0.0,x1;x<x2;i,x = i+1,x+1 {
    for k,y := 0.0,y1;y<y2;k,y = k+1,y+1 {
      r,g,b := color.R, color.G, color.B
      rl.DrawPixelV(rl.Vector2{X: float32(x),Y: float32(y)},rl.NewColor(r,g,b,uint8(f(i,k) * 255 / (x2-x1))))
    }
  }
}

func generateGrid(nX, nY int) ([][]Grid,map[Pos]int) {

  elements := make([]Grid, nX*nY/2)

  for index,i := range rand.Perm(nX*nY/2){
    elements[i].color = colors[index % len(colors)]
  }

  for index,i := range rand.Perm(nX*nY/2){
    elements[i].greyFunc = functions[index % len(functions)]
  }

  rand_pos := rand.Perm(nX*nY)
  index := 0
  grid := make([][]Grid,nX)

  ans := make(map[Pos]int)

  for i := range grid {
    grid[i] = make([]Grid,nY)
    for k := range grid[i] {
      indexE := rand_pos[index] % (nX * nY / 2)
      grid[i][k] = elements[indexE]
      ans[Pos{i,k,0,0}] = indexE
      //fmt.Println(index,indexE,i,k)
      index++
    }
  }

  //fmt.Println(rand_pos)
  //fmt.Println(elements[0])
  //fmt.Println(nX*nY/2)
  //fmt.Println(grid)
  return grid,ans
}

var colors []color.RGBA = []color.RGBA {
  rl.Red,
  rl.DarkGreen,
  rl.Blue,
  rl.Black,
  rl.Magenta,
  rl.Brown,
  rl.DarkPurple,
}

var functions []func(float64, float64) float64 = []func(float64, float64) float64 {
    func(f1, f2 float64) float64 {return 250 + 250 * math.Cos(math.Sqrt(f1*f1 + f2*f2))},
    func(f1, f2 float64) float64 {return (f1 + f2)/2.0},
    func(f1, f2 float64) float64 {return (f1 * f2)},
    func(f1, f2 float64) float64 {return f1*f2*3.14},
    func(f1, f2 float64) float64 {return (f1 - f2)/2.0},
    func(f1, f2 float64) float64 {return f1*f1 + f2*f2},
}
